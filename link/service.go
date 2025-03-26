package link

import (
	"bufio"
	"database/sql"
	"os"
	
	"github.com/rag-nar1/TCP-Duckdb/internal"
	"github.com/rag-nar1/TCP-Duckdb/response"
	"github.com/rag-nar1/TCP-Duckdb/server"
	"github.com/rag-nar1/TCP-Duckdb/utils"
)

func Link(writer *bufio.Writer, duck, postgres *sql.DB, connStr string, DBID int) {

	encryptedConnStr, err := utils.Encrypt(connStr, []byte(os.Getenv("ENCRYPTION_KEY")))
	if err != nil {
		response.InternalError(writer)
		server.Serv.ErrorLog.Println(err)
		return
	}

	// start a transaction to insert the connstr
	txServer, err := server.Serv.Sqlitedb.Begin()
	if err != nil {
		response.InternalError(writer)
		server.Serv.ErrorLog.Println(err)
		return
	}
	defer txServer.Rollback()

	txDuck, err := duck.Begin()
	if err != nil {
		response.InternalError(writer)
		server.Serv.ErrorLog.Println(err)
		return
	}
	defer txDuck.Rollback()

	txPg, err := postgres.Begin()
	if err != nil {
		response.InternalError(writer)
		server.Serv.ErrorLog.Println(err)
		return
	}
	defer txPg.Rollback()

	// insert the connstr
	_, err = txServer.Stmt(server.Serv.Dbstmt["CreateLink"]).Exec(DBID, encryptedConnStr)
	if err != nil {
		response.InternalError(writer)
		server.Serv.ErrorLog.Println(err)
		return
	}

	// migrate schema
	err = internal.Migrate(DBID, connStr, server.Serv.Dbstmt["CreateTable"], txPg, txDuck, txServer)
	if err != nil {
		server.Serv.ErrorLog.Println(err)
		response.InternalError(writer)
		return
	}

	err = internal.Audit(txPg)
	if err != nil {
		server.Serv.ErrorLog.Println(err)
		response.InternalError(writer)
		return
	}

	err = txPg.Commit()
	if err != nil {
		server.Serv.ErrorLog.Println(err)
		response.InternalError(writer)
		return
	}
	err = txDuck.Commit()
	if err != nil {
		server.Serv.ErrorLog.Println(err)
		response.InternalError(writer)
		return
	}
	err = txServer.Commit()
	if err != nil {
		server.Serv.ErrorLog.Println(err)
		response.InternalError(writer)
		return
	}

	response.Success(writer)
}
