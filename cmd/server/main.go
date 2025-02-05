package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"strings"
	"TCP-Duckdb/utils"

	"github.com/joho/godotenv"
	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3"

)
var preparedStmtStrings = [][]string{
	{"login", "SELECT userid FROM user WHERE username LIKE ? AND password LIKE ? ;"},
}
type PreparedStmts map[string] *sql.Stmt

type Server struct {
	// db connection bool
	Sqlitedb *sql.DB
	dbstmt PreparedStmts
	Port string
	Address string
}

// cread prepared statments to use in executing queries 
func (s *Server) prepareStmt() {
	var tmpStmt *sql.Stmt 
	var err error
	for _ , stmt := range preparedStmtStrings {
		tmpStmt , err = s.Sqlitedb.Prepare(stmt[1])
		if err != nil {
			log.Fatal(err)
		}

		s.dbstmt[stmt[0]] = tmpStmt
	}
}

// create the only superuser user if not already created "duck" with an initial password "duck"  
func (s *Server) CreateSuper() {
	res , err := s.Sqlitedb.Exec("INSERT OR IGNORE INTO user(username, password, usertype) values('duck', 'duck', 'super')")
	if err != nil {
		log.Fatal(err)
	}
	affected , _ := res.RowsAffected()
	if affected == 1 {
		log.Print("Super user created")
	}
}

func HandleConnection(conn *net.Conn , dbstmt PreparedStmts) {
	defer (*conn).Close()	
	log.Println("Serving " + (*conn).RemoteAddr().String())
	// read login request
	route := make([]byte, 1024)
	n , err := (*conn).Read(route)
	if err != nil {
		log.Println("ERROR" , err)
		return
	}
	// check for a valid request
	request := strings.Split(string(route[0 : n]) , " ")
	var UserName , password string
	var UID int
	if request[0] != "Login" || len(request) != 3 {
		(*conn).Write([]byte("ERROR: BAD request\n"))
		return
	}
	// validate the username and password
	UserName, password = utils.Trim(request[1]), utils.Trim(request[2]) 
	UID , err = LoginHandler(UserName, password, dbstmt["login"])
	if err != nil {
		(*conn).Write([]byte("Unauthorized\n"))
		log.Print(err)
		return	
	}

	log.Printf("USER SERVED %d" , UID)
}

	
func NewServer() (*Server , error){
	dbconn , err := sql.Open("sqlite3",os.Getenv("DBdir"))
	if err != nil {
		return nil , err
	}
	err = dbconn.Ping()
	if err != nil {
		return nil , err
	}

	server := &Server{
		Sqlitedb: dbconn,
		Port: os.Getenv("ServerPort"),
		dbstmt: make(PreparedStmts),
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

func main() {

	server , err := NewServer()
	if err != nil{
		log.Fatal(err)
	}
	server.prepareStmt()
	server.CreateSuper()
	// start listing to tcp connections
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
		// starts a go routin to handle every connection
		go HandleConnection(&conn, server.dbstmt)

		defer conn.Close()
	}

}