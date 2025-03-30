package xmlparser

import "encoding/xml"

type Create struct {
	XMLName  xml.Name  `xml:"create"`
	Accounts []Account `xml:"account"`
	Symbols  []Symbol  `xml:"symbol"`
}

// Transaction represents the root element for transactions operations
type Transaction struct {
	XMLName xml.Name `xml:"transactions"`
	ID      string   `xml:"id,attr"`
	Orders  []Order  `xml:"order"`
	Queries []Query  `xml:"query"`
	Cancels []Cancel `xml:"cancel"`
}

// Order represents an order request
type Order struct {
	Symbol     string  `xml:"sym,attr"`
	Amount     int     `xml:"amount,attr"`
	LimitPrice float64 `xml:"limit,attr"`
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
	ID      string `xml:"id,attr"`
	Balance int    `xml:"balance,attr"`
}

type Position struct {
	Symbol string  `xml:"symbol"`
	Amount float64 `xml:"amount"`
}

type Symbol struct {
	Symbol   string `xml:"sym,attr"`
	Accounts []struct {
		ID      string `xml:"id,attr"`
		Balance string `xml:",chardata"`
	} `xml:"account"`
}
