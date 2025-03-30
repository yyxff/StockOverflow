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
	}
}
