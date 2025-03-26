package connect

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	
	"github.com/rag-nar1/TCP-Duckdb/internal"
	"github.com/rag-nar1/TCP-Duckdb/response"
	"github.com/rag-nar1/TCP-Duckdb/utils"
	"github.com/rag-nar1/TCP-Duckdb/pool"
	"github.com/rag-nar1/TCP-Duckdb/server"
)

func QueryService(Dbconn pool.Connection,query, username, dbname, privilege string, UID, DBID int, writer *bufio.Writer) {

	data, err := Query(query, privilege, UID, DBID, Dbconn.DB())
	if err != nil {
		server.Serv.ErrorLog.Println(err)
		if err.Error() == fmt.Errorf(response.UnauthorizedMSG).Error() {
			response.UnauthorizedError(writer)
			return
		}
		if err.Error() == fmt.Errorf(response.AccessDeniedMsg).Error() {
			response.AccesDeniedOverTables(writer, username, data)
			return
		}
		response.Error(writer, []byte(err.Error()))
		return
	}

	response.WriteData(writer, data)
}

func Transaction(Dbconn pool.Connection, UID, DBID int, username, dbname, privilege string, reader *bufio.Reader, writer *bufio.Writer) {
	buffer := make([]byte, 4096)
	transaction, err := Dbconn.DB().Begin()
	if err != nil {
		response.InternalError(writer)
		server.Serv.ErrorLog.Println(err)
		return
	}
	defer transaction.Rollback()

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			response.InternalError(writer)
			server.Serv.ErrorLog.Println("ERROR", err)
			return
		}

		query := strings.ToLower(utils.Trim(string(buffer[0:n])))
		if strings.HasPrefix(query, "rollback") {
			return
		}
		if strings.HasPrefix(query, "commit") {
			err = transaction.Commit()
			if err != nil {
				server.Serv.ErrorLog.Println(err)
				response.Error(writer, []byte(err.Error()))
			}
			return
		}

		data, err := Query(query, privilege, UID, DBID, transaction)
		if err != nil {
			server.Serv.ErrorLog.Println(err)
			if err.Error() == fmt.Errorf(response.UnauthorizedMSG).Error() {
				response.UnauthorizedError(writer)
				return
			}
			if err.Error() == fmt.Errorf(response.AccessDeniedMsg).Error() {
				response.AccesDeniedOverTables(writer, username, data)
				return
			}
			response.Error(writer, []byte(err.Error()))
			return
		}
		response.WriteData(writer, data)
	}

}

func Query(query, privilege string, UID, DBID int, executer internal.SQLExecutor) ([]byte, error) {
	query = strings.ToLower(query)
	authraized, err := Access(query, privilege, UID, DBID)
	if err != nil {
		return nil, err
	}
	if !authraized {
		tables, _ := internal.ExtractTableNames(query)

		return []byte(strings.Join(tables, ",")), fmt.Errorf(response.AccessDeniedMsg)
	}
	if strings.HasPrefix(query, "select") {
		data, err := internal.SELECT(executer, query)
		if err != nil {
			return nil, err
		}
		data = append(data, '\n')
		return data, nil
	}

	if strings.HasPrefix(query, "create") {
		err := internal.CREATE(executer, server.Serv.Sqlitedb, server.Serv.Dbstmt["CreateTable"], query, DBID)
		if err != nil {
			return nil, err
		}
		return []byte("Created"), nil
	}

	// other statements
	LastInsertedID, RowsAffected, err := internal.EXEC(executer, query)
	if err != nil {
		return nil, err
	}

	data := []byte(strconv.Itoa(int(LastInsertedID)) + " " + strconv.Itoa(int(RowsAffected)) + "\n")
	return data, nil
}

func Access(query, privilege string, UID, DBID int) (bool, error) {
	query = strings.ToLower(query)

	if privilege != "super" {
		hasDDL, err := internal.CheckDDLActions(query)
		if err != nil {
			return false, fmt.Errorf("%s", response.UnauthorizedMSG)
		}
		hasaccess, err := internal.CheckAccesOverTable(server.Serv.Sqlitedb, server.Serv.Dbstmt["CheckTableAccess"], query, UID, DBID)
		if err != nil {
			return false, fmt.Errorf("%s", response.UnauthorizedMSG)
		}
		return (hasaccess && !hasDDL), nil
	}
	return true, nil
}
