package exchange_test

import (
	"StockOverflow/internal/exchange"
	"StockOverflow/internal/pool"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// setupMockDB sets up a mock database for testing
func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	return db, mock
}

// TestPlaceOrder tests the PlaceOrder function
func TestPlaceOrder(t *testing.T) {
	// Setup
	db, mock := setupMockDB(t)
	defer db.Close()

	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)
	stockPool := pool.NewPool(100)
	exchange := exchange.NewExchange(db, stockPool, logger)

	// Test data
	orderID := "12345"
	accountID := "acc123"
	symbol := "AAPL"
	amount := decimal.NewFromInt(10)
	price := decimal.NewFromFloat(150.50)
	now := time.Now().UnixNano()
	fmt.Println("Current time:", now)

	// Mock ExpectExec for order creation
	mock.ExpectExec("INSERT INTO orders").
		WithArgs(orderID, accountID, symbol, amount, price, "open", amount, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock for order match check
	mock.ExpectQuery("SELECT (.+) FROM orders").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "account_id", "symbol", "amount", "price", "status", "remaining", "timestamp", "canceled_time"}))

	// Call function under test
	fmt.Println("Before PlaceOrder call")
	err := exchange.PlaceOrder(orderID, accountID, symbol, amount, price)
	fmt.Println("After PlaceOrder call")

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
