package xmlresponse

import "encoding/xml"

// Results is the root element for responses
type Results struct {
	XMLName  xml.Name   `xml:"results"`
	Created  []Created  `xml:"created,omitempty"`
	Errors   []Error    `xml:"error,omitempty"`
	Opened   []Opened   `xml:"opened,omitempty"`
	Statuses []Status   `xml:"status,omitempty"`
	Canceled []Canceled `xml:"canceled,omitempty"`
}

// Created represents a successful creation response
type Created struct {
	ID     string `xml:"id,attr,omitempty"`
	Symbol string `xml:"sym,attr,omitempty"`
}

// Error represents an error response
type Error struct {
	ID      string  `xml:"id,attr,omitempty"`
	Symbol  string  `xml:"sym,attr,omitempty"`
	Amount  float64 `xml:"amount,attr,omitempty"`
	Limit   float64 `xml:"limit,attr,omitempty"`
	Message string  `xml:",chardata"`
}

// Opened represents a successfully opened order
type Opened struct {
	Symbol string  `xml:"sym,attr"`
	Amount float64 `xml:"amount,attr"`
	Limit  float64 `xml:"limit,attr"`
	ID     string  `xml:"id,attr"`
}

// Status represents an order status response
type Status struct {
	ID       string     `xml:"id,attr"`
	Open     []Open     `xml:"open,omitempty"`
	Canceled []Canceled `xml:"canceled,omitempty"`
	Executed []Executed `xml:"executed,omitempty"`
}

// Open represents an open portion of an order
type Open struct {
	Shares float64 `xml:"shares,attr"`
}

// Canceled represents a canceled order or portion
type Canceled struct {
	ID       string     `xml:"id,attr,omitempty"` // Only used at top level
	Shares   float64    `xml:"shares,attr,omitempty"`
	Time     int64      `xml:"time,attr,omitempty"`
	Executed []Executed `xml:"executed,omitempty"` // Only used at top level
}

// Executed represents an executed portion of an order
type Executed struct {
	Shares float64 `xml:"shares,attr"`
	Price  float64 `xml:"price,attr"`
	Time   int64   `xml:"time,attr"`
}

// Position represents a holding of a symbol in an account
type Position struct {
	Symbol string  `xml:"symbol"`
	Amount float64 `xml:"amount"`
}
