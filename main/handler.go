package main

import (
	utils		"TCP-Duckdb/utils"
	global 		"TCP-Duckdb/server"
	connect 	"TCP-Duckdb/connect" 
	internal 	"TCP-Duckdb/internal"

	"os"
	"net"
	"bufio"
	"strconv"
	"strings"
	"crypto/rand"
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	global.Serv.InfoLog.Println("Serving " + conn.RemoteAddr().String())
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	// read login request
	route := make([]byte, 1024)
	n , err := reader.Read(route)
	if err != nil {
		global.Serv.ErrorLog.Println(err)
		return
	}
	// check for a valid request
	request := strings.Split(string(route[0 : n]) , " ")
	var UserName , password , privilege string
	var UID int
	if request[0] != "login" || len(request) != 3 {
		utils.Write(writer, []byte("ERROR: BAD request\n"))
		return
	}
	// validate the username and password
	UserName, password = utils.Trim(request[1]), utils.Trim(request[2]) 
	UID, privilege, err = LoginHandler(UserName, password, global.Serv.Dbstmt["login"])
	if err != nil {
		utils.Write(writer, []byte("Unauthorized\n"))
		global.Serv.ErrorLog.Println(err)
		return	
	}
	utils.Write(writer, []byte("success\n"))
	DBHandler(UID, UserName, privilege, reader, writer)
}


func LoginHandler(UserName, password string, Dbstmt *sql.Stmt) (int , string , error){
	
	var UID int
	var privilige string
	err := Dbstmt.QueryRow(UserName , utils.Hash(password)).Scan(&UID, &privilige)
	if err != nil {
		return -1 , "", err
	}

	return UID, privilige, nil
}

func DBHandler(UID int, UserName, privilege string, reader *bufio.Reader, writer *bufio.Writer) {
	global.Serv.InfoLog.Println("Serving: " + UserName)
	rawreq := make([]byte, 1024)
	for {
		n , err := reader.Read(rawreq)
		if err != nil {
			utils.Write(writer, []byte("ERROR: while reading\n"))
			global.Serv.ErrorLog.Println(err)
			return
		}

		req := strings.Split(string(rawreq[0:n]) , " ")
		utils.TrimList(req)
		
		if req[0] != "connect" && req[0] != "create" && req[0] != "grant" && req[0] != "migrate" && req[0] != "link" {
			utils.Write(writer, []byte("ERROR: BAD request\n"))
			continue
		}
		
		if req[0] == "connect" {
			if len(req) != 2 {
				utils.Write(writer, []byte("ERROR: BAD request\n"))
				continue
			}
			connect.Handler(global.Serv, UID, UserName, privilege, req[1], reader, writer) 
			continue
		}
		
		if req[0] == "create" {
			CreatHandler(privilege, req[1:], writer)
			continue
		}
		
		if req[0] == "grant" {
			GrantHandler(privilege, req[1:], writer)
			continue
		}

		if req[0] == "link" {
			LinkHandler(privilege, req[1:], writer)
		}

		MigrateHandler(privilege, req[1:], writer)

	}
	
}
// link dbname connStr
func LinkHandler(privilege string, req []string, writer *bufio.Writer) {
	if privilege != "super" {
		utils.Write(writer, []byte("Unauthorized\n"))
		return
	}

	dbname, connStr := req[0], req[1]
	// check the existince of the database
	var DBID int
	err := global.Serv.Dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err != nil {
		utils.Write(writer, []byte("database: " + dbname + " does not exists\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}
	var hasLink int
	err = global.Serv.Dbstmt["CheckLink"].QueryRow(DBID).Scan(&hasLink)
	if err != nil || hasLink > 0 {
		utils.Write(writer, []byte("database: " + dbname + " already linked\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}

	// open duckdb
	duck, err := sql.Open("duckdb", os.Getenv("DBdir") + "/users/" + dbname + ".db")
	if err != nil {
		utils.Write(writer, []byte("error while connecting to the duckdb database\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}
	defer duck.Close()

	// check the connStr
	postgres, err := sql.Open("postgres", connStr)
	if err != nil {
		utils.Write(writer, []byte("error while connecting to the postgresql database\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}
	defer postgres.Close()

	err = postgres.Ping()
	if err != nil {
		utils.Write(writer, []byte("error while connecting to the postgresql database\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}

	// generete random 32 byte key for encryption
	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		utils.Write(writer, []byte("global.Serv Error\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}

	encryptedConnStr , err := utils.Encrypt(connStr, key)
	if err != nil {
		utils.Write(writer, []byte("global.Serv Error\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}

	// start a transaction to insert the key and the connstr
	txServer, err := global.Serv.Sqlitedb.Begin()
	if err != nil {
		utils.Write(writer, []byte("global.Serv Error\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}
	defer txServer.Rollback()

	txDuck, err := duck.Begin()
	if err != nil {
		utils.Write(writer, []byte("global.Serv Error\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}
	defer txDuck.Rollback()

	txPg, err := postgres.Begin()
	if err != nil {
		utils.Write(writer, []byte("global.Serv Error\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}
	defer txPg.Rollback()
	// insert the key
	_,err = txServer.Stmt(global.Serv.Dbstmt["CreateKey"]).Exec(DBID, string(key))
	if err != nil {
		utils.Write(writer, []byte("global.Serv Error\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}
	// insert the connstr
	_,err = txServer.Stmt(global.Serv.Dbstmt["CreateLink"]).Exec(DBID, encryptedConnStr)
	if err != nil {
		utils.Write(writer, []byte("global.Serv Error\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}
	utils.Write(writer, []byte("successful Linking\n starting the schema migration....\n"))
	
	// migrate schema
	err = internal.Migrate(DBID, connStr, global.Serv.Dbstmt["CreateTable"], txPg, txDuck, txServer)
	if err != nil {
		global.Serv.ErrorLog.Println(err)
		utils.Write(writer, []byte("Error while migrating"))
		return
	}

	err = internal.Audit(txPg)
	if err != nil {
		global.Serv.ErrorLog.Println(err)
		utils.Write(writer, []byte("Error while migrating"))
		return
	}

	err = txPg.Commit()
	if err != nil {
		global.Serv.ErrorLog.Println(err)
		utils.Write(writer, []byte("Error while migrating"))
		return
	}
	err = txDuck.Commit()
	if err != nil {
		global.Serv.ErrorLog.Println(err)
		utils.Write(writer, []byte("Error while migrating"))
		return
	}
	err = txServer.Commit()
	if err != nil {
		global.Serv.ErrorLog.Println(err)
		utils.Write(writer, []byte("Error while migrating"))
		return
	}

	utils.Write(writer, []byte("migration is successful"))
	
}

func MigrateHandler(privilege string, req []string, writer *bufio.Writer) { // todo
	if privilege != "super" {
		utils.Write(writer, []byte("Unauthorized\n"))
		return
	}
	dbname := req[0]
	// check the existince of the database
	var DBID int
	err := global.Serv.Dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err != nil {
		utils.Write(writer, []byte("database: " + dbname + " does not exists\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}
	var connStrEncrypted string
	err = global.Serv.Dbstmt["SelectLink"].QueryRow(DBID).Scan(&connStrEncrypted)
	if err != nil {
		utils.Write(writer, []byte("database: " + dbname + " is not linked to any postgreSQL database\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}
	var key string
	err = global.Serv.Dbstmt["SelectKey"].QueryRow(DBID).Scan(&key)
	if err != nil {
		utils.Write(writer, []byte("database: " + dbname + " key is missing\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}

	connStr, err := utils.Decrypt(connStrEncrypted, []byte(key))
	if err != nil {
		utils.Write(writer, []byte("database: " + dbname + " wrong key or global.Serv error\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}
	// open duckdb
	duck, err := sqlx.Open("duckdb", os.Getenv("DBdir") + "/users/" + dbname + ".db")
	if err != nil {
		utils.Write(writer, []byte("error while connecting to the duckdb database\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}
	defer duck.Close()

	postgres, err := sqlx.Open("postgres", connStr)
	if err != nil {
		utils.Write(writer, []byte("database: " + dbname + " could not reach to postgreSQL\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}
	defer postgres.Close()

	if err := internal.ReadAudit(duck, postgres); err != nil {
		utils.Write(writer, []byte("database: " + dbname + " error while migrating\n"))
		global.Serv.ErrorLog.Println(err)
		return
	}

	utils.Write(writer, []byte("migration is successful"))
}


// create database [dbname],
// create user [dbname] [username] [password]
func CreatHandler(privilege string, req []string, writer *bufio.Writer) {
	if privilege != "super" {
		utils.Write(writer, []byte("Unauthorized\n"))
		return
	}

	if (req[0] != "database" && req[0] != "user") || (req[0] == "database" && len(req) != 2) {
		utils.Write(writer, []byte("ERROR: BAD request\n"))
		return
	}

	if req[0] == "database" {
		CreateDB(req[1], writer)
		return
	}

	if len(req) != 4 {
		utils.Write(writer, []byte("ERROR: BAD request\n"))
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


// grant database [dbname] [username] [accesstype] ,
// grant table [dbname] [tablename] [username] [accesstype] 
func GrantHandler(privilege string, req []string, writer *bufio.Writer) { 
	if privilege != "super" {
		utils.Write(writer, []byte("Unauthorized\n"))
		return
	}
	if (req[0] != "database" && req[0] != "table") || (req[0] == "database" && len(req) != 4) || (req[0] == "table" && len(req) != 5) {
		utils.Write(writer, []byte("ERROR: BAD request\n"))
		return
	}

	if req[0] == "database" {
		GrantDB(req[1], req[2], req[3], writer)
		return
	}

	GrantTable(req[1], req[2], req[3], req[4], writer)
}

func GrantDB(dbname, username, accesstype string, writer *bufio.Writer) {
	accesstype = strings.ToLower(accesstype)
	// check for DB access
	if accesstype != "read" && accesstype != "write" {
		utils.Write(writer, []byte("unsupported Access\n"))
		return
	}
	// get DBID , UID
	var DBID, UID int
	err := global.Serv.Dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err != nil {
		utils.Write(writer, []byte("database: " + dbname + " does not exists\n"))
		return
	}

	err = global.Serv.Dbstmt["SelectUser"].QueryRow(username).Scan(&UID)
	if err != nil {
		utils.Write(writer, []byte("user: " + username + " does not exists\n"))
		return
	}

	// grant access

	_, err = global.Serv.Dbstmt["GrantDB"].Exec(DBID, UID, accesstype)

	if err != nil {
        utils.Write(writer, []byte("global.Serv ERROR\n"))
        global.Serv.ErrorLog.Println(err)
        return
    }

	utils.Write(writer, []byte("success\n"))
}

func GrantTable(dbname, tablename, username, accesstype string, writer *bufio.Writer) {
	accesstype = strings.ToLower(accesstype)
	// check for DB access
	if accesstype != "select" && accesstype != "insert" && accesstype != "update" && accesstype != "delete"{
		utils.Write(writer, []byte("unsupported Access\n"))
		return
	}

	var DBID, UID, TID int
	err := global.Serv.Dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err != nil {
		utils.Write(writer, []byte("database: " + dbname + " does not exists\n"))
		return
	}

	err = global.Serv.Dbstmt["SelectUser"].QueryRow(username).Scan(&UID)
	if err != nil {
		utils.Write(writer, []byte("user: " + username + " does not exists\n"))
		return
	}

	err = global.Serv.Dbstmt["SelectTable"].QueryRow(tablename, DBID).Scan(&TID)
	if err != nil {
		utils.Write(writer, []byte("table: " + tablename + " does not exists in database: " + strconv.Itoa(DBID) + "\n"))
		return
	}
	
	_, err = global.Serv.Dbstmt["GrantTable"].Exec(TID, UID, accesstype)
	if err != nil {
        utils.Write(writer, []byte("global.Serv ERROR\n"))
        global.Serv.ErrorLog.Println(err)
        return
    }

	utils.Write(writer, []byte("success\n"))
}