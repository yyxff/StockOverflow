package xmlparser

import (
	"encoding/xml"

	"github.com/shopspring/decimal"
)

type Create struct {
	XMLName  xml.Name `xml:"create"`
	Children []any    `xml:"any"`
}

// Transaction represents the root element for transactions operations
type Transaction struct {
	XMLName  xml.Name `xml:"transactions"`
	ID       string   `xml:"id,attr"`
	Children []any    `xm':"any"`
}

// Order represents an order request
type Order struct {
	Symbol     string          `xml:"sym,attr"`
	Amount     int             `xml:"amount,attr"`
	LimitPrice decimal.Decimal `xml:"limit,attr"`
}

// Query represents an order query
type Query struct {
	ID string `xml:"id,attr"`
}

// Cancel represents an order cancellation
type Cancel struct {
	ID string `xml:"id,attr"`
}
type Account struct {
	ID      string          `xml:"id,attr"`
	Balance decimal.Decimal `xml:"balance,attr"`
}

type Position struct {
	Symbol string  `xml:"symbol"`
	Amount float64 `xml:"amount"`
}

type Symbol struct {
	Symbol   string            `xml:"sym,attr"`
	Accounts []AccountInSymbol `xml:"account"`
}

type AccountInSymbol struct {
	ID     string          `xml:"id,attr"`
	Amount decimal.Decimal `xml:",chardata"`
}
