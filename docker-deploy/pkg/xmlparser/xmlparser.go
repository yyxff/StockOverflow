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
func (parser *Xmlparser) Parse(xmlData []byte) (interface{}, reflect.Type, error) {

	decoder := xml.NewDecoder(bytes.NewReader(xmlData))

	for {
		tok, err := decoder.Token()
		if err != nil {
			break
		}

		// get root element
		startElement, ok := tok.(xml.StartElement)
		fmt.Println(startElement.Name.Local)

		if !ok {
			fmt.Println("error in getting startElement")
		}

		switch startElement.Name.Local {
		case "create":
			{
				var create Create
				err := decoder.DecodeElement(&create, &startElement)
				if err != nil {
					fmt.Println("error:", err)
				}
				return create, reflect.TypeOf(create), err
			}
		case "transactions":
			{
				var transaction Transaction
				err := decoder.DecodeElement(&transaction, &startElement)
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
	return nil, nil, nil
}
