package test

import (
	"StockOverflow/internal/server"
	. "StockOverflow/pkg/xmlparser"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
)

func setupMockDB(t *testing.T) *sql.DB {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()
	return db

}

// get response
func captureTCPResponse(xmlData []byte) (string, error) {
	conn, err := net.Dial("tcp", "localhost:12345")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// read the datasize
	dataSize := int32(len(xmlData))

	num := string(dataSize)

	// 先发送字节数
	_, err = conn.Write([]byte(num))
	if err != nil {
		return "", fmt.Errorf("failed to send XML data: %v", err)
	}

	_, err = conn.Write([]byte{'\n'})
	if err != nil {
		return "", fmt.Errorf("failed to send XML data: %v", err)
	}

	_, err = conn.Write(xmlData)
	if err != nil {
		return "", fmt.Errorf("failed to send XML data: %v", err)
	}

	var response [1024]byte
	n, err := conn.Read(response[:])
	if err != nil {
		return "", err
	}

	log.Printf("Expected data size: %d bytes", dataSize)

	// 根据字节数读取 XML 数据
	data := make([]byte, dataSize)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		log.Printf("Failed to read XML data: %v", err)
		return "", err
	}

	return string(response[:n]), nil
}
func TestMain(t *testing.T) {
	serverEntry := server.ServerEntry{}
	go serverEntry.Enter(setupMockDB(t))

	// wait it to start
	time.Sleep(2 * time.Second)

	// prepare xml
	item := &Create{
		XMLName: xml.Name{Local: "create"},
		Accounts: []Account{
			{
				ID:      "testuserid",
				Balance: decimal.NewFromFloat(100),
			},
		},
		Symbols: []Symbol{
			{
				Symbol: "duke",
				Accounts: []AccountInSymbol{
					{
						ID:      "testuserid",
						Balance: decimal.NewFromFloat(999),
					},
				},
			},
		},
	}
	xmlData, err := xml.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal XML: %v", err)
	}

	fmt.Println("ready to send xml")
	// send it by tcp
	response, err := captureTCPResponse(xmlData)
	if err != nil {
		t.Fatalf("Failed to get response: %v", err)
	}
	fmt.Println("get response")

	expectedResponse := "Received item with value: 100"
	if response != expectedResponse {
		t.Errorf("Expected response %q, but got %q", expectedResponse, response)
	}

}
