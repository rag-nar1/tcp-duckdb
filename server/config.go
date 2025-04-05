package server

import (
	"database/sql"
	"log"
	"os"

	"github.com/rag-nar1/tcp-duckdb/globals"
	"github.com/rag-nar1/tcp-duckdb/request_handler"
	"github.com/rag-nar1/tcp-duckdb/utils"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	// db connection bool
	Sqlitedb *sql.DB
	Dbstmt   map[string]*sql.Stmt
	Pool     *request_handler.RequestHandler
	Port     string
	Address  string
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

var Serv *Server

// cread prepared statments to use in executing queries
func (s *Server) PrepareStmt() {
	var tmpStmt *sql.Stmt
	var err error
	for _, stmt := range globals.PreparedStmtStrings {
		tmpStmt, err = s.Sqlitedb.Prepare(stmt[1])
		if err != nil {
			s.ErrorLog.Fatal(err)
		}

		s.Dbstmt[stmt[0]] = tmpStmt
	}
}

// create the only superuser user if not already created "duck" with an initial password "duck"
func (s *Server) CreateSuper() {
	hashedPassword := utils.Hash("duck")
	res, err := s.Sqlitedb.Exec("INSERT OR IGNORE INTO user(username, password, usertype) values('duck', ?, 'super')", hashedPassword)
	if err != nil {
		s.ErrorLog.Fatal(err)
	}
	affected, _ := res.RowsAffected()
	if affected == 1 {
		s.InfoLog.Println("Super user created")
	}
}

func ExecuteScheme(db *sql.DB) error {
	scheme, err := os.ReadFile(os.Getenv("DBdir") + "server/scheme.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(scheme))
	if err != nil {
		return err
	}
	return nil
}

func NewServer() error {
	// Check if database file exists
	log.Println(os.Getenv("DBdir") + "server/" + os.Getenv("ServerDbFile"))	
	dbPath := os.Getenv("DBdir") + "server/" + os.Getenv("ServerDbFile")
	executeScheme := false
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		executeScheme = true
	}

	// Open database connection
	dbconn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// Test connection
	err = dbconn.Ping()
	if err != nil {
		return err
	}

	if executeScheme {
		err = ExecuteScheme(dbconn)
		if err != nil {
			return err
		}
	}

	Serv = &Server{
		Sqlitedb: dbconn,
		Port:     os.Getenv("ServerPort"),
		Dbstmt:   make(map[string]*sql.Stmt),
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		Pool:     request_handler.NewRequestHandler(),
	}
	Serv.Address = os.Getenv("ServerAddr") + ":" + Serv.Port
	return nil
}

func Init() {
	// Create default loggers for initialization errors before Serv is ready
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// First try to use existing environment variables

	// Only try to load .env file if ServerAddr is not set
	if os.Getenv("ServerAddr") == "" {
		// Try to load .env file from different locations
		err1 := godotenv.Load(".env")
		err2 := godotenv.Load("../.env")
		if err1 != nil && err2 != nil {
			errorLog.Fatal("Failed to load .env file:", err1, err2)
		}
	}

	if err := NewServer(); err != nil {
		panic(err)
	}

	Serv.CreateSuper()
	Serv.PrepareStmt()
	go Serv.Pool.Spin()
}
