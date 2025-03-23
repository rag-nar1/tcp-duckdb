package create

import (
	response "github.com/rag-nar1/TCP-Duckdb/response"
	global "github.com/rag-nar1/TCP-Duckdb/server"
	utils "github.com/rag-nar1/TCP-Duckdb/utils"

	"bufio"
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
)

func CreateDB(dbname string, writer *bufio.Writer) {
	var DBID int
	err := global.Serv.Dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err == nil && DBID != 0 {
		response.Error(writer, []byte("database: "+dbname+" already exists"))
		return
	}

	// create file
	_, err = sql.Open("duckdb", os.Getenv("DBdir")+"/users/"+dbname+".db")
	if err != nil {
		global.Serv.ErrorLog.Fatal(err)
	}

	_, err = global.Serv.Dbstmt["CreateDB"].Exec(dbname)
	if err != nil {
		response.Error(writer, []byte("could not create databse"))

		err = os.Remove(os.Getenv("DBdir") + "/users/" + dbname + ".db")
		if err != nil {
			global.Serv.ErrorLog.Fatal(err)
		}
		return
	}
	response.Success(writer)
}

func CreateUser(NewUser, password string, writer *bufio.Writer) {

	// create user
	_, err := global.Serv.Dbstmt["CreateUser"].Exec(NewUser, utils.Hash(password), "norm")

	if err != nil {
		response.Error(writer, []byte("user already exists\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}

	response.Success(writer)
}
