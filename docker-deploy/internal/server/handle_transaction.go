package server

import (
	"StockOverflow/internal/database"
	"StockOverflow/pkg/xmlparser"
	"StockOverflow/pkg/xmlresponse"
	"encoding/xml"
	"fmt"
	"reflect"

	"github.com/shopspring/decimal"
)

// handleTransactions processes a transactions XML request and returns the response
func (s *Server) handleTransactions(transactionData xmlparser.Transaction) ([]byte, error) {
	// Initialize response
	response := xmlresponse.Results{
		Children: make([]any, 0),
	}

	// Validate account exists
	s.accountsMutex.RLock()
	account, exists := s.accounts[transactionData.ID]
	s.accountsMutex.RUnlock()

	if !exists {
		// Try to load from database
		// TODO account's lru logic
		dbAccount, err := database.GetAccount(s.db, transactionData.ID)
		if err != nil {
			// Account not found, return error for all transactions
			s.logger.Printf("Account not found for transactions: %s", transactionData.ID)

			// Generate errors for all operations
			generateAccountNotFoundErrors(&response, transactionData)

			// Return response since account doesn't exist
			return marshalResponse(response)
		}

		// Create account in memory with positions
		account = &AccountNode{
			ID:        dbAccount.ID,
			Balance:   dbAccount.Balance,
			Positions: make(map[string]decimal.Decimal),
		}

		// Load all existing positions for this account
		positions, err := database.GetPositions(s.db, dbAccount.ID)
		if err == nil {
			for _, pos := range positions {
				account.Positions[pos.Symbol] = pos.Amount
			}
		} else {
			s.logger.Printf("Warning: Failed to load positions for account %s: %v", dbAccount.ID, err)
			// Continue with empty positions map
		}

		// Store in server memory
		s.accountsMutex.Lock()
		s.accounts[account.ID] = account
		s.accountsMutex.Unlock()
	}

	// process ele in order
	for _, child := range transactionData.Children {
		switch ele := child.(type) {
		case xmlparser.Order:
			s.processOrder(&ele, account, &response)
		case xmlparser.Query:
			s.processQuery(&ele, &response)
		case xmlparser.Cancel:
			s.processCancel(&ele, &response)
		default:
			s.logger.Fatalf("unknown type in children: %T", reflect.TypeOf(ele))
		}

	}

	// Marshal response to XML
	return marshalResponse(response)
}

// validateAndReserve validates an order and reserves the necessary funds or shares
func (s *Server) validateAndReserve(account *AccountNode, symbol string, amount, price decimal.Decimal, isBuy bool) string {
	if isBuy {
		// For buy order, check account balance
		totalCost := amount.Mul(price)

		// Get the latest account balance from memory
		s.accountsMutex.RLock()
		accountBalance := s.accounts[account.ID].Balance
		s.accountsMutex.RUnlock()

		if accountBalance.LessThan(totalCost) {
			return "Insufficient funds for account: " + account.ID
		}

		// Reserve the funds by updating account balance
		newBalance := accountBalance.Sub(totalCost)

		err := database.UpdateAccountBalance(s.db, account.ID, newBalance)
		if err != nil {
			return fmt.Sprintf("Failed to update account balance: %v", err)
		}

		// Update the account in memory
		s.accountsMutex.Lock()
		s.accounts[account.ID].Balance = newBalance
		s.accountsMutex.Unlock()

		return ""
	} else {
		// For sell order
		sellAmount := amount.Abs()

		// Get the position from memory instead of database
		s.accountsMutex.RLock()
		currentPosition, exists := s.accounts[account.ID].Positions[symbol]
		s.accountsMutex.RUnlock()

		if !exists || currentPosition.LessThan(sellAmount) {
			posAmount := decimal.Zero
			if exists {
				posAmount = currentPosition
			}
			return "Insufficient shares for: " + posAmount.String() + " in account: " + account.ID
		}

		// Calculate new position amount
		newAmount := currentPosition.Sub(sellAmount)
		// Update position in database
		err := database.CreateOrUpdatePosition(s.db, account.ID, symbol, newAmount)
		if err != nil {
			return fmt.Sprintf("Failed to update position: %v", err)
		}

		// Update position in memory
		s.accountsMutex.Lock()
		s.accounts[account.ID].Positions[symbol] = newAmount
		s.accountsMutex.Unlock()
		return ""
	}
}

// createStatusResponse creates a status response from an order and its executions
func createStatusResponse(orderID string, order *database.Order, executions []database.Execution) xmlresponse.Status {
	status := xmlresponse.Status{
		ID: orderID,
	}

	// Add appropriate status details
	switch order.Status {
	case "open":
		status.Open = []xmlresponse.Open{
			{Shares: float64(order.Remaining.InexactFloat64())},
		}
	case "canceled":
		status.Canceled = []xmlresponse.Canceled{
			{Shares: float64(order.Remaining.InexactFloat64()), Time: order.CanceledTime},
		}
	}

	// Add executions
	for _, exec := range executions {
		status.Executed = append(status.Executed, xmlresponse.Executed{
			Shares: float64(exec.Shares.InexactFloat64()),
			Price:  float64(exec.Price.InexactFloat64()),
			Time:   exec.Timestamp,
		})
	}

	return status
}

// createCanceledResponse creates a canceled response from an order and its executions
func createCanceledResponse(orderID string, order *database.Order, executions []database.Execution) xmlresponse.Canceled {
	canceled := xmlresponse.Canceled{
		ID:     orderID,
		Shares: float64(order.Remaining.InexactFloat64()),
		Time:   order.CanceledTime,
	}

	// Add executions
	for _, exec := range executions {
		canceled.Executed = append(canceled.Executed, xmlresponse.Executed{
			Shares: float64(exec.Shares.InexactFloat64()),
			Price:  float64(exec.Price.InexactFloat64()),
			Time:   exec.Timestamp,
		})
	}

	return canceled
}

// generateAccountNotFoundErrors generates errors for all operations when account doesn't exist
func generateAccountNotFoundErrors(response *xmlresponse.Results, transaction xmlparser.Transaction) {

	for _, child := range transaction.Children {
		switch ele := child.(type) {
		case xmlparser.Order:
			{
				response.Children = append(response.Children, xmlresponse.Error{
					Symbol:  ele.Symbol,
					Amount:  float64(ele.Amount),
					Limit:   float64(ele.LimitPrice.InexactFloat64()),
					Message: "Account not found",
				})
			}
		case xmlparser.Query:
			{
				response.Children = append(response.Children, xmlresponse.Error{
					ID:      ele.ID,
					Message: "Account not found",
				})
			}
		case xmlparser.Cancel:
			{
				response.Children = append(response.Children, xmlresponse.Error{
					ID:      ele.ID,
					Message: "Account not found",
				})
			}
		}
	}

}

// marshalResponse marshals a response to XML
func marshalResponse(response xmlresponse.Results) ([]byte, error) {
	xmlHeader := []byte(xml.Header)
	xmlBody, err := xml.MarshalIndent(response, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %v", err)
	}

	return append(xmlHeader, xmlBody...), nil
}

func (s *Server) processOrder(orderRequest *xmlparser.Order, account *AccountNode, response *xmlresponse.Results) {
	// Generate order ID
	orderID := s.generateOrderID()
	s.logger.Printf("Processing order: %s, symbol: %s, amount: %d, price: %s",
		orderID, orderRequest.Symbol, orderRequest.Amount, orderRequest.LimitPrice.String())

	// Convert amount to decimal for calculations
	amount := decimal.NewFromInt(int64(orderRequest.Amount))

	// Negative amount means sell, positive means buy
	isBuy := orderRequest.Amount > 0

	// Validate and reserve funds/shares
	errorMsg := s.validateAndReserve(account, orderRequest.Symbol, amount, orderRequest.LimitPrice, isBuy)

	// If there was an error, add it to response and continue
	if errorMsg != "" {
		response.Children = append(response.Children, xmlresponse.Error{
			Symbol:  orderRequest.Symbol,
			Amount:  float64(orderRequest.Amount),
			Limit:   float64(orderRequest.LimitPrice.InexactFloat64()),
			Message: errorMsg,
		})
		return
	}

	// Place the order in the exchange
	err := s.exchange.PlaceOrder(orderID, account.ID, orderRequest.Symbol, amount, orderRequest.LimitPrice)
	if err != nil {
		s.logger.Printf("Failed to place order: %v", err)
		response.Children = append(response.Children, xmlresponse.Error{
			Symbol:  orderRequest.Symbol,
			Amount:  float64(orderRequest.Amount),
			Limit:   float64(orderRequest.LimitPrice.InexactFloat64()),
			Message: fmt.Sprintf("Exchange error: %v", err),
		})
		return
	}

	// Add success response
	response.Children = append(response.Children, xmlresponse.Opened{
		Symbol: orderRequest.Symbol,
		Amount: float64(orderRequest.Amount),
		Limit:  float64(orderRequest.LimitPrice.InexactFloat64()),
		ID:     orderID,
	})

	s.logger.Printf("Successfully created order %s for %s %s at %s",
		orderID, amount.String(), orderRequest.Symbol, orderRequest.LimitPrice.String())
}

func (s *Server) processQuery(query *xmlparser.Query, response *xmlresponse.Results) {
	s.logger.Printf("Processing query for order: %s", query.ID)

	// Get order status from exchange
	order, executions, err := s.exchange.GetOrderStatus(query.ID)
	if err != nil {
		s.logger.Printf("Failed to get order status: %v", err)
		response.Children = append(response.Children, xmlresponse.Error{
			ID:      query.ID,
			Message: err.Error(),
		})
		return
	}

	// Convert to response format
	status := createStatusResponse(query.ID, order, executions)

	// Add to response
	response.Children = append(response.Children, status)
	s.logger.Printf("Processed query for order %s with status %s", query.ID, order.Status)
}

func (s *Server) processCancel(cancel *xmlparser.Cancel, response *xmlresponse.Results) {
	s.logger.Printf("Processing cancel for order: %s", cancel.ID)

	// Cancel the order in the exchange
	err := s.exchange.CancelOrder(cancel.ID)
	if err != nil {
		s.logger.Printf("Failed to cancel order: %v", err)
		response.Children = append(response.Children, xmlresponse.Error{
			ID:      cancel.ID,
			Message: err.Error(),
		})
		return
	}

	// Get updated order status
	order, executions, err := s.exchange.GetOrderStatus(cancel.ID)
	if err != nil {
		s.logger.Printf("Failed to get order status after cancel: %v", err)
		response.Children = append(response.Children, xmlresponse.Error{
			ID:      cancel.ID,
			Message: fmt.Sprintf("Order was canceled but error retrieving status: %v", err),
		})
		return
	}

	// Create canceled response
	canceled := createCanceledResponse(cancel.ID, order, executions)

	// Add to response
	response.Children = append(response.Children, canceled)
	s.logger.Printf("Successfully canceled order %s", cancel.ID)

	// if it's a buy order, refund the money, if it's a sell order, return the shares
	if order.Amount.IsPositive() { // Buy order
		// Refund money to the account
		s.accountsMutex.Lock()
		if account, exists := s.accounts[order.AccountID]; exists {
			refundAmount := order.Remaining.Mul(order.Price)
			account.Balance = account.Balance.Add(refundAmount)
			s.logger.Printf("Updated in-memory balance for account %s after buy order cancelation", order.AccountID)
		}
		s.accountsMutex.Unlock()
	} else { // Sell order
		// Return shares to the account
		s.accountsMutex.Lock()
		if account, exists := s.accounts[order.AccountID]; exists {
			if account.Positions == nil {
				account.Positions = make(map[string]decimal.Decimal)
			}
			currentShares := account.Positions[order.Symbol]
			account.Positions[order.Symbol] = currentShares.Add(order.Remaining)
			s.logger.Printf("Updated in-memory position for account %s after sell order cancelation", order.AccountID)
		}
		s.accountsMutex.Unlock()
	}
}
