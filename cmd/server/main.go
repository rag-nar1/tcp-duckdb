package main

import (
	"TCP-Duckdb/utils"
	"bufio"
	"database/sql"
	"log"
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3"
)
var preparedStmtStrings = [][]string{
	{"login",               "SELECT userid , usertype FROM user WHERE username LIKE ? AND password LIKE ? ;"},
	{"SelectUser",          "SELECT userid FROM user WHERE username LIKE ? ;"},
	{"CreateUser",          "INSERT INTO user(username, password, usertype) VALUES(?, ?, ?);"},
	{"SelectDB",            "SELECT dbid FROM DB WHERE dbname LIKE ? ;"},
	{"CreateDB",            "INSERT INTO DB(dbname) VALUES(?);"},
	{"GrantDB",             "INSERT OR IGNORE INTO dbprivilege(dbid, userid, privilegetype) VALUES(?, ?, ?);"},
	{"CheckDbAccess",       "SELECT COUNT(*) FROM dbprivilege WHERE userid == ? AND dbid == ?"},
	{"SelectTable",         "SELECT tableid FROM tables WHERE tablename LIKE ? AND dbid == ?;"},
	{"CheckTableAccess",    "SELECT COUNT(*) FROM tableprivilege WHERE userid == ? AND tableid == ? AND tableprivilege LIKE ?;"},
	{"GrantTable",          "INSERT OR IGNORE INTO tableprivilege(tableid, userid, tableprivilege) VALUES(?, ?, ?);"},
	{"CreateTable",         "INSERT OR IGNORE INTO tables(tablename, dbid) VALUES(?, ?);"},
	{"CreateLink",          "INSERT OR IGNORE INTO postgres(dbid, connstr) VALUES(?, ?);"},
	{"CreateKey",           "INSERT OR IGNORE INTO keys(dbid, key) VALUES(?, ?);"},
}

var infoLog, errorLog *log.Logger

type PreparedStmts map[string] *sql.Stmt

func Write(writer *bufio.Writer, data []byte) error {
	data = append(data, '\n')
	if _, err := writer.Write(data); err != nil {
        return err
    }
    return writer.Flush()
}

type Server struct {
	// db connection bool
	Sqlitedb 	*sql.DB
	dbstmt 		PreparedStmts
	Port 		string
	Address 	string
}

// cread prepared statments to use in executing queries 
func (s *Server) prepareStmt() {
	var tmpStmt *sql.Stmt 
	var err error
	for _ , stmt := range preparedStmtStrings {
		tmpStmt , err = s.Sqlitedb.Prepare(stmt[1])
		if err != nil {
			errorLog.Fatal(err)
		}

		s.dbstmt[stmt[0]] = tmpStmt
	}
}
var server *Server

// create the only superuser user if not already created "duck" with an initial password "duck"  
func (s *Server) CreateSuper() {
	hashedPassword := utils.Hash("duck")
	res , err := s.Sqlitedb.Exec("INSERT OR IGNORE INTO user(username, password, usertype) values('duck', ?, 'super')",hashedPassword)
	if err != nil {
		errorLog.Fatal(err)
	}
	affected , _ := res.RowsAffected()
	if affected == 1 {
		infoLog.Println("Super user created")
	}
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	infoLog.Println("Serving " + conn.RemoteAddr().String())
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	// read login request
	route := make([]byte, 1024)
	n , err := reader.Read(route)
	if err != nil {
		errorLog.Println(err)
		return
	}
	// check for a valid request
	request := strings.Split(string(route[0 : n]) , " ")
	var UserName , password , privilege string
	var UID int
	if request[0] != "login" || len(request) != 3 {
		Write(writer, []byte("ERROR: BAD request\n"))
		return
	}
	// validate the username and password
	UserName, password = utils.Trim(request[1]), utils.Trim(request[2]) 
	UID, privilege, err = LoginHandler(UserName, password, server.dbstmt["login"])
	if err != nil {
		Write(writer, []byte("Unauthorized\n"))
		errorLog.Println(err)
		return	
	}
	Write(writer, []byte("success\n"))
	DBHandler(UID, UserName, privilege, reader, writer)
}

	
func NewServer() (error){
	dbconn , err := sql.Open("sqlite3",os.Getenv("DBdir") + os.Getenv("ServerDbFile"))
	if err != nil {
		return err
	}
	err = dbconn.Ping()
	if err != nil {
		return err
	}

	server = &Server{
		Sqlitedb: dbconn,
		Port: os.Getenv("ServerPort"),
		dbstmt: make(PreparedStmts),
	}
	server.Address = os.Getenv("ServerAddr") + ":" + server.Port
	return nil
}

func init() {
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate | log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate | log.Ltime | log.Lshortfile) 
	err := godotenv.Load() // load environement variables
	if err != nil {
		errorLog.Fatal(err)
	}

	err = NewServer()
	if err != nil {
		errorLog.Fatal(err)
	}
	server.prepareStmt()
	server.CreateSuper()
}

func main() {
	// start listing to tcp connections
	listener , err := net.Listen("tcp", server.Address)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer listener.Close()

	infoLog.Println("listening to " + server.Address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			errorLog.Fatal(err)
		}
		// starts a go routin to handle every connection
		go HandleConnection(conn)
	}

}