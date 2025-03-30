package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type databaseMaster struct {
	connStr string
	dbName  string
	db      *sql.DB
}

// connect to db
func (dbm *databaseMaster) connect() {
	db, err := sql.Open("postgres", dbm.connStr)
	dbm.db = db
	if err != nil {
		log.Fatal(err)
	}
}

// check if db exist
func (dbm *databaseMaster) checkIfExist() bool {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname='%s')", dbm.dbName)
	err := dbm.db.QueryRow(query).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	return exists
}

// if no db, create it
func (dbm *databaseMaster) createDB() {

	exists := dbm.checkIfExist()
	if !exists {
		fmt.Println("db doesn't exist, creating...")
		_, err := dbm.db.Exec("CREATE DATABASE " + dbm.dbName)
		if err != nil {
			log.Fatal("fail to create db:", err)
		}
		fmt.Println("create db successfully!")
	} else {
		fmt.Println("db exists")
	}
}
