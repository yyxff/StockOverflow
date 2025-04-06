package database

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
)

// ===================== Account Operations =====================

// CreateAccount creates a new account in the database
func CreateAccount(db *sql.DB, id string, balance decimal.Decimal) error {
	// Check if account already exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if account exists: %v", err)
	}

	if exists {
		return fmt.Errorf("account already exists: %s", id)
	}

	// Create the account
	_, err = db.Exec("INSERT INTO accounts (id, balance) VALUES ($1, $2)", id, balance)
	if err != nil {
		return fmt.Errorf("error creating account: %v", err)
	}

	return nil
}

// GetAccount retrieves an account from the database
func GetAccount(db *sql.DB, id string) (*Account, error) {
	var account Account
	err := db.QueryRow("SELECT id, balance FROM accounts WHERE id = $1", id).Scan(&account.ID, &account.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found: %s", id)
		}
		return nil, fmt.Errorf("error retrieving account: %v", err)
	}

	return &account, nil
}

// UpdateAccountBalance updates an account's balance
func UpdateAccountBalance(db *sql.DB, id string, balance decimal.Decimal) error {
	_, err := db.Exec("UPDATE accounts SET balance = $1 WHERE id = $2", balance, id)
	if err != nil {
		return fmt.Errorf("error updating account balance: %v", err)
	}

	return nil
}

// ===================== Position Operations =====================

// GetPositions retrieves all positions for an account
func GetPositions(db *sql.DB, accountID string) ([]Position, error) {
	rows, err := db.Query("SELECT account_id, symbol, amount FROM positions WHERE account_id = $1", accountID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving positions: %v", err)
	}
	defer rows.Close()

	var positions []Position
	for rows.Next() {
		var pos Position
		if err := rows.Scan(&pos.AccountID, &pos.Symbol, &pos.Amount); err != nil {
			return nil, fmt.Errorf("error scanning position: %v", err)
		}
		positions = append(positions, pos)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating positions: %v", err)
	}

	return positions, nil
}

// GetPosition retrieves a specific position for an account and symbol
func GetPosition(db *sql.DB, accountID string, symbol string) (*Position, error) {
	var position Position
	err := db.QueryRow(
		"SELECT account_id, symbol, amount FROM positions WHERE account_id = $1 AND symbol = $2",
		accountID, symbol).Scan(&position.AccountID, &position.Symbol, &position.Amount)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("position not found for account %s and symbol %s", accountID, symbol)
		}
		return nil, fmt.Errorf("error retrieving position: %v", err)
	}

	return &position, nil
}

// CreateOrUpdatePosition creates or updates a position
func CreateOrUpdatePosition(db *sql.DB, accountID string, symbol string, amount decimal.Decimal) error {
	// Check if position exists
	var exists bool
	err := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM positions WHERE account_id = $1 AND symbol = $2)",
		accountID, symbol).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if position exists: %v", err)
	}

	if exists {
		// Update existing position
		_, err = db.Exec("UPDATE positions SET amount = $1 WHERE account_id = $2 AND symbol = $3",
			amount, accountID, symbol)
		if err != nil {
			return fmt.Errorf("error updating position: %v", err)
		}
	} else {
		// Create new position
		_, err = db.Exec("INSERT INTO positions (account_id, symbol, amount) VALUES ($1, $2, $3)",
			accountID, symbol, amount)
		if err != nil {
			return fmt.Errorf("error creating position: %v", err)
		}
	}

	return nil
}

// ===================== Symbol Operations =====================

// CreateSymbol creates a new symbol if it doesn't exist
func CreateSymbol(db *sql.DB, symbol string) error {
	// Check if symbols table exists in the schema
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.tables 
			WHERE table_name = 'symbols'
		)
	`).Scan(&exists)

	if err != nil {
		return fmt.Errorf("error checking if symbols table exists: %v", err)
	}

	// If symbols table doesn't exist, nothing to do
	if !exists {
		return nil
	}

	// Check if symbol exists
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM symbols WHERE symbol = $1)", symbol).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if symbol exists: %v", err)
	}

	if exists {
		// Symbol already exists, which is fine
		return nil
	}

	// Create the symbol
	_, err = db.Exec("INSERT INTO symbols (symbol) VALUES ($1)", symbol)
	if err != nil {
		return fmt.Errorf("error creating symbol: %v", err)
	}

	return nil
}

// GetAllSymbols retrieves all symbols from the database
func GetAllSymbols(db *sql.DB) ([]Symbol, error) {
	// Check if symbols table exists in the schema
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.tables 
			WHERE table_name = 'symbols'
		)
	`).Scan(&exists)

	if err != nil {
		return nil, fmt.Errorf("error checking if symbols table exists: %v", err)
	}

	// If symbols table doesn't exist, return empty slice
	if !exists {
		return []Symbol{}, nil
	}

	rows, err := db.Query("SELECT symbol FROM symbols")
	if err != nil {
		return nil, fmt.Errorf("error retrieving symbols: %v", err)
	}
	defer rows.Close()

	var symbols []Symbol
	for rows.Next() {
		var sym Symbol
		if err := rows.Scan(&sym.Symbol); err != nil {
			return nil, fmt.Errorf("error scanning symbol: %v", err)
		}
		symbols = append(symbols, sym)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating symbols: %v", err)
	}

	return symbols, nil
}

// ===================== Order Operations =====================

// CreateOrder creates a new order in the database
func CreateOrder(db *sql.DB, order *Order) error {
	_, err := db.Exec(
		"INSERT INTO orders (id, account_id, symbol, amount, price, status, remaining, timestamp) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		order.ID, order.AccountID, order.Symbol, order.Amount, order.Price,
		order.Status, order.Remaining, order.Timestamp)
	if err != nil {
		return fmt.Errorf("error creating order: %v", err)
	}

	return nil
}

// GetOrder retrieves an order from the database
func GetOrder(db *sql.DB, orderID string) (*Order, error) {
	var order Order
	var canceledTime sql.NullInt64 // Use sql.NullInt64 struct to handle NULL values

	err := db.QueryRow(
		"SELECT id, account_id, symbol, amount, price, status, remaining, timestamp, canceled_time "+
			"FROM orders WHERE id = $1", orderID).Scan(
		&order.ID, &order.AccountID, &order.Symbol, &order.Amount, &order.Price,
		&order.Status, &order.Remaining, &order.Timestamp, &canceledTime)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found: %s", orderID)
		}
		return nil, fmt.Errorf("error retrieving order: %v", err)
	}

	// Convert NullInt64 to int64 (0 if NULL)
	if canceledTime.Valid {
		order.CanceledTime = canceledTime.Int64
	} else {
		order.CanceledTime = 0 // Use 0 or another default value for NULL
	}

	return &order, nil
}

// UpdateOrderStatus updates an order's status and remaining amount
func UpdateOrderStatus(db *sql.DB, orderID string, status string, remaining decimal.Decimal, canceledTime int64) error {
	var query string
	var args []interface{}

	if canceledTime > 0 {
		query = "UPDATE orders SET status = $1, remaining = $2, canceled_time = $3 WHERE id = $4"
		args = []interface{}{status, remaining, canceledTime, orderID}
	} else {
		query = "UPDATE orders SET status = $1, remaining = $2 WHERE id = $3"
		args = []interface{}{status, remaining, orderID}
	}

	_, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error updating order status: %v", err)
	}

	return nil
}

// GetOpenOrdersBySymbol retrieves all open orders for a symbol
func GetOpenOrdersBySymbol(db *sql.DB, symbol string) ([]Order, error) {
	rows, err := db.Query(
		"SELECT id, account_id, symbol, amount, price, status, remaining, timestamp, canceled_time "+
			"FROM orders WHERE symbol = $1 AND status = 'open'", symbol)
	if err != nil {
		return nil, fmt.Errorf("error retrieving open orders: %v", err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.AccountID, &order.Symbol, &order.Amount,
			&order.Price, &order.Status, &order.Remaining, &order.Timestamp, &order.CanceledTime); err != nil {
			return nil, fmt.Errorf("error scanning order: %v", err)
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %v", err)
	}

	return orders, nil
}

func GetOpenOrdersBySymbolForHeap(db *sql.DB, symbol string, target string, limit int) ([]Order, error) {

	condition := ""
	orderStr := ""
	if target == "buyer" {
		condition = " AND amount > 0"
		orderStr = " ORDER BY price DESC, timestamp ASC"
	} else if target == "seller" {
		condition = " AND amount < 0"
		orderStr = " ORDER BY price ASC, timestamp ASC"
	} else {
		return nil, fmt.Errorf("target should be buyer or seller, but get%s", target)
	}

	limitStr := strconv.Itoa(limit)
	sqlStr := "SELECT id, account_id, symbol, amount, price, status, remaining, timestamp, canceled_time " +
		"FROM orders WHERE symbol = $1 AND status = 'open'" + condition + orderStr + " LIMIT " + limitStr
	rows, err := db.Query(sqlStr, symbol)
	if err != nil {
		return nil, fmt.Errorf("error retrieving open orders: %v", err)
	}
	defer rows.Close()

	var orders []Order
	var canceledTime sql.NullInt64 // Use sql.NullInt64 to handle NULL values

	for rows.Next() {
		var order Order

		if err := rows.Scan(&order.ID, &order.AccountID, &order.Symbol, &order.Amount,
			&order.Price, &order.Status, &order.Remaining, &order.Timestamp, &canceledTime); err != nil {
			return nil, fmt.Errorf("error scanning order: %v", err)
		}

		// Convert NullInt64 to int64 (0 if NULL)
		if canceledTime.Valid {
			order.CanceledTime = canceledTime.Int64
		} else {
			order.CanceledTime = 0 // Use 0 or another default value for NULL
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %v", err)
	}

	return orders, nil
}

// ===================== Execution Operations =====================

// RecordExecution creates a new execution record in the database
func RecordExecution(db *sql.DB, execution *Execution) error {
	_, err := db.Exec(
		"INSERT INTO executions (order_id, shares, price, timestamp) VALUES ($1, $2, $3, $4)",
		execution.OrderID, execution.Shares, execution.Price, execution.Timestamp)
	if err != nil {
		return fmt.Errorf("error creating execution: %v", err)
	}

	return nil
}

// GetOrderExecutions retrieves all executions for an order
func GetOrderExecutions(db *sql.DB, orderID string) ([]Execution, error) {
	rows, err := db.Query(
		"SELECT order_id, shares, price, timestamp FROM executions WHERE order_id = $1 ORDER BY timestamp",
		orderID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving executions: %v", err)
	}
	defer rows.Close()

	var executions []Execution
	for rows.Next() {
		var exec Execution
		if err := rows.Scan(&exec.OrderID, &exec.Shares, &exec.Price, &exec.Timestamp); err != nil {
			return nil, fmt.Errorf("error scanning execution: %v", err)
		}
		executions = append(executions, exec)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating executions: %v", err)
	}

	return executions, nil
}

// ===================== Transaction Helpers =====================

// ExecuteWithTransaction executes a function within a database transaction
func ExecuteWithTransaction(db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; rollback
		} else {
			err = tx.Commit() // err is nil; commit
			if err != nil {
				err = fmt.Errorf("failed to commit transaction: %v", err)
			}
		}
	}()

	err = fn(tx)
	return err
}

// CommonTxFunctions contains common database operations that can be used within a transaction
type CommonTxFunctions struct {
	Tx *sql.Tx
}

// CreateAccount creates a new account within a transaction
func (f *CommonTxFunctions) CreateAccount(id string, balance decimal.Decimal) error {
	_, err := f.Tx.Exec("INSERT INTO accounts (id, balance) VALUES ($1, $2)", id, balance)
	if err != nil {
		return fmt.Errorf("error creating account in transaction: %v", err)
	}
	return nil
}

// UpdateAccountBalance updates an account's balance within a transaction
func (f *CommonTxFunctions) UpdateAccountBalance(id string, balance decimal.Decimal) error {
	_, err := f.Tx.Exec("UPDATE accounts SET balance = $1 WHERE id = $2", balance, id)
	if err != nil {
		return fmt.Errorf("error updating account balance in transaction: %v", err)
	}
	return nil
}

// CreateOrUpdatePosition creates or updates a position within a transaction
func (f *CommonTxFunctions) CreateOrUpdatePosition(accountID string, symbol string, amount decimal.Decimal) error {
	// Check if position exists
	var exists bool
	err := f.Tx.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM positions WHERE account_id = $1 AND symbol = $2)",
		accountID, symbol).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if position exists in transaction: %v", err)
	}

	if exists {
		// Update existing position
		_, err = f.Tx.Exec("UPDATE positions SET amount = $1 WHERE account_id = $2 AND symbol = $3",
			amount, accountID, symbol)
		if err != nil {
			return fmt.Errorf("error updating position in transaction: %v", err)
		}
	} else {
		// Create new position
		_, err = f.Tx.Exec("INSERT INTO positions (account_id, symbol, amount) VALUES ($1, $2, $3)",
			accountID, symbol, amount)
		if err != nil {
			return fmt.Errorf("error creating position in transaction: %v", err)
		}
	}

	return nil
}

// RecordExecution creates a new execution record within a transaction
func (f *CommonTxFunctions) RecordExecution(orderID string, shares decimal.Decimal, price decimal.Decimal, timestamp int64) error {
	_, err := f.Tx.Exec(
		"INSERT INTO executions (order_id, shares, price, timestamp) VALUES ($1, $2, $3, $4)",
		orderID, shares, price, timestamp)
	if err != nil {
		return fmt.Errorf("error creating execution in transaction: %v", err)
	}
	return nil
}

// UpdateOrderStatus updates an order's status within a transaction
func (f *CommonTxFunctions) UpdateOrderStatus(orderID string, status string, remaining decimal.Decimal, canceledTime int64) error {
	var query string
	var args []interface{}

	if canceledTime > 0 {
		query = "UPDATE orders SET status = $1, remaining = $2, canceled_time = $3 WHERE id = $4"
		args = []interface{}{status, remaining, canceledTime, orderID}
	} else {
		query = "UPDATE orders SET status = $1, remaining = $2 WHERE id = $3"
		args = []interface{}{status, remaining, orderID}
	}

	_, err := f.Tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error updating order status in transaction: %v", err)
	}

	return nil
}

// ===================== Server Start Helpers =====================

// GetMaxOrderID retrieves the highest order ID from the database starting from server start
func GetMaxOrderID(db *sql.DB) (int, error) {
	// Check if orders table exists and has records
	var exists bool
	err := db.QueryRow(`
        SELECT EXISTS (
            SELECT 1 
            FROM information_schema.tables 
            WHERE table_name = 'orders'
        )
    `).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("error checking if orders table exists: %v", err)
	}

	// If orders table doesn't exist, return default
	if !exists {
		return 0, nil
	}

	// Get all order IDs from the database
	rows, err := db.Query("SELECT id FROM orders")
	if err != nil {
		return 0, fmt.Errorf("error querying order IDs: %v", err)
	}
	defer rows.Close()

	maxID := 0

	// Iterate through all IDs and find the maximum numerical value
	for rows.Next() {
		var idStr string
		if err := rows.Scan(&idStr); err != nil {
			return 0, fmt.Errorf("error scanning order ID: %v", err)
		}

		// Try to convert string ID to integer
		id, err := strconv.Atoi(idStr)
		if err != nil {
			// Skip non-numeric IDs
			continue
		}

		// Update maxID if this ID is larger
		if id > maxID {
			maxID = id
		}
	}

	if err = rows.Err(); err != nil {
		return 0, fmt.Errorf("error iterating order IDs: %v", err)
	}

	return maxID, nil
}
