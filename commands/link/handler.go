package commands

import (
	"bufio"
	"database/sql"
	
	"github.com/rag-nar1/tcp-duckdb/response"
	"github.com/rag-nar1/tcp-duckdb/server"
	"github.com/rag-nar1/tcp-duckdb/utils"
	
	_ "github.com/lib/pq"
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
	err := server.Serv.Dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err != nil {
		response.DoesNotExistDatabse(writer, dbname)
		server.Serv.ErrorLog.Println(err)
		return
	}

	var hasLink int
	err = server.Serv.Dbstmt["CheckLink"].QueryRow(DBID).Scan(&hasLink)
	if err != nil || hasLink > 0 {
		response.Error(writer, []byte("database: "+dbname+" already linked"))
		server.Serv.ErrorLog.Println(err)
		return
	}

	DbConn, err := utils.OpenDb(server.Serv.Pool, dbname)
	if err != nil {
		response.InternalError(writer)
		server.Serv.ErrorLog.Println(err)
		return
	}
	defer DbConn.Destroy()

	// check the connStr
	postgres, err := sql.Open("postgres", connStr)
	if err != nil {
		response.InternalError(writer)
		server.Serv.ErrorLog.Println(err)
		return
	}
	defer postgres.Close()

	Link(writer, DbConn.DB(), postgres, connStr, DBID)
}
