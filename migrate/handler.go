package migrate

import (
	internal "TCP-Duckdb/internal"
	response "TCP-Duckdb/response"
	global "TCP-Duckdb/server"
	utils "TCP-Duckdb/utils"

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
	
	if err := global.Serv.Dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID); err != nil {
		response.DoesNotExistDatabse(writer, dbname)
		global.Serv.ErrorLog.Println(err)
		return
	}

	var connStrEncrypted string
	
	if err := global.Serv.Dbstmt["SelectLink"].QueryRow(DBID).Scan(&connStrEncrypted); err != nil {
		response.Error(writer, []byte("database: " + dbname + " is not linked to any postgreSQL database\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}

	connStr, err := utils.Decrypt(connStrEncrypted, []byte(os.Getenv("ENCRYPTION_KEY")))
	if err != nil {
		response.InternalError(writer)
		global.Serv.ErrorLog.Println(err)
		return
	}
	// open duckdb
	duck, err := sqlx.Open("duckdb", os.Getenv("DBdir") + "/users/" + dbname + ".db")
	if err != nil {
		response.InternalError(writer)
		global.Serv.ErrorLog.Println(err)
		return
	}
	defer duck.Close()

	postgres, err := sqlx.Open("postgres", connStr)
	if err != nil {
		response.InternalError(writer)
		global.Serv.ErrorLog.Println(err)
		return
	}
	defer postgres.Close()

	
	if err := internal.ReadAudit(duck, postgres); err != nil {
		response.InternalError(writer)
		global.Serv.ErrorLog.Println(err)
		return
	}

	utils.Write(writer, []byte("migration is successful"))
}
