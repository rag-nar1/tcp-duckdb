package create

import (
	utils		"TCP-Duckdb/utils"
	global 		"TCP-Duckdb/server"

	"os"
	"bufio"
	"database/sql"

	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
)

// create database [dbname],
// create user [dbname] [username] [password]
func Handler(privilege string, req []string, writer *bufio.Writer) {
	if privilege != "super" {
		Error(writer, []byte("Unauthorized\n"))
		return
	}

	if (req[0] != "database" && req[0] != "user") || (req[0] == "database" && len(req) != 2) {
		Error(writer, []byte("BAD request\n"))
		return
	}

	if req[0] == "database" {
		CreateDB(req[1], writer)
		return
	}

	if len(req) != 4 {
		Error(writer, []byte("BAD request\n"))
		return
	}

	CreateUser(req[2], req[3], req[4], writer) 
}

func CreateDB(dbname string, writer *bufio.Writer) {
	var DBID int
	err := global.Serv.Dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err == nil && DBID != 0 {
		utils.Write(writer, []byte("database: " + dbname + " already exists\n"))
		return
	}
	
	// create file
	_ ,err = sql.Open("duckdb", os.Getenv("DBdir") + "/users/" + dbname + ".db")
	if err != nil {
		global.Serv.ErrorLog.Fatal(err)
	}

	_ , err = global.Serv.Dbstmt["CreateDB"].Exec(dbname)
	if err != nil {
		utils.Write(writer, []byte("ERROR: could not create databse\n"))

		err = os.Remove(os.Getenv("DBdir") + "/users/" + dbname + ".db")
		if err != nil {
			global.Serv.ErrorLog.Fatal(err)
		}
		return
	}
	utils.Write(writer, []byte("success\n"))
}

func CreateUser(dbname, NewUser, password string, writer *bufio.Writer) {
	var DBID, UID int
	err := global.Serv.Dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err != nil {
		utils.Write(writer, []byte("database: " + dbname + " does not exists\n"))
		return
	}

	//start transaction
	transaction, err := global.Serv.Sqlitedb.Begin()
	defer transaction.Rollback()

	// create user
	_, err = transaction.Stmt(global.Serv.Dbstmt["CreateUser"]).Exec(NewUser, utils.Hash(password) , "norm")

	if err != nil {
		utils.Write(writer, []byte("user already exists\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}

	var privilege string
	err = transaction.Stmt(global.Serv.Dbstmt["login"]).QueryRow(NewUser , utils.Hash(password)).Scan(&UID , &privilege)
	if err != nil {
		utils.Write(writer, []byte("global.Serv error\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}

	_, err = transaction.Stmt(global.Serv.Dbstmt["GrantDB"]).Exec(DBID, UID, "read")
	if err != nil {
		utils.Write(writer, []byte("global.Serv error\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}
	
	err = transaction.Commit()
	if err != nil {
		utils.Write(writer, []byte("global.Serv error\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}
	utils.Write(writer, []byte("success\n"))
}