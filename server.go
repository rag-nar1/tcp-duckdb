package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/marcboeker/go-duckdb"
)


const (
	dbscheme string = `
		CREATE TABLE USER (
			id INT PRIMARY KEY,
			username VARCHAR UNIQUE
		);
		CREATE TABLE DB (
			id INT PRIMARY KEY,
			dbname VARCHAR,
			userid INT REFERENCES USER(id)
		);
		CREATE UNIQUE INDEX USER_INDEX ON USER (username);
		CREATE UNIQUE INDEX DB_INDEX ON DB (userid , dbname);
	`
	dbname string = "db.duckdb"
)

var (
	currentdir string
	dbpath string
	db *sql.DB
)




/**
 * Function to create a tcp listener
 * @pram -> address string : address to listen to 
 * @return -> net.Listener , error
*/
func CreateListener(address string) (net.Listener , error) { 
	listner , err := net.Listen("tcp", address)
	if err != nil {
		return nil , err
	}
	return listner , nil
}


func HandleConnection(connection net.Conn) error {
	defer connection.Close()

	return nil
}

/**
 * Function to start Accepting connections 
 * @pram -> listner net.Listener
 * @return -> error
*/

func Start(listner net.Listener) error {
	defer listner.Close()
	for {
		connection , err := listner.Accept()
		if err != nil {
			return err
		}
		go HandleConnection(connection)
	}
}


func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}


func init() {
	currentdir , _ = os.Getwd()
	dbpath = currentdir + "/DB/" + dbname
	var err error
	db , err = sql.Open("duckdb" , dbpath)
	check(err)
}

func main() {
	defer db.Close()
	_ , err := db.Exec(dbscheme)
	check(err)
	fmt.Println("created db")

}