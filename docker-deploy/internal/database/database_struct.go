package database

import (
	"github.com/shopspring/decimal"
)

// Account represents an account in the database
type Account struct {
	ID      string          // account ID
	Balance decimal.Decimal // account balance in USD
}

// Position represents a position (holding) in the database
type Position struct {
	AccountID string          // account ID
	Symbol    string          // symbol of the stock/commodity
	Amount    decimal.Decimal // amount of shares/units held
}

// Order represents an order in the database
type Order struct {
	ID           string          // order ID
	AccountID    string          // account ID that placed the order
	Symbol       string          // symbol being traded
	Amount       decimal.Decimal // amount to trade (negative for sell, positive for buy)
	Price        decimal.Decimal // limit price
	Status       string          // "open", "executed", or "canceled"
	Remaining    decimal.Decimal // remaining amount to be executed
	Timestamp    int64           // timestamp when order was placed
	CanceledTime int64           // timestamp when order was canceled (if applicable)
}

// Execution represents an order execution (trade) in the database
type Execution struct {
	OrderID   string          // order ID that was executed
	Shares    decimal.Decimal // number of shares/units executed
	Price     decimal.Decimal // execution price
	Timestamp int64           // timestamp when execution occurred
}

// Symbol represents a symbol in the database
type Symbol struct {
	Symbol string // symbol name
}
