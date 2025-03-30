package main

import (
	_ "github.com/lib/pq"
)

func main() {

	dbm := databaseMaster{
		connStr: "user=postgres password=passw0rd host=localhost port=5432 sslmode=disable",
		dbName:  "stockoverflow",
	}
	dbm.connect()
	dbm.createDB()

}
