package grant


import (
	response	"TCP-Duckdb/response"
	global 		"TCP-Duckdb/server"
	utils 		"TCP-Duckdb/utils"

	"bufio"
	"strings"
)

func GrantDB(writer *bufio.Writer, dbname, username, accesstype string) {
	accesstype = strings.ToLower(accesstype)
	// check for DB access
	if accesstype != "read" && accesstype != "write" {
		response.BadRequest(writer)
		return
	}
	// get DBID , UID
	var DBID, UID int
	err := global.Serv.Dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err != nil {
		response.DoesNotExistDatabse(writer,dbname)
		return
	}

	err = global.Serv.Dbstmt["SelectUser"].QueryRow(username).Scan(&UID)
	if err != nil {
		response.DoesNotExistUser(writer, username)
		return
	}

	// grant access
	_, err = global.Serv.Dbstmt["GrantDB"].Exec(DBID, UID, accesstype)

	if err != nil {
        response.InternalError(writer)
        global.Serv.ErrorLog.Println(err)
        return
    }

	response.Success(writer)
}

func GrantTable(writer *bufio.Writer, dbname, tablename, username string, accesstypes ...string) {
	accesstypes = utils.ToLower(accesstypes...)
	// check for DB access
	for _, accesstype := range accesstypes {
		if accesstype != "select" && accesstype != "insert" && accesstype != "update" && accesstype != "delete"{
			response.BadRequest(writer)
			return
		}
	}

	var DBID, UID, TID int
	err := global.Serv.Dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err != nil {
		response.DoesNotExistDatabse(writer,dbname)
		return
	}

	err = global.Serv.Dbstmt["SelectUser"].QueryRow(username).Scan(&UID)
	if err != nil {
		response.DoesNotExistUser(writer, username)
		return
	}

	err = global.Serv.Dbstmt["SelectTable"].QueryRow(tablename, DBID).Scan(&TID)
	if err != nil {
		response.DoesNotExistTables(writer, tablename)
		return
	}

	transaction, err := global.Serv.Sqlitedb.Begin()
	if err != nil {
		response.InternalError(writer)
		global.Serv.ErrorLog.Println(err)
		return
	}
	defer transaction.Rollback()

	for _, accesstype := range accesstypes {
		_, err = transaction.Stmt(global.Serv.Dbstmt["GrantTable"]).Exec(TID, UID, accesstype)
		if err != nil {
			response.InternalError(writer)
			global.Serv.ErrorLog.Println(err)
			return
		}
	}
	
	response.Success(writer)
}