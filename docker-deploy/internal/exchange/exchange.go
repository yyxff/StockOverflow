package exchange

import (
	"StockOverflow/internal/database"
	"StockOverflow/internal/pool"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/shopspring/decimal"
)

// Exchange represents the core matching engine
type Exchange struct {
	db        *sql.DB
	stockPool *pool.StockPool
	logger    *log.Logger
}

// NewExchange creates a new exchange instance
func NewExchange(db *sql.DB, stockPool *pool.StockPool, logger *log.Logger) *Exchange {
	return &Exchange{
		db:        db,
		stockPool: stockPool,
		logger:    logger,
	}
}

// PlaceOrder places a new order in the exchange
func (e *Exchange) PlaceOrder(orderID, accountID, symbol string, amount, price decimal.Decimal) error {
	// Create order in database
	now := time.Now().UnixNano()
	order := &database.Order{
		ID:        orderID,
		AccountID: accountID,
		Symbol:    symbol,
		Amount:    amount,
		Price:     price,
		Status:    "open",
		Remaining: amount.Abs(),
		Timestamp: now,
	}

	err := database.CreateOrder(e.db, order)
	if err != nil {
		return fmt.Errorf("failed to create order in database: %v", err)
	}

	// Call matching logic
	e.MatchOrder(orderID, accountID, symbol, amount.IsPositive(), price, amount.Abs())
	return nil
}

// MatchOrder handles the matching of an order with existing orders
func (e *Exchange) MatchOrder(orderID string, accountID string, symbol string, isBuy bool, price decimal.Decimal, amount decimal.Decimal) {
	// Try to get the stock node for this symbol
	stockNode, err := e.stockPool.Get(symbol)
	if err != nil {
		// Symbol doesn't exist in pool, create a new node
		stockNode = pool.NewStockNode(symbol, 10)

		buyers := stockNode.GetValue().GetBuyers()
		buyers.SetDB(e.db)
		buyers.CheckMin()
		sellers := stockNode.GetValue().GetSellers()
		sellers.SetDB(e.db)
		sellers.CheckMin()

		err = e.stockPool.Put(stockNode)
		if err != nil {
			e.logger.Printf("Warning: Failed to add stock node to pool: %v", err)
			return
		}
	}

	// Lock the stock node for matching
	stockNode.Lock()
	defer stockNode.Unlock()

	if isBuy {
		e.matchBuyOrder(stockNode, orderID, accountID, symbol, price, amount)
	} else {
		e.matchSellOrder(stockNode, orderID, accountID, symbol, price, amount)
	}
}

// matchBuyOrder handles matching a buy order with existing sell orders
func (e *Exchange) matchBuyOrder(stockNode *pool.LruNode[*pool.StockNode], orderID string, accountID string, symbol string, price decimal.Decimal, amount decimal.Decimal) {
	// This is a buy order, try to match with sell orders
	sellersHeap := stockNode.GetValue().GetSellers()

	// If sellers heap is empty, just add to buyers heap and return
	if sellersHeap.Len() == 0 {
		e.addRemainingBuyOrder(stockNode, orderID, price, amount)
		return
	}

	// Check if we can match with the best sell order
	sellOrderData, err := sellersHeap.SafePop()
	if err != nil {
		e.logger.Printf("Error popping from sellers heap: %v", err)
		e.addRemainingBuyOrder(stockNode, orderID, price, amount)
		return
	}

	sellOrderInfo := sellOrderData.(pool.Order)
	sellPrice := sellOrderInfo.GetPrice()

	// If the best sell price is higher than our buy price, no match is possible
	if sellPrice.GreaterThan(price) {
		// Put the order back in the heap
		sellOrderPtr := &sellOrderInfo
		sellersHeap.SafePush(sellOrderPtr)

		// Add our buy order to the buyers heap without attempting to match
		e.logger.Printf("No match possible: best sell price %s > buy price %s",
			sellPrice.String(), price.String()+" For Order ID: "+orderID)
		e.addRemainingBuyOrder(stockNode, orderID, price, amount)
		return
	}

	// Put the sell order back for the matching process
	sellOrderPtr := &sellOrderInfo
	sellersHeap.SafePush(sellOrderPtr)

	// Keep matching until no compatible sellers or order is fully executed
	remainingAmount := amount
	for remainingAmount.GreaterThan(decimal.Zero) && sellersHeap.Len() > 0 {
		// Get the best sell order (lowest price)
		sellOrderData, err := sellersHeap.SafePop()
		if err != nil {
			e.logger.Printf("Error popping from sellers heap: %v", err)
			break
		}

		sellOrderInfo := sellOrderData.(pool.Order)
		sellPrice := sellOrderInfo.GetPrice()

		// Check if prices are compatible
		if sellPrice.GreaterThan(price) {
			// No compatible price, put the order back
			sellOrderPtr := &sellOrderInfo
			sellersHeap.SafePush(sellOrderPtr)
			e.logger.Printf("No compatible price: sell price %s > buy price %s",
				sellPrice.String(), price.String()+" For Order ID: "+orderID)
			break
		}

		// Get the sell order from the database
		sellOrderID := sellOrderInfo.GetID()
		sellOrder, err := database.GetOrder(e.db, sellOrderID)

		if err != nil || sellOrder.Status != "open" {
			// Skip this order and continue
			e.logger.Printf("Invalid sell order: %v", err)
			continue
		}

		// Determine the execution price (price of the earlier order)
		var executionPrice decimal.Decimal

		if sellOrder.Timestamp < time.Now().UnixNano() {
			executionPrice = sellOrder.Price
		} else {
			executionPrice = price
		}
		var refundPrice decimal.Decimal
		if (price).GreaterThan(executionPrice) {
			refundPrice = price.Sub(executionPrice)
		} else {
			refundPrice = decimal.Zero
		}

		// Determine execution amount
		var executionAmount decimal.Decimal
		if remainingAmount.LessThanOrEqual(sellOrder.Remaining) {
			executionAmount = remainingAmount
		} else {
			executionAmount = sellOrder.Remaining
		}

		// Execute the match
		err = e.executeMatch(
			orderID, accountID, sellOrderID, sellOrder.AccountID,
			symbol, executionAmount, executionPrice, refundPrice, time.Now().UnixNano(),
		)

		if err != nil {
			e.logger.Printf("Error executing match: %v", err)
			continue
		}

		// Update remaining amount
		remainingAmount = remainingAmount.Sub(executionAmount)
	}

	// If order still has remaining amount, add to buyers heap
	if remainingAmount.GreaterThan(decimal.Zero) {
		e.addRemainingBuyOrder(stockNode, orderID, price, remainingAmount)
	}
}

// Helper function to add a buy order with remaining amount to the buyers heap
func (e *Exchange) addRemainingBuyOrder(stockNode *pool.LruNode[*pool.StockNode], orderID string, price decimal.Decimal, remainingAmount decimal.Decimal) {
	// Update remaining amount in database
	err := database.UpdateOrderStatus(e.db, orderID, "open", remainingAmount, 0)
	if err != nil {
		e.logger.Printf("Error updating order status: %v", err)
	}

	// Add to buyers heap
	buyerOrder := pool.NewOrder(
		orderID,
		uint(remainingAmount.IntPart()),
		price,
		time.Now(),
	)
	stockNode.GetValue().GetBuyers().SafePush(buyerOrder)
}

// matchSellOrder handles matching a sell order with existing buy orders
func (e *Exchange) matchSellOrder(stockNode *pool.LruNode[*pool.StockNode], orderID string, accountID string, symbol string, price decimal.Decimal, amount decimal.Decimal) {
	// This is a sell order, try to match with buy orders
	buyersHeap := stockNode.GetValue().GetBuyers()

	// If buyers heap is empty, just add to sellers heap and return
	if buyersHeap.Len() == 0 {
		e.addRemainingSellOrder(stockNode, orderID, price, amount)
		return
	}

	// Check if we can match with the best buy order
	buyOrderData, err := buyersHeap.SafePop()
	if err != nil {
		e.logger.Printf("Error popping from buyers heap: %v", err)
		e.addRemainingSellOrder(stockNode, orderID, price, amount)
		return
	}

	buyOrderInfo := buyOrderData.(pool.Order)
	buyPrice := buyOrderInfo.GetPrice()

	// If the best buy price is lower than our sell price, no match is possible
	if buyPrice.LessThan(price) {
		// Put the order back in the heap
		buyOrderPtr := &buyOrderInfo
		buyersHeap.SafePush(buyOrderPtr)

		// Add our sell order to the sellers heap without attempting to match
		e.logger.Printf("No match possible: best buy price %s < sell price %s",
			buyPrice.String(), price.String())
		e.addRemainingSellOrder(stockNode, orderID, price, amount)
		return
	}

	// Put the buy order back for the matching process
	buyOrderPtr := &buyOrderInfo
	buyersHeap.SafePush(buyOrderPtr)

	// Keep matching until no compatible buyers or order is fully executed
	remainingAmount := amount
	for remainingAmount.GreaterThan(decimal.Zero) && buyersHeap.Len() > 0 {
		// Get the best buy order (highest price)
		buyOrderData, err := buyersHeap.SafePop()
		if err != nil {
			e.logger.Printf("Error popping from buyers heap: %v", err)
			break
		}

		buyOrderInfo := buyOrderData.(pool.Order)
		buyPrice := buyOrderInfo.GetPrice()

		// Check if prices are compatible
		if buyPrice.LessThan(price) {
			// No compatible price, put the order back
			buyOrderPtr := &buyOrderInfo
			buyersHeap.SafePush(buyOrderPtr)
			e.logger.Printf("No compatible price: buy price %s < sell price %s",
				buyPrice.String(), price.String())
			break
		}

		// Get the buy order from the database
		buyOrderID := buyOrderInfo.GetID()
		buyOrder, err := database.GetOrder(e.db, buyOrderID)

		if err != nil || buyOrder.Status != "open" {
			// Skip this order and continue
			e.logger.Printf("Invalid buy order: %v", err)
			continue
		}

		// Determine the execution price (price of the earlier order)
		var executionPrice decimal.Decimal

		if buyOrder.Timestamp < time.Now().UnixNano() {
			executionPrice = buyOrder.Price
		} else {
			executionPrice = price
		}
		// Determine refund price
		var refundPrice decimal.Decimal
		if (buyOrder.Price).GreaterThan(executionPrice) {
			refundPrice = buyOrder.Price.Sub(executionPrice)
		} else {
			refundPrice = decimal.Zero
		}
		// Determine execution amount
		var executionAmount decimal.Decimal
		if remainingAmount.LessThanOrEqual(buyOrder.Remaining) {
			executionAmount = remainingAmount
		} else {
			executionAmount = buyOrder.Remaining
		}

		// Execute the match
		err = e.executeMatch(
			buyOrderID, buyOrder.AccountID, orderID, accountID,
			symbol, executionAmount, executionPrice, refundPrice, time.Now().UnixNano(),
		)

		if err != nil {
			e.logger.Printf("Error executing match: %v", err)
			continue
		}

		// Update remaining amount
		remainingAmount = remainingAmount.Sub(executionAmount)
	}

	// If order still has remaining amount, add to sellers heap
	if remainingAmount.GreaterThan(decimal.Zero) {
		e.addRemainingSellOrder(stockNode, orderID, price, remainingAmount)
	}
}

// Helper function to add a sell order with remaining amount to the sellers heap
func (e *Exchange) addRemainingSellOrder(stockNode *pool.LruNode[*pool.StockNode], orderID string, price decimal.Decimal, remainingAmount decimal.Decimal) {
	// Update remaining amount in database
	err := database.UpdateOrderStatus(e.db, orderID, "open", remainingAmount, 0)
	if err != nil {
		e.logger.Printf("Error updating order status: %v", err)
	}

	// Add to sellers heap
	sellerOrder := pool.NewOrder(
		orderID,
		uint(remainingAmount.IntPart()),
		price,
		time.Now(),
	)
	stockNode.GetValue().GetSellers().SafePush(sellerOrder)
}

// executeMatch executes a trade between a buy order and a sell order
func (e *Exchange) executeMatch(buyOrderID, buyerAccountID, sellOrderID, sellerAccountID,
	symbol string, amount, executionPrice decimal.Decimal, refundPrice decimal.Decimal, timestamp int64) error {
	// Execute the match within a transaction to ensure atomicity
	return database.ExecuteWithTransaction(e.db, func(tx *sql.Tx) error {
		txFuncs := &database.CommonTxFunctions{Tx: tx}

		// 1. Get current orders
		buyOrder, err := database.GetOrder(e.db, buyOrderID)
		if err != nil {
			return fmt.Errorf("failed to get buy order: %v", err)
		}

		sellOrder, err := database.GetOrder(e.db, sellOrderID)
		if err != nil {
			return fmt.Errorf("failed to get sell order: %v", err)
		}

		// 2. Record execution for both orders
		err = txFuncs.RecordExecution(buyOrderID, amount, executionPrice, timestamp)
		if err != nil {
			return fmt.Errorf("failed to record buy execution: %v", err)
		}

		err = txFuncs.RecordExecution(sellOrderID, amount, executionPrice, timestamp)
		if err != nil {
			return fmt.Errorf("failed to record sell execution: %v", err)
		}

		// 3. Update remaining amounts for both orders
		newBuyRemaining := buyOrder.Remaining.Sub(amount)
		buyStatus := "open"
		if newBuyRemaining.IsZero() {
			buyStatus = "executed"
		}
		err = txFuncs.UpdateOrderStatus(buyOrderID, buyStatus, newBuyRemaining, 0)
		if err != nil {
			return fmt.Errorf("failed to update buy order status: %v", err)
		}

		newSellRemaining := sellOrder.Remaining.Sub(amount)
		sellStatus := "open"
		if newSellRemaining.IsZero() {
			sellStatus = "executed"
		}
		err = txFuncs.UpdateOrderStatus(sellOrderID, sellStatus, newSellRemaining, 0)
		if err != nil {
			return fmt.Errorf("failed to update sell order status: %v", err)
		}

		// 4. Process refund for buyer if applicable
		if refundPrice.GreaterThan(decimal.Zero) {
			refundAmount := amount.Mul(refundPrice)

			// Get buyer account
			buyerAccount, err := database.GetAccount(e.db, buyerAccountID)
			if err != nil {
				return fmt.Errorf("failed to get buyer account: %v", err)
			}

			// Update buyer's balance with refund
			newBuyerBalance := buyerAccount.Balance.Add(refundAmount)
			err = txFuncs.UpdateAccountBalance(buyerAccountID, newBuyerBalance)
			if err != nil {
				return fmt.Errorf("failed to update buyer balance with refund: %v", err)
			}

			e.logger.Printf("Buyer %s gets refund: %s for order %s",
				buyerAccountID, refundAmount.String(), buyOrderID)
		}

		// 5. Update seller's balance
		sellerAccount, err := database.GetAccount(e.db, sellerAccountID)
		if err != nil {
			return fmt.Errorf("failed to get seller account: %v", err)
		}

		// Calculate trade amount based on execution price
		tradeAmount := amount.Mul(executionPrice)
		newSellerBalance := sellerAccount.Balance.Add(tradeAmount)
		err = txFuncs.UpdateAccountBalance(sellerAccountID, newSellerBalance)
		if err != nil {
			return fmt.Errorf("failed to update seller balance: %v", err)
		}

		// 6. Update buyer's position
		buyerPositions, err := database.GetPositions(e.db, buyerAccountID)
		if err != nil {
			return fmt.Errorf("failed to get buyer positions: %v", err)
		}

		// Find the position for this symbol
		var currentBuyerPosition *database.Position
		for i := range buyerPositions {
			if buyerPositions[i].Symbol == symbol {
				currentBuyerPosition = &buyerPositions[i]
				break
			}
		}

		// Update or create buyer's position
		if currentBuyerPosition == nil {
			// Create new position
			err = txFuncs.CreateOrUpdatePosition(buyerAccountID, symbol, amount)
		} else {
			// Update existing position
			newAmount := currentBuyerPosition.Amount.Add(amount)
			err = txFuncs.CreateOrUpdatePosition(buyerAccountID, symbol, newAmount)
		}

		if err != nil {
			return fmt.Errorf("failed to update buyer position: %v", err)
		}

		// Log successful execution
		e.logger.Printf("Executed match: %s bought %s %s from %s at %s",
			buyerAccountID, amount.String(), symbol, sellerAccountID, executionPrice.String())

		return nil
	})
}

// CancelOrder cancels an open order
func (e *Exchange) CancelOrder(orderID string) error {
	// Get order from database
	order, err := database.GetOrder(e.db, orderID)
	if err != nil {
		return fmt.Errorf("failed to find order: %v", err)
	}

	// Check if order is already completed or canceled
	if order.Status != "open" {
		return fmt.Errorf("order is not open")
	}

	// Get the current timestamp
	now := time.Now().UnixNano()

	// Return funds or shares based on order type
	isBuy := order.Amount.IsPositive()

	return database.ExecuteWithTransaction(e.db, func(tx *sql.Tx) error {
		txFuncs := &database.CommonTxFunctions{Tx: tx}

		// Update order status
		err := txFuncs.UpdateOrderStatus(orderID, "canceled", order.Remaining, now)
		if err != nil {
			return fmt.Errorf("failed to update order status: %v", err)
		}

		// Return funds or shares
		if isBuy {
			// Get latest account from DB to ensure accurate balance
			dbAccount, err := database.GetAccount(e.db, order.AccountID)
			if err != nil {
				return fmt.Errorf("failed to get account: %v", err)
			}

			// Calculate refund amount
			refundAmount := order.Remaining.Mul(order.Price)
			newBalance := dbAccount.Balance.Add(refundAmount)

			// Update account balance
			err = txFuncs.UpdateAccountBalance(order.AccountID, newBalance)
			if err != nil {
				return fmt.Errorf("failed to update account balance: %v", err)
			}
		} else {
			// Get the position
			positions, err := database.GetPositions(e.db, order.AccountID)
			if err != nil {
				return fmt.Errorf("failed to get positions: %v", err)
			}

			// Find the position for this symbol
			var currentPosition *database.Position
			for i := range positions {
				if positions[i].Symbol == order.Symbol {
					currentPosition = &positions[i]
					break
				}
			}

			if currentPosition == nil {
				// Create new position
				err = txFuncs.CreateOrUpdatePosition(order.AccountID, order.Symbol, order.Remaining)
			} else {
				// Update existing position
				newAmount := currentPosition.Amount.Add(order.Remaining)
				err = txFuncs.CreateOrUpdatePosition(order.AccountID, order.Symbol, newAmount)
			}

			if err != nil {
				return fmt.Errorf("failed to update position: %v", err)
			}
		}

		return nil
	})
}

// GetOrderStatus returns the current status of an order
func (e *Exchange) GetOrderStatus(orderID string) (*database.Order, []database.Execution, error) {
	// Get order from database
	order, err := database.GetOrder(e.db, orderID)
	if err != nil {
		return nil, nil, err
	}

	// Get executions for this order
	executions, err := database.GetOrderExecutions(e.db, orderID)
	if err != nil {
		// Return the order with empty executions
		return order, []database.Execution{}, nil
	}

	return order, executions, nil
}
