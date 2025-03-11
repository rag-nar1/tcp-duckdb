package connect

import (
	global "TCP-Duckdb/server"
	utils  "TCP-Duckdb/utils"
	"os"
	"bufio"
	"strings"
	"database/sql"

	// "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
)

// connect dbname
func Handler(server *global.Server, UID int, UserName, privilege, dbname string, reader *bufio.Reader, writer *bufio.Writer) {
	// check for db existense
	var DBID int
	if err := server.Dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID); err != nil {
		Error(writer, []byte("database: " + dbname + " does not exists\n"))
		return
	}

	// check for authrization
	var access int 
	if err := server.Dbstmt["CheckDbAccess"].QueryRow(UID , DBID).Scan(&access); err != nil {
		Error(writer, []byte("server error\n"))
		server.ErrorLog.Println(err)
		return
	}

	if access == 0 && privilege != "super" {
		Error(writer, []byte("user: " + UserName + " does not have access over database: " + dbname + "\n"))
		return
	}

	buffer := make([]byte, 4096)
    _ , err := sql.Open("duckdb" , os.Getenv("DBdir") + "/users/" + dbname + ".db")
    if err != nil {
        Error(writer, []byte("server error\n"))
        server.ErrorLog.Println(err)
        return
    }
    Write(writer, []byte("success\n"))
    
	for {
		n , err := reader.Read(buffer)
		if err != nil {
			Error(writer, []byte("while reading\n"))
			server.ErrorLog.Println(err)
			return
		}

        query := strings.ToLower(strings.Split(string(buffer[0:n]) , " ")[0])

        if query == "start" {
            if strings.ToLower(utils.Trim(string(buffer[0:n]))) != "start transaction" {
                Error(writer, []byte("bad request\n"))
                continue
            }
            Transaction(server, UID, DBID, dbname, privilege, reader, writer)
            continue
        }

        // single query
        QueryHandler(server, utils.Trim(string(buffer[0:n])), UserName, dbname, privilege, UID, DBID, writer)
	}

}

