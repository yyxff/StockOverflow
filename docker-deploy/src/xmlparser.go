package main

import (
	"encoding/xml"
	"fmt"
)

type Xmlparser struct {
}

// parse xml
// return the xml struct
func (parser *Xmlparser) Parse(xmlData []byte) (interface{}, error) {
	var root struct {
		XMLName  xml.Name  `xml:"account"`
		Account  *Account  `xml:"account"`
		Position *Position `xml:"position"`
		Symbol   *Symbol   `xml:"symbol"`
		Order    *Order    `xml:"order"`
	}
	err := xml.Unmarshal(xmlData, &root)

	if err != nil {
		fmt.Println("fail to parse xml: ", err)
		return nil, err
	}

	// return struct
	switch root.XMLName.Local {
	case "account":
		return root.Account, nil
	case "position":
		return root.Position, nil
	case "symbol":
		return root.Symbol, nil
	case "order":
		return root.Order, nil
	default:
		return nil, fmt.Errorf("unknown root element: %s", root.XMLName.Local)
	}
}
