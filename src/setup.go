package src

import (
	"database/sql"
	"log"

	//postgresDB driver
	_ "github.com/lib/pq"
)

//connection is a global postgres connection variable
var connection *sql.DB

//Connect is a function that is used to open connection with a dataBase.
func Connect() {
	var err error
	//For localhosts setup "postgres://postgres:root@localhost:5432/postgres?sslmode=disable"
	//where ://username:password@host:port/dbname
	dbURL := "postgres://postgres:1234@localhost:5432/postgres?sslmode=disable"
	connection, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	err = connection.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return
}
