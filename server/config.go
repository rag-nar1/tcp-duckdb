package server

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rag-nar1/TCP-Duckdb/utils"

	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	// db connection bool
	Sqlitedb *sql.DB
	Dbstmt   map[string]*sql.Stmt
	Port     string
	Address  string
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// cread prepared statments to use in executing queries
func (s *Server) PrepareStmt() {
	var tmpStmt *sql.Stmt
	var err error
	for _, stmt := range preparedStmtStrings {
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

func NewServer() error {
	dbconn, err := sql.Open("sqlite3", os.Getenv("DBdir")+os.Getenv("ServerDbFile"))
	if err != nil {
		return err
	}
	err = dbconn.Ping()
	if err != nil {
		return err
	}

	Serv = &Server{
		Sqlitedb: dbconn,
		Port:     os.Getenv("ServerPort"),
		Dbstmt:   make(map[string]*sql.Stmt),
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	Serv.Address = os.Getenv("ServerAddr") + ":" + Serv.Port
	return nil
}

func Init() {
	if err := godotenv.Load("../.env"); err != nil {
		Serv.ErrorLog.Fatal(err)
	}

	if err := NewServer(); err != nil {
		panic(err)
	}

	Serv.CreateSuper()
	Serv.PrepareStmt()
}
