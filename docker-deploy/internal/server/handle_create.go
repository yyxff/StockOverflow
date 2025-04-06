package server

import (
	"StockOverflow/internal/database"
	"StockOverflow/internal/pool"
	"StockOverflow/pkg/xmlparser"
	"StockOverflow/pkg/xmlresponse"
	"fmt"

	"github.com/shopspring/decimal"
)

// handleCreate processes a create XML request and returns the response
func (s *Server) handleCreate(createData xmlparser.Create) ([]byte, error) {
	// Initialize response
	response := xmlresponse.Results{
		Children: make([]any, 0),
	}

	// process children in order
	for _, child := range createData.Children {
		switch ele := child.(type) {
		case xmlparser.Account:
			s.processAccount(&ele, &response)
		case xmlparser.Symbol:
			s.processSymbol(&ele, &response)
		}
	}

	return marshalResponse(response)
}

// process account ele
func (s *Server) processAccount(account *xmlparser.Account, response *xmlresponse.Results) {

	// Process accounts first
	s.logger.Printf("Processing create account request for ID: %s", account.ID)

	// Store in database
	err := database.CreateAccount(s.db, account.ID, account.Balance)
	if err != nil {
		s.logger.Printf("Failed to create account %s: %v", account.ID, err)
		response.Children = append(response.Children, xmlresponse.Error{
			ID:      account.ID,
			Message: err.Error(),
		})
		return
	}

	// Store in server memory
	s.accountsMutex.Lock()
	s.accounts[account.ID] = &AccountNode{
		ID:        account.ID,
		Balance:   account.Balance,
		Positions: make(map[string]decimal.Decimal),
	}
	s.accountsMutex.Unlock()

	// Add success response
	response.Children = append(response.Children, xmlresponse.Created{
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
		account, exists := s.accounts[allocation.ID]
		s.accountsMutex.RUnlock()

		if !exists {
			// Try to load from database
			dbAccount, err := database.GetAccount(s.db, allocation.ID)
			if err != nil {
				s.logger.Printf("Account not found for allocation: %s", allocation.ID)
				response.Children = append(response.Children, xmlresponse.Error{
					Symbol:  symbol.Symbol,
					ID:      allocation.ID,
					Message: "Account not found",
				})
				continue
			}

			// Create account in memory with Positions map
			account = &AccountNode{
				ID:        dbAccount.ID,
				Balance:   dbAccount.Balance,
				Positions: make(map[string]decimal.Decimal),
			}

			// Load all existing positions for this account
			positions, err := database.GetPositions(s.db, allocation.ID)
			if err == nil {
				for _, pos := range positions {
					account.Positions[pos.Symbol] = pos.Amount
				}
			}

			s.accountsMutex.Lock()
			s.accounts[account.ID] = account
			s.accountsMutex.Unlock()
		}

		// Check if position already exists in database
		position, err := database.GetPosition(s.db, allocation.ID, symbol.Symbol)

		var newAmount decimal.Decimal
		if err == nil && position != nil {
			// Position exists, add to it
			newAmount = position.Amount.Add(allocation.Amount)
		} else {
			// Position doesn't exist, create with allocation amount
			newAmount = allocation.Amount
		}

		// Update position in database with the new amount
		err = database.CreateOrUpdatePosition(s.db, allocation.ID, symbol.Symbol, newAmount)
		if err != nil {
			s.logger.Printf("Failed to update position: %v", err)
			response.Children = append(response.Children, xmlresponse.Error{
				Symbol:  symbol.Symbol,
				ID:      allocation.ID,
				Message: fmt.Sprintf("Database error: %v", err),
			})
			continue
		}

		// Update position in memory
		s.accountsMutex.Lock()
		if s.accounts[allocation.ID].Positions == nil {
			s.accounts[allocation.ID].Positions = make(map[string]decimal.Decimal)
		}

		// Add to existing position in memory, if any
		currentAmount := s.accounts[allocation.ID].Positions[symbol.Symbol] // if just make, it's zero.
		s.accounts[allocation.ID].Positions[symbol.Symbol] = currentAmount.Add(allocation.Amount)
		s.accountsMutex.Unlock()

		// Add success response
		response.Children = append(response.Children, xmlresponse.Created{
			Symbol: symbol.Symbol,
			ID:     allocation.ID,
		})
		s.logger.Printf("Successfully allocated %s shares of %s to account %s (total now: %s)",
			allocation.Amount.String(), symbol.Symbol, allocation.ID, newAmount.String())
	}
}
