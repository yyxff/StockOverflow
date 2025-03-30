package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DatabaseMaster struct {
	ConnStr string
	DbName  string
	Db      *sql.DB
}

// connect to db
func (dbm *DatabaseMaster) Connect() {
	db, err := sql.Open("postgres", dbm.ConnStr)
	dbm.Db = db
	if err != nil {
		log.Fatal(err)
	}
}

// check if db exist
func (dbm *DatabaseMaster) CheckIfExist() bool {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname='%s')", dbm.DbName)
	err := dbm.Db.QueryRow(query).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	return exists
}

// if no db, create it
func (dbm *DatabaseMaster) CreateDB() {

	exists := dbm.CheckIfExist()
	if !exists {
		fmt.Println("db doesn't exist, creating...")
		_, err := dbm.Db.Exec("CREATE DATABASE " + dbm.DbName)
		if err != nil {
			log.Fatal("fail to create db:", err)
		}
		fmt.Println("create db successfully!")
	} else {
		fmt.Println("db exists")
		connStr := dbm.ConnStr + "host=localhost port=5432 user=postgres password=passw0rd dbname=stockoverflow sslmode=disable"
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			fmt.Println("failed to switch to stockoverflow")
		}
		dbm.Db = db
	}
}

// init all tables
func (dbm *DatabaseMaster) Init() {
	dbm.initAccountTable()
	dbm.initPositionTable()
	dbm.initOrderTable()
	dbm.initExecutionTable()
}

// init account table
func (dbm *DatabaseMaster) initAccountTable() {

	createTableSQL := `CREATE TABLE IF NOT EXISTS accounts (
		id VARCHAR(255) PRIMARY KEY,
		balance NUMERIC(20, 2) NOT NULL 
		);`

	_, err := dbm.Db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	} else {
		fmt.Println("Table <Account> checked/created successfully.")
	}
}

// init position table
func (dbm *DatabaseMaster) initPositionTable() {

	createTableSQL := `CREATE TABLE IF NOT EXISTS positions (
    account_id VARCHAR(255) REFERENCES accounts(id),
    symbol VARCHAR(255) NOT NULL,
    amount NUMERIC(20, 6) NOT NULL,
    PRIMARY KEY (account_id, symbol)
);`

	_, err := dbm.Db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	} else {
		fmt.Println("Table <Position> checked/created successfully.")
	}
}

// init orders table
func (dbm *DatabaseMaster) initOrderTable() {

	createTableSQL := `CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(255) PRIMARY KEY,
    account_id VARCHAR(255) REFERENCES accounts(id),
    symbol VARCHAR(255) NOT NULL,
    amount NUMERIC(20, 6) NOT NULL,
	price NUMERIC(20, 6) NOT NULL,
    status VARCHAR(10) NOT NULL,
    remaining NUMERIC(20, 6) NOT NULL,
    timestamp BIGINT NOT NULL,
    canceled_time BIGINT
);`

	_, err := dbm.Db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	} else {
		fmt.Println("Table <Orders> checked/created successfully.")
	}
}

// init execution table
func (dbm *DatabaseMaster) initExecutionTable() {

	createTableSQL := `CREATE TABLE IF NOT EXISTS executions (
    order_id VARCHAR(255) REFERENCES orders(id),
    shares NUMERIC(20, 6) NOT NULL,
    price NUMERIC(20, 6) NOT NULL,
    timestamp BIGINT NOT NULL,
    PRIMARY KEY (order_id, timestamp)
);`

	_, err := dbm.Db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	} else {
		fmt.Println("Table <Exxcution> checked/created successfully.")
	}
}
