package server

import (
	"StockOverflow/internal/database"
	"StockOverflow/internal/exchange"
	"StockOverflow/internal/pool"
	"StockOverflow/pkg/xmlparser"
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/shopspring/decimal"
)

// AccountNode represents a simple account structure
type AccountNode struct {
	ID        string
	Balance   decimal.Decimal
	Positions map[string]decimal.Decimal //Map of an account to it's hold of symbol
}

// Server represents the exchange server
// Server represents the exchange server
type Server struct {
	// Server configuration
	listener    net.Listener
	logger      *log.Logger
	wg          sync.WaitGroup
	connections map[net.Conn]struct{}
	mutex       sync.Mutex

	// Database connection
	db *sql.DB

	// Exchange state
	stockPool   *pool.StockPool         // Stock trading nodes
	accounts    map[string]*AccountNode // Simple account storage
	nextOrderID int                     // For generating unique order IDs
	exchange    *exchange.Exchange      // Exchange engine for matching orders

	// Mutexes for concurrent access
	accountsMutex sync.RWMutex
	idMutex       sync.Mutex
}

// NewServer creates a new exchange server
func NewServer(logger *log.Logger) *Server {
	stockPool := pool.NewPool(1000)

	server := &Server{
		logger:      logger,
		connections: make(map[net.Conn]struct{}),
		stockPool:   stockPool,
		accounts:    make(map[string]*AccountNode),
	}

	return server
}

// SetDB sets the database connection and initializes the exchange
func (s *Server) SetDB(db *sql.DB) {
	s.db = db
	s.exchange = exchange.NewExchange(db, s.stockPool, s.logger)
	// Initialize nextOrderID from database
	maxID, err := database.GetMaxOrderID(db)
	if err != nil {
		s.logger.Printf("Warning: Failed to get max order ID from database: %v", err)
		// Fall back to default if there's an error
		s.nextOrderID = 1
	} else {
		// If there are existing orders, use max ID + 1
		// Otherwise, start from 11 as requested
		if maxID > 0 {
			s.nextOrderID = maxID + 1
			s.logger.Printf("Initialized nextOrderID to %d based on database", s.nextOrderID)
		} else {
			s.nextOrderID = 1
			s.logger.Printf("No existing orders found, initializing nextOrderID to %d", s.nextOrderID)
		}
	}
}

// Start begins listening for connections on the specified address
func (s *Server) Start(addr string) error {
	var err error
	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start listener: %v", err)
	}

	s.logger.Printf("Server listening on %s", addr)

	// Accept connections in a loop
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			// Check if the listener was closed
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil
			}
			s.logger.Printf("Error accepting connection: %v", err)
			continue
		}

		// Track the connection
		s.mutex.Lock()
		s.connections[conn] = struct{}{}
		s.mutex.Unlock()

		// Handle the connection in a separate goroutine
		s.wg.Add(1)
		go func(c net.Conn) {
			defer s.wg.Done()
			defer func() {
				s.mutex.Lock()
				s.logger.Printf("Closing connection from %s", c.RemoteAddr())
				delete(s.connections, c)
				s.mutex.Unlock()
				c.Close()
			}()

			s.handleConnection(c)
		}(conn)
	}
}

// handleConnection processes a single client connection
func (s *Server) handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)

	// Keep handling messages until connection is closed
	for {
		// Read the message length
		lengthStr, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				s.logger.Printf("Error reading message length: %v", err)
			}
			// EOF or other error, connection is closed or broken
			return
		}

		// Parse the length
		lengthStr = strings.TrimSpace(lengthStr)
		if lengthStr == "" {
			// Skip empty lines
			continue
		}

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			s.logger.Printf("Invalid message length: %v", err)
			return
		}

		// Read the XML data
		xmlData := make([]byte, length)
		_, err = io.ReadFull(reader, xmlData)
		if err != nil {
			s.logger.Printf("Error reading XML data: %v", err)
			return
		}

		// Parse the XML using xmlparser
		parser := &xmlparser.Xmlparser{}
		parsedXML, xmlType, err := parser.Parse(xmlData)
		if err != nil {
			s.logger.Printf("Error parsing XML: %v", err)
			continue // Try to read the next message instead of closing the connection
		}

		// Process based on the type of XML structure
		var response []byte
		switch xmlType.Name() {
		case "Create":
			createData, ok := parsedXML.(xmlparser.Create)
			if !ok {
				s.logger.Printf("Error: Failed to cast to Create type")
				continue // Try to read the next message
			}
			response, err = s.handleCreate(createData)
		case "Transaction":
			transactionData, ok := parsedXML.(xmlparser.Transaction)
			if !ok {
				s.logger.Printf("Error: Failed to cast to Transaction type")
				continue // Try to read the next message
			}
			response, err = s.handleTransactions(transactionData)
		default:
			s.logger.Printf("Unknown XML type: %s", xmlType.Name())
			continue // Try to read the next message
		}

		if err != nil {
			s.logger.Printf("Error processing request: %v", err)
			continue // Try to read the next message
		}

		// Send the response
		respStr := string(response)
		_, err = conn.Write([]byte(fmt.Sprintf("%d\n%s", len(respStr), respStr)))
		if err != nil {
			s.logger.Printf("Error sending response: %v", err)
			return
		}
	}
}

// generateOrderID creates a unique order ID
func (s *Server) generateOrderID() string {
	s.idMutex.Lock()
	defer s.idMutex.Unlock()

	id := s.nextOrderID
	s.nextOrderID++

	return strconv.Itoa(id)
}

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			return fmt.Errorf("failed to close listener: %v", err)
		}
	}

	// Close all existing connections
	s.mutex.Lock()
	for conn := range s.connections {
		conn.Close()
	}
	s.mutex.Unlock()

	// Wait for all connection handlers to finish
	s.wg.Wait()
	return nil
}
