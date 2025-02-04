package main

import (
	"log"
	"net"
	"os"
	"database/sql"

	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/joho/godotenv"
)

type Server struct {
	// db connection bools
	Duckdb *sql.DB  
	Sqlitedb *sql.DB
}


func init() {
	err := godotenv.Load() // load environement variables
	if err != nil {
		log.Fatal(err)
	}
	
}


func main() {





}