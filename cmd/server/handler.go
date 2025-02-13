package main

import (
	"TCP-Duckdb/utils"
	"database/sql"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	_ "github.com/marcboeker/go-duckdb"
)

func LoginHandler(UserName, password string, dbstmt *sql.Stmt) (int , string , error){
	
	var UID int
	var privilige string
	err := dbstmt.QueryRow(UserName , utils.Hash(password)).Scan(&UID, &privilige)
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
			(*conn).Write([]byte("ERROR: while reading\n"))
			log.Println("ERROR" , err)
			return
		}


		req := strings.Split(string(rawreq[0:n]) , " ")
		utils.TrimList(req)

		if req[0] != "connect" && req[0] != "create" && req[0] != "grant" && req[0] != "migrate" {
			(*conn).Write([]byte("ERROR: BAD request\n"))
			continue
		}
		
		if req[0] == "connect" {
			if len(req) != 2 {
				(*conn).Write([]byte("ERROR: BAD request\n"))
				continue
			}
			DbConnectionHandler(UID, UserName, privilege, req[1], conn) 
			continue
		}

		if req[0] == "create" {
			CreatHandler(privilege, req[1:], conn)
			continue
		}

		if req[1] == "grant" {
			GrantHandler(privilege, req[1:], conn)
			continue
		}

		MigrateHandler()

	}
	
}

// create database [dbname],
// create user [dbname] [username] [password]
func CreatHandler(privilege string, req []string, conn *net.Conn) {
	if privilege != "super" {
		(*conn).Write([]byte("Unauthorized\n"))
		return
	}

	if (req[0] != "database" && req[0] != "user") || (req[0] == "database" && len(req) != 2) {
		(*conn).Write([]byte("ERROR: BAD request\n"))
		return
	}

	if req[1] == "database" {
		CreateDB(req[2], conn)
		return
	}

	if len(req) != 4 {
		(*conn).Write([]byte("ERROR: BAD request\n"))
		return
	}

	CreateUser(req[2], req[3], req[4], conn) 
}
// grant database [dbname] [username] [accesstype] ,
// grant table [dbname] [tablename] [username] [accesstype] 
func GrantHandler(privilege string, req []string, conn *net.Conn) { 
	if privilege != "super" {
		(*conn).Write([]byte("Unauthorized\n"))
		return
	}

	if (req[0] != "database" && req[0] != "user") || (req[0] == "database" && len(req) != 4) || (req[0] == "table" && len(req) != 5) {
		(*conn).Write([]byte("ERROR: BAD request\n"))
		return
	}

	if req[0] == "database" {
		GrantDB(req[1], req[2], req[3], conn)
		return
	}

	GrantTable(req[1], req[2], req[3], req[4], conn)
}

func MigrateHandler() { // todo

}

func CreateDB(dbname string, conn *net.Conn) {
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
		(*conn).Write([]byte("ERROR: could not create databse\n"))

		err = os.Remove(os.Getenv("DBdir") + "/users/" + dbname + ".db")
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	(*conn).Write([]byte("success\n"))
}

func CreateUser(dbname, NewUser, password string, conn *net.Conn) {
	var DBID, UID int
	err := server.dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err != nil {
		(*conn).Write([]byte("database: " + dbname + " does not exists\n"))
		return
	}

	//start transaction
	transaction, err := server.Sqlitedb.Begin()
	defer transaction.Rollback()

	// create user
	_, err = transaction.Stmt(server.dbstmt["CreateUser"]).Exec(NewUser, utils.Hash(password) , "norm")

	if err != nil {
		(*conn).Write([]byte("user already exists\n"))
		log.Println(err)
		return
	}

	var privilege string
	err = transaction.Stmt(server.dbstmt["login"]).QueryRow(NewUser , utils.Hash(password)).Scan(&UID , &privilege)
	if err != nil {
		(*conn).Write([]byte("server error\n"))
		log.Println(err)
		return
	}

	_, err = transaction.Stmt(server.dbstmt["GrantDB"]).Exec(DBID, UID, "read")
	if err != nil {
		(*conn).Write([]byte("server error\n"))
		log.Println(err)
		return
	}
	
	err = transaction.Commit()
	if err != nil {
		(*conn).Write([]byte("server error\n"))
		log.Println(err)
		return
	}
	(*conn).Write([]byte("success\n"))
}

func DbConnectionHandler(UID int, UserName, privilege, dbname string, conn *net.Conn) {
	// check for db existense
	var DBID int
	err := server.dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err != nil {
		(*conn).Write([]byte("database: " + dbname + " does not exists\n"))
		return
	}

	// check for authrization
	var access int 
	err = server.dbstmt["CheckDbAccess"].QueryRow(UID , DBID).Scan(&access)
	if err != nil {
		(*conn).Write([]byte("server error\n"))
		log.Println(err)
		return
	}

	if access == 0 {
		(*conn).Write([]byte("user: " + UserName + " does not have access over database: " + dbname))
		return
	}

	buffer := make([]byte, 4096)

    _ , err = sql.Open("duckdb" , os.Getenv("DBdir") + "/users/" + dbname + ".db")
    if err != nil {
        (*conn).Write([]byte("SERVER ERROR"))
        log.Println(err)
        return
    }
    (*conn).Write([]byte("success"))
	for {
		n , err := (*conn).Read(buffer)
		if err != nil {
			(*conn).Write([]byte("ERROR: while reading\n"))
			log.Println("ERROR" , err)
			return
		}

        query := strings.Split(string(buffer[0:n]) , " ")[0]

        if query != "SELECT" && query != "TRANSACTION" {
            (*conn).Write([]byte("unsupported Query"))
            continue
        }

        if query == "SELECT" {
            // todo: add selecet query and design the response design
        }
        
        
        
	}

}

func GrantDB(dbname, username, accesstype string, conn *net.Conn) {
	accesstype = strings.ToLower(accesstype)
	// check for DB access
	if accesstype != "read" && accesstype != "write" {
		(*conn).Write([]byte("unsupported Access"))
		return
	}
	// get DBID , UID
	var DBID, UID int
	err := server.dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err != nil {
		(*conn).Write([]byte("database: " + dbname + " does not exists\n"))
		return
	}

	err = server.dbstmt["SelectUser"].QueryRow(username).Scan(&UID)
	if err != nil {
		(*conn).Write([]byte("user: " + username + " does not exists\n"))
		return
	}

	// grant access

	_, err = server.dbstmt["GrantDB"].Exec(DBID, UID, accesstype)

	if err != nil {
        (*conn).Write([]byte("SERVER ERROR"))
        log.Println(err)
        return
    }

	(*conn).Write([]byte("success"))
}

func GrantTable(dbname, tablename, username, accesstype string, conn *net.Conn) {
	accesstype = strings.ToLower(accesstype)
	// check for DB access
	if accesstype != "select" && accesstype != "insert" && accesstype != "update" && accesstype != "delete"{
		(*conn).Write([]byte("unsupported Access"))
		return
	}

	var DBID, UID, TID int
	err := server.dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err != nil {
		(*conn).Write([]byte("database: " + dbname + " does not exists\n"))
		return
	}

	err = server.dbstmt["SelectUser"].QueryRow(username).Scan(&UID)
	if err != nil {
		(*conn).Write([]byte("user: " + username + " does not exists\n"))
		return
	}

	err = server.dbstmt["SelectTable"].QueryRow(tablename, DBID).Scan(&TID)
	if err != nil {
		(*conn).Write([]byte("table: " + tablename + " does not exists in database: " + strconv.Itoa(DBID) + "\n"))
		return
	}
	
	_, err = server.dbstmt["GrantTable"].Exec(TID, UID, accesstype)
	if err != nil {
        (*conn).Write([]byte("SERVER ERROR"))
        log.Println(err)
        return
    }

	(*conn).Write([]byte("success"))
}