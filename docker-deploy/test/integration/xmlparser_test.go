package test

import (
	. "StockOverflow/pkg/xmlparser"
	"fmt"
	"testing"
)

// func TestParse(t *testing.T) {

// 	str := "<account id=\"123456\" balance=\"1000\"/>"

// 	byteArray := []byte(str)

// 	parser := Xmlparser{}
// 	xmlData, dataType, err := parser.Parse(byteArray)

// 	if err != nil {
// 		t.Errorf("failed")
// 	}

// 	fmt.Println(xmlData)
// 	fmt.Println(dataType)
// 	accountData, ok := xmlData.(Account)
// 	if !ok {
// 		t.Errorf("failed to cast to accout")
// 	}
// 	if accountData.ID != "123456" {
// 		t.Errorf("failed to get id , should be <123456>")
// 	}
// 	fmt.Println("id:", accountData.ID)
// 	fmt.Println("balance:", reflect.TypeOf(accountData.Balance))
// }

func TestParseCreate(t *testing.T) {

	str :=
		`<?xml version="1.0" encoding="UTF-8"?>
<create>
	<account id="123456" balance="1000"/>
	<symbol sym="SPY">
		<account id="123456">100000</account>
	</symbol>
</create>`

	byteArray := []byte(str)
	fmt.Println(str)

	parser := Xmlparser{}
	xmlData, dataType, err := parser.Parse(byteArray)

	if err != nil {
		t.Errorf("failed")
	}

	fmt.Println(xmlData)
	fmt.Println(dataType)
	create, ok := xmlData.(Create)
	// print(reflect.TypeOf(create.Accounts[0].Balance))
	if !ok {
		t.Errorf("failed to cast to accout")
	}
	if create.XMLName.Local != "create" {
		t.Errorf("failed to get name , should be <create>")
	}
	if len(create.Accounts) != 1 {
		t.Errorf("accounts should be 1, but get %d\n", len(create.Accounts))
	}
	if len(create.Symbols) != 1 {
		t.Errorf("symbols should be 1, but get %d\n", len(create.Symbols))
	}
	if len(create.Symbols[0].Accounts) != 1 {
		t.Errorf("symbols should have 1 account, but get %d\n", len(create.Symbols[0].Accounts))
	}
}

func TestParseTransaction(t *testing.T) {

	str :=
		`<transactions id="ACCOUNT_ID"> #contains 1 or more of the below
children
	<order sym="SYM" amount="1000" limit="100.7"/>
	<query id="TRANS_ID"/>
	<cancel id="TRANS_ID"/>
</transactions>`

	byteArray := []byte(str)
	fmt.Println(str)

	parser := Xmlparser{}
	xmlData, dataType, err := parser.Parse(byteArray)

	if err != nil {
		t.Errorf("failed")
	}

	fmt.Println(xmlData)
	fmt.Println(dataType)
	transaction, ok := xmlData.(Transaction)
	if !ok {
		t.Errorf("failed to cast to accout")
	}
	if transaction.XMLName.Local != "transactions" {
		t.Errorf("failed to get name , should be <transaction>")
	}
	if transaction.ID != "ACCOUNT_ID" {
		t.Errorf("accounts should be ACCOUNT_ID, but get %s\n", transaction.ID)
	}
	if transaction.Orders[0].Symbol != "SYM" {
		t.Errorf("symbols should be SYM, but get %s\n", transaction.Orders[0].Symbol)
	}
}
