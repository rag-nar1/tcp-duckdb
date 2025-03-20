package link

import (
	response 	"TCP-Duckdb/response"
	global 		"TCP-Duckdb/server"

	"bufio"
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
)

// link dbname connStr
func Handler(privilege string, req []string, writer *bufio.Writer) {
	if privilege != "super" {
		response.UnauthorizedError(writer)
		return
	}

	dbname, connStr := req[0], req[1]
	// check the existince of the database
	var DBID int
	err := global.Serv.Dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err != nil {
		response.DoesNotExistDatabse(writer, dbname)
		global.Serv.ErrorLog.Println(err)
		return
	}

	var hasLink int
	err = global.Serv.Dbstmt["CheckLink"].QueryRow(DBID).Scan(&hasLink)
	if err != nil || hasLink > 0 {
		response.Error(writer, []byte("database: " + dbname + " already linked"))
		global.Serv.ErrorLog.Println(err)
		return
	}

	// open duckdb
	duck, err := sql.Open("duckdb", os.Getenv("DBdir") + "/users/" + dbname + ".db")
	if err != nil {
		response.InternalError(writer)
		global.Serv.ErrorLog.Println(err)
		return
	}
	defer duck.Close()

	// check the connStr
	postgres, err := sql.Open("postgres", connStr)
	if err != nil {
		response.InternalError(writer)
		global.Serv.ErrorLog.Println(err)
		return
	}
	defer postgres.Close()

	err = postgres.Ping()
	if err != nil {
		response.InternalError(writer)
		global.Serv.ErrorLog.Println(err)
		return
	}

	
	
}
