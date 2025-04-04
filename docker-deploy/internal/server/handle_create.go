package server

import (
	"StockOverflow/internal/database"
	"StockOverflow/internal/pool"
	"StockOverflow/pkg/xmlparser"
	"StockOverflow/pkg/xmlresponse"
	"encoding/xml"
	"fmt"
)

// handleCreate processes a create XML request and returns the response
func (s *Server) handleCreate(createData xmlparser.Create) ([]byte, error) {
	// Initialize response
	response := xmlresponse.Results{}

	// process children in order
	for _, child := range createData.Children {
		switch ele := child.(type) {
		case xmlparser.Account:
			s.processAccount(&ele, &response)
		case xmlparser.Symbol:
			s.processSymbol(&ele, &response)
		}
	}

	// Marshal response to XML
	xmlHeader := []byte(xml.Header)
	xmlBody, err := xml.MarshalIndent(response, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %v", err)
	}

	return append(xmlHeader, xmlBody...), nil
}

// process account ele
func (s *Server) processAccount(account *xmlparser.Account, response *xmlresponse.Results) {

	// Process accounts first
	s.logger.Printf("Processing create account request for ID: %s", account.ID)

	// Store in database
	err := database.CreateAccount(s.db, account.ID, account.Balance)
	if err != nil {
		s.logger.Printf("Failed to create account %s: %v", account.ID, err)
		response.Errors = append(response.Errors, xmlresponse.Error{
			ID:      account.ID,
			Message: err.Error(),
		})
		return
	}

	// Store in server memory
	s.accountsMutex.Lock()
	s.accounts[account.ID] = &AccountNode{
		ID:      account.ID,
		Balance: account.Balance,
	}
	s.accountsMutex.Unlock()

	// Add success response
	response.Created = append(response.Created, xmlresponse.Created{
		ID: account.ID,
	})
	s.logger.Printf("Successfully created account: %s", account.ID)
}

// process symbol ele
func (s *Server) processSymbol(symbol *xmlparser.Symbol, response *xmlresponse.Results) {
	s.logger.Printf("Processing create symbol request for: %s", symbol.Symbol)

	// Get or create stock node for trading
	var stockNode *pool.LruNode[*pool.StockNode]
	node, err := s.stockPool.Get(symbol.Symbol)
	if err != nil {
		// Symbol doesn't exist in pool, create a new node
		stockNode = pool.NewStockNode(symbol.Symbol, 10)
		err = s.stockPool.Put(stockNode)
		if err != nil {
			s.logger.Printf("Warning: Failed to add stock node to pool: %v", err)
			// Non-fatal, continue processing
		}
	} else {
		stockNode = node
	}

	// Process allocations for this symbol
	for _, allocation := range symbol.Accounts {
		// Validate account exists
		s.accountsMutex.RLock()
		// TODO change the account map to lru accounts
		account, exists := s.accounts[allocation.ID]
		s.accountsMutex.RUnlock()

		if !exists {
			// Try to load from database
			dbAccount, err := database.GetAccount(s.db, allocation.ID)
			if err != nil {
				s.logger.Printf("Account not found for allocation: %s", allocation.ID)
				response.Errors = append(response.Errors, xmlresponse.Error{
					Symbol:  symbol.Symbol,
					ID:      allocation.ID,
					Message: "Account not found",
				})
				continue
			}

			// Create account in memory
			account = &AccountNode{
				ID:      dbAccount.ID,
				Balance: dbAccount.Balance,
			}

			s.accountsMutex.Lock()
			s.accounts[account.ID] = account
			s.accountsMutex.Unlock()
		}

		// Update position in database
		err = database.CreateOrUpdatePosition(s.db, allocation.ID, symbol.Symbol, allocation.Balance)
		if err != nil {
			s.logger.Printf("Failed to update position: %v", err)
			response.Errors = append(response.Errors, xmlresponse.Error{
				Symbol:  symbol.Symbol,
				ID:      allocation.ID,
				Message: fmt.Sprintf("Database error: %v", err),
			})
			continue
		}

		// Add success response
		response.Created = append(response.Created, xmlresponse.Created{
			Symbol: symbol.Symbol,
			ID:     allocation.ID,
		})
		s.logger.Printf("Successfully allocated %s shares of %s to account %s",
			allocation.Balance.String(), symbol.Symbol, allocation.ID)
	}
}
