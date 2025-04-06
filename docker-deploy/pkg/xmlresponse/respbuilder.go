package xmlresponse

import (
	"encoding/xml"
	"fmt"
)

// MarshalXML
func (r Results) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "results"}
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for _, child := range r.Children {
		switch v := child.(type) {
		case Created:
			createdStart := xml.StartElement{
				Name: xml.Name{Local: "created"},
				Attr: []xml.Attr{},
			}

			if v.ID != "" {
				createdStart.Attr = append(createdStart.Attr, xml.Attr{Name: xml.Name{Local: "id"}, Value: v.ID})
			}
			if v.Symbol != "" {
				createdStart.Attr = append(createdStart.Attr, xml.Attr{Name: xml.Name{Local: "sym"}, Value: v.Symbol})
			}

			if err := e.EncodeToken(createdStart); err != nil {
				return err
			}
			if err := e.EncodeToken(xml.EndElement{Name: createdStart.Name}); err != nil {
				return err
			}

		case Error:
			errorStart := xml.StartElement{
				Name: xml.Name{Local: "error"},
				Attr: []xml.Attr{},
			}

			if v.ID != "" {
				errorStart.Attr = append(errorStart.Attr, xml.Attr{Name: xml.Name{Local: "id"}, Value: v.ID})
			}
			if v.Symbol != "" {
				errorStart.Attr = append(errorStart.Attr, xml.Attr{Name: xml.Name{Local: "sym"}, Value: v.Symbol})
			}
			if v.Amount != 0 {
				errorStart.Attr = append(errorStart.Attr, xml.Attr{Name: xml.Name{Local: "amount"}, Value: fmt.Sprintf("%g", v.Amount)})
			}
			if v.Limit != 0 {
				errorStart.Attr = append(errorStart.Attr, xml.Attr{Name: xml.Name{Local: "limit"}, Value: fmt.Sprintf("%g", v.Limit)})
			}

			if err := e.EncodeToken(errorStart); err != nil {
				return err
			}
			if err := e.EncodeToken(xml.CharData([]byte(v.Message))); err != nil {
				return err
			}
			if err := e.EncodeToken(xml.EndElement{Name: errorStart.Name}); err != nil {
				return err
			}

		case Opened:
			if err := e.EncodeElement(v, xml.StartElement{Name: xml.Name{Local: "opened"}}); err != nil {
				return err
			}

		case Status:
			if err := e.EncodeElement(v, xml.StartElement{Name: xml.Name{Local: "status"}}); err != nil {
				return err
			}

		case Canceled:
			if err := e.EncodeElement(v, xml.StartElement{Name: xml.Name{Local: "canceled"}}); err != nil {
				return err
			}
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}
