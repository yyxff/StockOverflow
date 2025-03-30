package test

import (
	. "StockOverflow/pkg/xmlparser"
	"fmt"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {

	str := "<account id=\"123456\" balance=\"1000\"/>"

	byteArray := []byte(str)

	parser := Xmlparser{}
	xmlData, dataType, err := parser.Parse(byteArray)

	if err != nil {
		t.Errorf("failed")
	}

	fmt.Println(xmlData)
	fmt.Println(dataType)
	accountData, ok := xmlData.(Account)
	if !ok {
		t.Errorf("failed to cast to accout")
	}
	if accountData.ID != "123456" {
		t.Errorf("failed to get id , should be <123456>")
	}
	fmt.Println("id:", accountData.ID)
	fmt.Println("balance:", reflect.TypeOf(accountData.Balance))
}
