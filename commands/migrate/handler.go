package commands

import (
	"github.com/rag-nar1/tcp-duckdb/internal"
	"github.com/rag-nar1/tcp-duckdb/response"
	"github.com/rag-nar1/tcp-duckdb/server"
	"github.com/rag-nar1/tcp-duckdb/utils"

	"bufio"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
)

func Handler(privilege string, req []string, writer *bufio.Writer) {
	if privilege != "super" {
		response.UnauthorizedError(writer)
		return
	}

	dbname := req[0]
	// check the existince of the database
	var DBID int

	if err := server.Serv.Dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID); err != nil {
		response.DoesNotExistDatabse(writer, dbname)
		server.Serv.ErrorLog.Println(err)
		return
	}

	var connStrEncrypted string

	if err := server.Serv.Dbstmt["SelectLink"].QueryRow(DBID).Scan(&connStrEncrypted); err != nil {
		response.Error(writer, []byte("database: "+dbname+" is not linked to any postgreSQL database\n"))
		server.Serv.ErrorLog.Println(err)
		return
	}

	connStr, err := utils.Decrypt(connStrEncrypted, []byte(os.Getenv("ENCRYPTION_KEY")))
	if err != nil {
		response.InternalError(writer)
		server.Serv.ErrorLog.Println(err)
		return
	}

	// open duckdb
	Dbconn, err := utils.OpenDb(server.Serv.Pool, dbname)
	if err != nil {
		response.InternalError(writer)
		server.Serv.ErrorLog.Println(err)
		return
	}
	defer Dbconn.Destroy()

	duck := sqlx.NewDb(Dbconn.DB(), "duckdb")
	postgres, err := sqlx.Open("postgres", connStr)
	if err != nil {
		response.InternalError(writer)
		server.Serv.ErrorLog.Println(err)
		return
	}
	defer postgres.Close()

	if err := internal.ReadAudit(duck, postgres); err != nil {
		response.InternalError(writer)
		server.Serv.ErrorLog.Println(err)
		return
	}

	utils.Write(writer, []byte("migration is successful"))
}
