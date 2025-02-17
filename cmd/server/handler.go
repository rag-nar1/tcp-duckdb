package main

import (
	"TCP-Duckdb/internal"
	"TCP-Duckdb/utils"
	"database/sql"
	"fmt"
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
		
		if req[0] == "grant" {
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
	if (req[0] != "database" && req[0] != "table") || (req[0] == "database" && len(req) != 4) || (req[0] == "table" && len(req) != 5) {
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
		(*conn).Write([]byte("user: " + UserName + " does not have access over database: " + dbname + "\n"))
		return
	}

	buffer := make([]byte, 4096)

    _ , err = sql.Open("duckdb" , os.Getenv("DBdir") + "/users/" + dbname + ".db")
    if err != nil {
        (*conn).Write([]byte("SERVER ERROR\n"))
        log.Println(err)
        return
    }
    (*conn).Write([]byte("success\n"))
    
	for {
		n , err := (*conn).Read(buffer)
		if err != nil {
			(*conn).Write([]byte("ERROR: while reading\n"))
			log.Println("ERROR" , err)
			return
		}

        query := strings.ToLower(strings.Split(string(buffer[0:n]) , " ")[0])

        if query == "start" { // transaction
            
        }

        // single query
        QueryHandler(utils.Trim(string(buffer[0:n])), UserName, dbname, privilege, UID, DBID, conn)
        
	}

}

func QueryHandler(query, username, dbname, privilege string, UID, DBID int, conn *net.Conn) {
	query = strings.ToLower(query)
	fmt.Println(query)
	if privilege != "super" {
		hasaccess , err := internal.CheckAccesOverTable(server.Sqlitedb, server.dbstmt["CheckTableAccess"], query, UID, DBID)
		if err != nil || !hasaccess{
			(*conn).Write([]byte("Access denied\n"))
            log.Println(err)
			return
		}
	} else {
		hasDDL , err := internal.CheckDDLActions(query)
		if err != nil || !hasDDL {
			(*conn).Write([]byte("Access denied\n"))
            log.Println(err)
			return
		}
	}
		
	db , err := sql.Open("duckdb", os.Getenv("DBdir") + "/users/" + dbname + ".db")
	if err != nil {
		(*conn).Write([]byte("SERVER ERROR\n"))
		log.Println(err)
		return
	}

	if strings.HasPrefix(query, "select") {
        data, err := internal.SELECT(db, query)
        if err != nil {
            (*conn).Write([]byte("SERVER ERROR\n"))
            log.Println(err)
            return
        }
		data = append(data, '\n')
        (*conn).Write(data)
        return
	}

	if strings.HasPrefix(query, "Create") { 
		/*
			todo:
			1- create handler in internals 2- parse query create the table in sqlite database
		*/

	}

    // other statements
	LastInsertedID, RowsAffected, err := internal.EXEC(db, query)
	if err != nil {
		(*conn).Write([]byte("SERVER ERROR\n"))
		log.Println(err)
		return
	}

	data := []byte(strconv.Itoa(int(LastInsertedID)) + " " + strconv.Itoa(int(RowsAffected)) + "\n")
	(*conn).Write(data)
}

func GrantDB(dbname, username, accesstype string, conn *net.Conn) {
	accesstype = strings.ToLower(accesstype)
	// check for DB access
	if accesstype != "read" && accesstype != "write" {
		(*conn).Write([]byte("unsupported Access\n"))
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
        (*conn).Write([]byte("SERVER ERROR\n"))
        log.Println(err)
        return
    }

	(*conn).Write([]byte("success\n"))
}

func GrantTable(dbname, tablename, username, accesstype string, conn *net.Conn) {
	accesstype = strings.ToLower(accesstype)
	// check for DB access
	if accesstype != "select" && accesstype != "insert" && accesstype != "update" && accesstype != "delete"{
		(*conn).Write([]byte("unsupported Access\n"))
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
        (*conn).Write([]byte("SERVER ERROR\n"))
        log.Println(err)
        return
    }

	(*conn).Write([]byte("success\n"))
}