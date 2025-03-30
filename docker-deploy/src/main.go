package main

import (
	_ "github.com/lib/pq"
)

func main() {

	// new a dbmaster
	dbm := databaseMaster{
		connStr: "user=postgres password=passw0rd host=localhost port=5432 sslmode=disable",
		dbName:  "stockoverflow",
	}

	// connec to db
	dbm.connect()

	// create db if doesn't exist
	dbm.createDB()

}
