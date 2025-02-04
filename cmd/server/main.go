package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	// db connection bool
	Sqlitedb *sql.DB
	Port string
	Address string
}

func NewServer() (*Server , error){
	dbconn , err := sql.Open("sqlite3",os.Getenv("DBdir"))
	if err != nil {
		return nil , err
	}
	
	server := &Server{
		Sqlitedb: dbconn,
		Port: os.Getenv("ServerPort"),
	}
	server.Address = os.Getenv("ServerAddr") + ":" + server.Port
	return server , nil
}

func init() {
	err := godotenv.Load() // load environement variables
	if err != nil {
		log.Fatal(err)
	}
}

func HandleConnection(conn *net.Conn) {
	log.Println("Serving " + (*conn).RemoteAddr().String())

	route := make([]byte, 1024)
	n , err := (*conn).Read(route)
	if err != nil {
		log.Println("ERROR" , err)
		return
	}

	request := strings.Split(string(route[0 : n]) , " ")
	var UserName , password string
	var UID int
	if request[0] == "Login" {
		// Todo: impelement authantication handler
		
	} else {
		log.Println("ERROR: BAD request")
	}

}


func main() {

	server , err := NewServer()
	if err != nil{
		log.Fatal(err)
	}

	listener , err := net.Listen("tcp", server.Address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("listening to " + server.Address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go HandleConnection(&conn)

		defer conn.Close()
	}


}