package test

import (
	"StockOverflow/internal/server"
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func setupMockDB(t *testing.T) (*sql.DB, *sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	return db, &mock

}

// get response
func captureTCPResponse(xmlData []byte) (string, error) {
	conn, err := net.Dial("tcp", "localhost:12345")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// read the datasize
	dataSize := len(xmlData)
	num := strconv.Itoa(dataSize) + "\n"
	xmlData = append([]byte(num), xmlData...)

	// send data
	_, err = conn.Write(xmlData)
	if err != nil {
		return "", fmt.Errorf("failed to send XML data: %v", err)
	}

	reader := bufio.NewReader(conn)
	// 1. 读取长度（直到 \n）
	lenLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("failed to read length:", err)
		return "", err
	}

	// 去掉 '\n' 并转为整数
	lengthStr := lenLine[:len(lenLine)-1]
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		fmt.Println("length format problem:", err)
		return "", err
	}

	// 2. 读取指定字节数的 XML
	xmlData = make([]byte, length)
	_, err = io.ReadFull(reader, xmlData)
	if err != nil {
		fmt.Println("failed to read xml:", err)
		return "", err
	}

	return string(xmlData), nil
}
func TestMainCreate(t *testing.T) {

	db, mock := setupMockDB(t)
	defer db.Close()

	(*mock).ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM accounts WHERE id = \\$1\\)").
		WithArgs("123456").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	(*mock).ExpectExec("INSERT INTO accounts \\(id, balance\\) VALUES \\(\\$1, \\$2\\)").
		WithArgs("123456", "1000").
		WillReturnResult(sqlmock.NewResult(1, 1))

	(*mock).ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM positions WHERE account_id = \\$1 AND symbol = \\$2\\)").
		WithArgs("123456", "SPY").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	(*mock).ExpectExec("UPDATE positions SET amount = \\$1 WHERE account_id = \\$2 AND symbol = \\$3").
		WithArgs("100000", "123456", "SPY").
		WillReturnResult(sqlmock.NewResult(1, 1))
	serverEntry := server.ServerEntry{}
	go serverEntry.Enter(db)
	// go serverEntry.Enter(nil)

	// wait it to start
	time.Sleep(1 * time.Second)

	// prepare xml
	item := `<?xml version="1.0" encoding="UTF-8"?>
<create>
<account id="123456" balance="1000"/>
<symbol sym="SPY">
<account id="123456">100000</account>
</symbol>
</create>`

	xmlData := []byte(item)

	fmt.Println("ready to send xml")
	// send it by tcp
	response, err := captureTCPResponse(xmlData)
	if err != nil {
		t.Fatalf("Failed to get response: %v", err)
	}
	fmt.Println("get response")

	expectedResponse := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<results>\n  <created id=\"123456\"></created>\n  <created id=\"123456\" sym=\"SPY\"></created>\n</results>"
	if response != expectedResponse {
		t.Errorf("Expected response %q, but got %q", expectedResponse, response)
	}

}
