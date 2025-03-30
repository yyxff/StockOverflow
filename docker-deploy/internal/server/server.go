package server

import (
	"StockOverflow/pkg/xmlparser"
	"bufio"
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

// Server represents the exchange server
type Server struct {
	listener    net.Listener
	logger      *log.Logger
	wg          sync.WaitGroup
	connections map[net.Conn]struct{}
	mutex       sync.Mutex
	// Exchange state
	accounts      map[string]*xmlparser.Account
	symbols       map[string]*xmlparser.Symbol
	orders        map[string]*xmlparser.Order
	nextOrderID   int
	accountsMutex sync.RWMutex
	symbolsMutex  sync.RWMutex
	ordersMutex   sync.RWMutex
}

// NewServer creates a new exchange server
func NewServer(logger *log.Logger) *Server {
	return &Server{
		logger:      logger,
		connections: make(map[net.Conn]struct{}),
		accounts:    make(map[string]*xmlparser.Account),
		symbols:     make(map[string]*xmlparser.Symbol),
		orders:      make(map[string]*xmlparser.Order),
		nextOrderID: 1,
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

	// Read the message length
	lengthStr, err := reader.ReadString('\n')
	if err != nil {
		s.logger.Printf("Error reading message length: %v", err)
		return
	}

	// Parse the length
	lengthStr = strings.TrimSpace(lengthStr)
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		s.logger.Printf("Invalid message length: %v", err)
		return
	}

	// Read the XML data
	xmlData := make([]byte, length)
	_, err = reader.Read(xmlData)
	if err != nil {
		s.logger.Printf("Error reading XML data: %v", err)
		return
	}

	// Process the XML and get a response
	response, err := s.processXML(xmlData)
	if err != nil {
		s.logger.Printf("Error processing XML: %v", err)
		return
	}

	// Send the response
	respStr := string(response)
	_, err = conn.Write([]byte(fmt.Sprintf("%d\n%s", len(respStr), respStr)))
	if err != nil {
		s.logger.Printf("Error sending response: %v", err)
		return
	}
}

// processXML handles the parsing and execution of XML commands
func (s *Server) processXML(data []byte) ([]byte, error) {
	// Determine the root element to decide how to process
	var root struct {
		XMLName xml.Name
	}

	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %v", err)
	}

	// Process based on the root element
	switch root.XMLName.Local {
	case "create":
		return s.handleCreate(data)
	case "transactions":
		return s.handleTransactions(data)
	default:
		return nil, fmt.Errorf("unknown XML root element: %s", root.XMLName.Local)
	}
}

// handleCreate processes a create request (placeholder implementation)
func (s *Server) handleCreate(data []byte) ([]byte, error) {
	// Placeholder implementation - will be filled with actual logic
	result := `<?xml version="1.0" encoding="UTF-8"?>
<results>
  <created id="test"/>
</results>`
	return []byte(result), nil
}

// handleTransactions processes a transactions request (placeholder implementation)
func (s *Server) handleTransactions(data []byte) ([]byte, error) {
	// Placeholder implementation - will be filled with actual logic
	result := `<?xml version="1.0" encoding="UTF-8"?>
<results>
  <opened sym="SPY" amount="100" limit="145.67" id="1"/>
</results>`
	return []byte(result), nil
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
