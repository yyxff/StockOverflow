package exchange

import (
	"StockOverflow/internal/database"
	"StockOverflow/internal/pool"
	"database/sql"
	"fmt"
	"log"
	"strconv"
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
	now := time.Now().Unix()
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
		stockNode = pool.NewStockNode(symbol)
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
		// This is a buy order, try to match with sell orders
		sellersHeap := stockNode.GetValue().GetSellers()

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
					sellPrice.String(), price.String())
				break
			}

			// Get the sell order from the database
			sellOrderID := strconv.Itoa(int(sellOrderInfo.GetID()))
			sellOrder, err := database.GetOrder(e.db, sellOrderID)

			if err != nil || sellOrder.Status != "open" {
				// Skip this order and continue
				e.logger.Printf("Invalid sell order: %v", err)
				continue
			}

			// Determine the execution price (price of the earlier order)
			var executionPrice decimal.Decimal
			if sellOrder.Timestamp < time.Now().Unix() {
				executionPrice = sellOrder.Price
			} else {
				executionPrice = price
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
				symbol, executionAmount, executionPrice, time.Now().Unix(),
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
			// Update remaining amount in database
			err := database.UpdateOrderStatus(e.db, orderID, "open", remainingAmount, 0)
			if err != nil {
				e.logger.Printf("Error updating order status: %v", err)
			}

			// Add to buyers heap
			orderIDInt, _ := strconv.Atoi(orderID)
			buyerOrder := pool.NewOrder(
				uint(orderIDInt),
				uint(remainingAmount.IntPart()),
				price,
				time.Now(),
			)
			stockNode.GetValue().GetBuyers().SafePush(buyerOrder)
		}
	} else {
		// This is a sell order, try to match with buy orders
		buyersHeap := stockNode.GetValue().GetBuyers()

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
			buyOrderID := strconv.Itoa(int(buyOrderInfo.GetID()))
			buyOrder, err := database.GetOrder(e.db, buyOrderID)

			if err != nil || buyOrder.Status != "open" {
				// Skip this order and continue
				e.logger.Printf("Invalid buy order: %v", err)
				continue
			}

			// Determine the execution price (price of the earlier order)
			var executionPrice decimal.Decimal
			if buyOrder.Timestamp < time.Now().Unix() {
				executionPrice = buyOrder.Price
			} else {
				executionPrice = price
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
				symbol, executionAmount, executionPrice, time.Now().Unix(),
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
			// Update remaining amount in database
			err := database.UpdateOrderStatus(e.db, orderID, "open", remainingAmount, 0)
			if err != nil {
				e.logger.Printf("Error updating order status: %v", err)
			}

			// Add to sellers heap
			orderIDInt, _ := strconv.Atoi(orderID)
			sellerOrder := pool.NewOrder(
				uint(orderIDInt),
				uint(remainingAmount.IntPart()),
				price,
				time.Now(),
			)
			stockNode.GetValue().GetSellers().SafePush(sellerOrder)
		}
	}
}

// executeMatch executes a trade between a buy order and a sell order
func (e *Exchange) executeMatch(buyOrderID, buyerAccountID, sellOrderID, sellerAccountID,
	symbol string, amount, price decimal.Decimal, timestamp int64) error {
	// Execute the match within a transaction to ensure atomicity
	return database.ExecuteWithTransaction(e.db, func(tx *sql.Tx) error {
		txFuncs := &database.CommonTxFunctions{Tx: tx}
		// TODO check the false handle
		// 1. Record execution for buy order
		err := txFuncs.RecordExecution(buyOrderID, amount, price, timestamp)
		if err != nil {
			return err
		}

		// 2. Record execution for sell order
		err = txFuncs.RecordExecution(sellOrderID, amount, price, timestamp)
		if err != nil {
			return err
		}

		// 3. Get current buy order status
		buyOrder, err := database.GetOrder(e.db, buyOrderID)
		if err != nil {
			return err
		}

		// 4. Get current sell order status
		sellOrder, err := database.GetOrder(e.db, sellOrderID)
		if err != nil {
			return err
		}

		// 5. Update buy order remaining amount
		newBuyRemaining := buyOrder.Remaining.Sub(amount)
		buyStatus := "open"
		if newBuyRemaining.IsZero() {
			buyStatus = "executed"
		}
		err = txFuncs.UpdateOrderStatus(buyOrderID, buyStatus, newBuyRemaining, 0)
		if err != nil {
			return err
		}

		// 6. Update sell order remaining amount
		newSellRemaining := sellOrder.Remaining.Sub(amount)
		sellStatus := "open"
		if newSellRemaining.IsZero() {
			sellStatus = "executed"
		}
		err = txFuncs.UpdateOrderStatus(sellOrderID, sellStatus, newSellRemaining, 0)
		if err != nil {
			return err
		}

		// 7. Transfer funds to seller
		sellerAccount, err := database.GetAccount(e.db, sellerAccountID)
		if err != nil {
			return err
		}

		// Calculate trade amount
		tradeAmount := amount.Mul(price)
		newSellerBalance := sellerAccount.Balance.Add(tradeAmount)
		err = txFuncs.UpdateAccountBalance(sellerAccountID, newSellerBalance)
		if err != nil {
			return err
		}

		// 8. Update buyer's position
		positions, err := database.GetPositions(e.db, buyerAccountID)
		if err != nil {
			return err
		}

		// Find the position for this symbol
		var currentPosition *database.Position
		for i := range positions {
			if positions[i].Symbol == symbol {
				currentPosition = &positions[i]
				break
			}
		}

		if currentPosition == nil {
			// Create new position
			err = txFuncs.CreateOrUpdatePosition(buyerAccountID, symbol, amount)
		} else {
			// Update existing position
			newAmount := currentPosition.Amount.Add(amount)
			err = txFuncs.CreateOrUpdatePosition(buyerAccountID, symbol, newAmount)
		}

		if err != nil {
			return err
		}

		e.logger.Printf("Executed match: %s bought %s %s from %s at %s",
			buyerAccountID, amount.String(), symbol, sellerAccountID, price.String())

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
	now := time.Now().Unix()

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
