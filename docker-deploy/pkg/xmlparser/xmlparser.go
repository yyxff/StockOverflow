package xmlparser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"reflect"
)

type Xmlparser struct {
}

// parse xml
// return the xml struct
func (parser *Xmlparser) Parse(xmlData []byte) (any, reflect.Type, error) {

	decoder := xml.NewDecoder(bytes.NewReader(xmlData))

	for {
		tok, err := decoder.Token()
		if err != nil {
			break
		}

		// get root element
		if startElement, ok := tok.(xml.StartElement); ok {
			switch startElement.Name.Local {
			case "create":
				{
					var create Create
					// create.XMLName.Local = "create"

					err := create.parse(decoder, startElement)
					if err != nil {
						fmt.Println("error:", err)
					}
					return create, reflect.TypeOf(create), err
				}
			case "transactions":
				{
					var transaction Transaction
					err := transaction.parse(decoder, startElement)
					if err != nil {
						fmt.Println("error:", err)
					}
					return transaction, reflect.TypeOf(transaction), err
				}
			default:
				{
					fmt.Println("default")
					return nil, nil, fmt.Errorf("unknown root element: %s", startElement.Name.Local)
				}
			}
		}
		// fmt.Println(startElement.Name.Local)

		// if !ok {
		// 	fmt.Println("error in getting startElement")
		// }

	}
	return nil, nil, nil
}

// parse create in order
func (create *Create) parse(decoder *xml.Decoder, start xml.StartElement) error {
	// var create Create
	create.XMLName = start.Name

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		// check if parent ele ends
		if end, ok := token.(xml.EndElement); ok && end.Name == start.Name {
			return nil
		}

		// process sub element
		if startElem, ok := token.(xml.StartElement); ok {

			// switch by element type label
			var child any
			switch startElem.Name.Local {
			case "account":
				var account Account
				err := decoder.DecodeElement(&account, &startElem)
				if err != nil {
					return err
				}
				child = account
			case "symbol":
				var symbol Symbol
				err := decoder.DecodeElement(&symbol, &startElem)
				if err != nil {
					return err
				}
				child = symbol
			default:
				if err := decoder.Skip(); err != nil {
					return err
				}
				continue
			}

			// record order
			create.Children = append(create.Children, child)
		}
	}
}

// parse transaction in order
func (transaction *Transaction) parse(decoder *xml.Decoder, start xml.StartElement) error {
	// var create Create

	// get name
	transaction.XMLName = start.Name
	// get id
	for _, attr := range start.Attr {
		if attr.Name.Local == "id" {
			transaction.ID = attr.Value
			break
		}
	}
	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		// check if parent ele ends
		if end, ok := token.(xml.EndElement); ok && end.Name == start.Name {
			return nil
		}

		// process sub element
		if startElem, ok := token.(xml.StartElement); ok {

			// switch by element type label
			var child any
			switch startElem.Name.Local {
			case "order":
				var order Order
				err := decoder.DecodeElement(&order, &startElem)
				if err != nil {
					return err
				}
				child = order
			case "query":
				var query Query
				err := decoder.DecodeElement(&query, &startElem)
				if err != nil {
					return err
				}
				child = query
			case "cancel":
				var cancel Cancel
				err := decoder.DecodeElement(&cancel, &startElem)
				if err != nil {
					return err
				}
				child = cancel
			default:
				if err := decoder.Skip(); err != nil {
					return err
				}
				continue
			}

			// record order
			transaction.Children = append(transaction.Children, child)
		}
	}
}
