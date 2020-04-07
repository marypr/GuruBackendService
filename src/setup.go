package src

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

//Connection is a global postgres connection variable
var Connection *sql.DB

//Connect is a function that is used to open Connection with a dataBase.
func Connect() () {
	var err error
	//For localhosts setup "postgres://postgres:root@localhost:5432/postgres?sslmode=disable"
	//where ://username:password@host:port/dbname
	dbURL := "postgres://postgres:1234@localhost:5432/postgres?sslmode=disable"
	Connection, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	err = Connection.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return
}
