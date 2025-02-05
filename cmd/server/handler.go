package main

import (
	"TCP-Duckdb/utils"
	"database/sql"
	"log"
	"net"
	"os"
	"strings"
	_ "github.com/marcboeker/go-duckdb"
)

func LoginHandler(UserName, password string, dbstmt *sql.Stmt) (int , string , error){
	
	var UID int
	var privilige string
	err := dbstmt.QueryRow(UserName , password).Scan(&UID, &privilige)
	if err != nil {
		return -1 , "", err
	}

	return UID, privilige, nil
}

func DBHandler(UID int, UserName, privilege string, conn *net.Conn) {
	log.Println("Serving: " + UserName)
	rawreq := make([]byte, 1024)
	for {
		n , err := (*conn).Read(rawreq)
		if err != nil {
			(*conn).Write([]byte("ERROR: while reading"))
			log.Println("ERROR" , err)
			return
		}


		req := strings.Split(string(rawreq[0:n]) , " ")
		utils.TrimList(req)

		if req[0] != "connect" && req[0] != "create" {
			(*conn).Write([]byte("ERROR: BAD request\n"))
			continue
		}

		if req[0] == "connect" && len(req) != 2 {
			(*conn).Write([]byte("ERROR: BAD request\n"))
			continue
		}

		if req[0] == "connect" {
			// DbConnectionHandler(UID, UserName, privilege, dbname, conn) // TODO
			continue
		}

		if req[0] == "create" && privilege != "super" {
			(*conn).Write([]byte("Unauthorized\n"))
			continue
		}

		if req[1] != "database" && req[1] != "user" {
			(*conn).Write([]byte("ERROR: BAD request\n"))
			continue
		}

		if req[1] == "database" && len(req) != 3 {
			(*conn).Write([]byte("ERROR: BAD request\n"))
			continue
		}

		if req[1] == "database" {
			CreateDB(UID, UserName, privilege, req[2], conn)
			continue
		}

		// CreateUser(UID, UserName, privilege, NewUser, dbname, conn) // TODO
	}
	
}

func CreateDB(UID int, UserName, privilege, dbname string, conn *net.Conn) {
	var DBID int
	err := server.dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err == nil {
		(*conn).Write([]byte("database: " + dbname + " already exists\n"))
		return
	}
	
	// create file
	_ ,err = sql.Open("duckdb", os.Getenv("DBdir") + "/users/" + dbname + ".db")
	if err != nil {
		log.Fatal(err)
	}

	_ , err = server.dbstmt["CreateDB"].Exec(dbname)
	if err != nil {
		(*conn).Write([]byte("ERROR: could not create databse"))

		err = os.Remove(os.Getenv("DBdir") + "/users/" + dbname + ".db")
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	(*conn).Write([]byte("success\n"))
}