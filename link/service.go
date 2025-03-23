package link

import (
	internal 	"TCP-Duckdb/internal"
	response 	"TCP-Duckdb/response"
	global 		"TCP-Duckdb/server"
	utils 		"TCP-Duckdb/utils"

	"os"
	"bufio"
	"database/sql"

	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
)
func Link(writer *bufio.Writer, duck, postgres *sql.DB, connStr string, DBID int) {

	encryptedConnStr , err := utils.Encrypt(connStr, []byte(os.Getenv("ENCRYPTION_KEY")))
	if err != nil {
		response.InternalError(writer)
		global.Serv.ErrorLog.Println(err)
		return
	}

	// start a transaction to insert the connstr
	txServer, err := global.Serv.Sqlitedb.Begin()
	if err != nil {
		response.InternalError(writer)
		global.Serv.ErrorLog.Println(err)
		return
	}
	defer txServer.Rollback()

	txDuck, err := duck.Begin()
	if err != nil {
		response.InternalError(writer)
		global.Serv.ErrorLog.Println(err)
		return
	}
	defer txDuck.Rollback()

	txPg, err := postgres.Begin()
	if err != nil {
		response.InternalError(writer)
		global.Serv.ErrorLog.Println(err)
		return
	}
	defer txPg.Rollback()

	// insert the connstr
	_,err = txServer.Stmt(global.Serv.Dbstmt["CreateLink"]).Exec(DBID, encryptedConnStr)
	if err != nil {
		response.InternalError(writer)
		global.Serv.ErrorLog.Println(err)
		return
	}
	
	// migrate schema
	err = internal.Migrate(DBID, connStr, global.Serv.Dbstmt["CreateTable"], txPg, txDuck, txServer)
	if err != nil {
		global.Serv.ErrorLog.Println(err)
		response.InternalError(writer)
		return
	}

	err = internal.Audit(txPg)
	if err != nil {
		global.Serv.ErrorLog.Println(err)
		response.InternalError(writer)
		return
	}

	err = txPg.Commit()
	if err != nil {
		global.Serv.ErrorLog.Println(err)
		response.InternalError(writer)
		return
	}
	err = txDuck.Commit()
	if err != nil {
		global.Serv.ErrorLog.Println(err)
		response.InternalError(writer)
		return
	}
	err = txServer.Commit()
	if err != nil {
		global.Serv.ErrorLog.Println(err)
		response.InternalError(writer)
		return
	}

	response.Success(writer)
}