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
		case "account":
			{
				var account Account
				err := decoder.DecodeElement(&account, &startElement)
				fmt.Println(account.ID)
				fmt.Println(account.Balance)
				if err != nil {
					fmt.Println("error:", err)
				}
				return account, reflect.TypeOf(account), err
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
