package connect

import (
	"bufio"

	"strings"

	"github.com/rag-nar1/tcp-duckdb/request_handler"
	"github.com/rag-nar1/tcp-duckdb/response"
	"github.com/rag-nar1/tcp-duckdb/server"
	"github.com/rag-nar1/tcp-duckdb/utils"
)

// connect dbname
func Handler(UID int, UserName, privilege, dbname string, reader *bufio.Reader, writer *bufio.Writer) {
	// check for db existense
	var DBID int
	if err := server.Serv.Dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID); err != nil {
		response.DoesNotExistDatabse(writer, dbname)
		return
	}

	// check for authrization
	var access int
	if err := server.Serv.Dbstmt["CheckDbAccess"].QueryRow(UID, DBID).Scan(&access); err != nil {
		response.InternalError(writer)
		server.Serv.ErrorLog.Println(err)
		return
	}

	if access == 0 && privilege != "super" {
		response.AccesDeniedOverDatabase(writer, UserName, dbname)
		return
	}

	req := request_handler.NewRequest(dbname)
	server.Serv.Pool.Push(req)

	DbConn, err := utils.OpenDb(server.Serv.Pool, dbname)
	if err != nil {
		response.InternalError(writer)
		server.Serv.ErrorLog.Println(err)
		return
	}
	defer DbConn.Destroy()
	
	response.Success(writer)
	
	buffer := make([]byte, 4096)
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			response.InternalError(writer)
			server.Serv.ErrorLog.Println(err)
			return
		}

		query := strings.ToLower(strings.Split(string(buffer[0:n]), " ")[0])

		if query == "start" {
			if strings.ToLower(utils.Trim(string(buffer[0:n]))) != "start transaction" {
				response.BadRequest(writer)
				continue
			}
			Transaction(DbConn, UID, DBID, UserName, dbname, privilege, reader, writer)
			continue
		}

		// single query
		QueryService(DbConn, utils.Trim(string(buffer[0:n])), UserName, dbname, privilege, UID, DBID, writer)
	}
}
