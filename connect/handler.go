package connect

import (
	"bufio"
	"database/sql"
	"os"
	"strings"

	response "github.com/rag-nar1/TCP-Duckdb/response"
	global "github.com/rag-nar1/TCP-Duckdb/server"
	utils "github.com/rag-nar1/TCP-Duckdb/utils"

	// "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
)

// connect dbname
func Handler(server *global.Server, UID int, UserName, privilege, dbname string, reader *bufio.Reader, writer *bufio.Writer) {
	// check for db existense
	var DBID int
	if err := server.Dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID); err != nil {
		response.DoesNotExistDatabse(writer, dbname)
		return
	}

	// check for authrization
	var access int
	if err := server.Dbstmt["CheckDbAccess"].QueryRow(UID, DBID).Scan(&access); err != nil {
		response.InternalError(writer)
		server.ErrorLog.Println(err)
		return
	}

	if access == 0 && privilege != "super" {
		response.AccesDeniedOverDatabase(writer, UserName, dbname)
		return
	}

	buffer := make([]byte, 4096)
	_, err := sql.Open("duckdb", os.Getenv("DBdir")+"/users/"+dbname+".db")
	if err != nil {
		response.InternalError(writer)
		server.ErrorLog.Println(err)
		return
	}
	response.Success(writer)

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			response.InternalError(writer)
			server.ErrorLog.Println(err)
			return
		}

		query := strings.ToLower(strings.Split(string(buffer[0:n]), " ")[0])

		if query == "start" {
			if strings.ToLower(utils.Trim(string(buffer[0:n]))) != "start transaction" {
				response.BadRequest(writer)
				continue
			}
			Transaction(server, UID, DBID, UserName, dbname, privilege, reader, writer)
			continue
		}

		// single query
		QueryService(server, utils.Trim(string(buffer[0:n])), UserName, dbname, privilege, UID, DBID, writer)
	}

}
