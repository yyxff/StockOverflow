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

	var c Create
	err := xml.Unmarshal([]byte(xmlData), &c)
	if err != nil {
		fmt.Println("XML parse error:", err)
	}
	fmt.Printf("%+v\n", c)

	fmt.Printf("Raw: %x\n", xmlData[:10])

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
		// fmt.Println(startElement.Name.Local)

		// if !ok {
		// 	fmt.Println("error in getting startElement")
		// }

	}
	return nil, nil, nil
}
