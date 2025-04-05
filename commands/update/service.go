package commands

import (
	"bufio"
	"os"

	"github.com/rag-nar1/tcp-duckdb/response"
	"github.com/rag-nar1/tcp-duckdb/server"
	"github.com/rag-nar1/tcp-duckdb/utils"
)



func UpdateDatabase(writer *bufio.Writer, oldDbname, newDbName string) {
	var DBID int
	if err := server.Serv.Dbstmt["SelectDB"].QueryRow(oldDbname).Scan(&DBID); err != nil {
		response.DoesNotExistDatabse(writer, oldDbname)
		server.Serv.ErrorLog.Println(err)
		return
	}

	if err := os.Rename(utils.UserDbPath(oldDbname), utils.UserDbPath(newDbName)); err != nil {
		response.InternalError(writer)
		server.Serv.ErrorLog.Println(err)
		return
	}

	if _, err := server.Serv.Dbstmt["UpdateDB"].Exec(newDbName, DBID); err != nil {
		if err := os.Rename(utils.UserDbPath(newDbName), utils.UserDbPath(oldDbname)); err != nil {
			response.InternalError(writer)
			server.Serv.ErrorLog.Println(err)
		}
		response.InternalError(writer)
		server.Serv.ErrorLog.Println(err)
		return
	}

	response.Success(writer)
}

func UpdateUserUsername(writer *bufio.Writer, oldUserName, NewUserName string) {
	var UID int
	if err := server.Serv.Dbstmt["SelectUser"].QueryRow(oldUserName).Scan(&UID); err != nil {
		response.DoesNotExistUser(writer, oldUserName)
		server.Serv.ErrorLog.Println(err)
		return
	}

	if _, err := server.Serv.Dbstmt["UpdateUsername"].Exec(NewUserName, UID); err != nil {
		response.InternalError(writer)
		server.Serv.ErrorLog.Println(err)
		return
	}
	response.Success(writer)
}

func UpdateUserPassword(writer *bufio.Writer, UserName, password string) {
	var UID int
	if err := server.Serv.Dbstmt["SelectUser"].QueryRow(UserName).Scan(&UID); err != nil {
		response.DoesNotExistUser(writer, UserName)
		server.Serv.ErrorLog.Println(err)
		return
	}

	password = utils.Hash(password)
	if _, err := server.Serv.Dbstmt["UpdateUserPassword"].Exec(password, UID); err != nil {
		response.InternalError(writer)
		server.Serv.ErrorLog.Println(err)
		return
	}
	response.Success(writer)
}