package connect

import (
	internal "TCP-Duckdb/internal"
	global "TCP-Duckdb/server"
	utils "TCP-Duckdb/utils"
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
)

func QueryHandler(server *global.Server, query, username, dbname, privilege string, UID, DBID int, writer *bufio.Writer) {
	db , err := sql.Open("duckdb", os.Getenv("DBdir") + "/users/" + dbname + ".db")
	if err != nil {
		Error(writer, []byte("server error\n"))
		server.ErrorLog.Println(err)
		return
	}

	data, err := Query(server, query, privilege, UID, DBID, db)
	if err != nil{
		server.ErrorLog.Println(err)
		Error(writer, []byte("Error while executing query\n"))
		return
	}

	Write(writer, data)

}

func Transaction(server *global.Server, UID, DBID int, dbname, privilege string, reader *bufio.Reader, writer *bufio.Writer) {
    buffer := make([]byte, 4096)
    db , err := sql.Open("duckdb", os.Getenv("DBdir") + "/users/" + dbname + ".db")
    if err != nil {
		Error(writer, []byte("server error\n"))
		server.ErrorLog.Println(err)
		return
	}

    transaction, err := db.Begin()
    if err != nil {
		Error(writer, []byte("server error\n"))
		server.ErrorLog.Println(err)
		return
	}
    defer transaction.Rollback()

    for {
        n , err := reader.Read(buffer)
		if err != nil {
			Error(writer, []byte("while reading\n"))
			server.ErrorLog.Println("ERROR" , err)
			return
		}

        query := strings.ToLower(utils.Trim(string(buffer[0:n])))
        if strings.HasPrefix(query, "rollback") {
            return
        }
        if strings.HasPrefix(query, "commit") {
            err = transaction.Commit()
            if err != nil {
                server.ErrorLog.Println(err)
                Error(writer, []byte("Error while commiting transaction\n"))
            }
            return
        }
        
        data, err := Query(server, query, privilege, UID, DBID, transaction)
        if err != nil{
			server.ErrorLog.Println(err)
			Error(writer, []byte("Error while executing query\n"))
			return
        }

		Write(writer, data)
    }

}

func Query(server *global.Server, query, privilege string, UID, DBID int, executer internal.SQLExecutor)  ([]byte, error) {
	query = strings.ToLower(query)
    authraized, err := Access(server, query, privilege, UID, DBID)
	if err != nil {
        return nil, err
	}
    if !authraized {
        return nil, fmt.Errorf("Access denied\n")
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
		err := internal.CREATE(executer, server.Sqlitedb, server.Dbstmt["CreateTable"], query, DBID)
		if err != nil {
			return nil, err 
		}
        return []byte("Created"),nil
	}

    // other statements
	LastInsertedID, RowsAffected, err := internal.EXEC(executer, query)
	if err != nil {
		return nil, err 
	}

	data := []byte(strconv.Itoa(int(LastInsertedID)) + " " + strconv.Itoa(int(RowsAffected)) + "\n")
	return data, nil
}

func Access(server *global.Server, query, privilege string, UID, DBID int) (bool, error){
    query = strings.ToLower(query)
	
	if privilege != "super" {
		hasDDL , err := internal.CheckDDLActions(query)
		if err != nil {
			return false, err
		}
		hasaccess , err := internal.CheckAccesOverTable(server.Sqlitedb, server.Dbstmt["CheckTableAccess"], query, UID, DBID)
		if err != nil {
			return false, err
		}
        return (hasaccess && !hasDDL), nil
	}
    return true, nil
}
