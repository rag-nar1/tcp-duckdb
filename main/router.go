package main

import (
	connect 	"TCP-Duckdb/connect"
	create 		"TCP-Duckdb/create"
	grant		"TCP-Duckdb/grant"
	internal 	"TCP-Duckdb/internal"
	response 	"TCP-Duckdb/response"
	global 		"TCP-Duckdb/server"
	utils 		"TCP-Duckdb/utils"
	login		"TCP-Duckdb/login"

	"bufio"
	"crypto/rand"
	"database/sql"
	"net"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	global.Serv.InfoLog.Println("Serving " + conn.RemoteAddr().String())
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	var (
		UID int
		userName string
		privilege string
	)

	if err := login.Handler(reader, writer, &UID, &userName, &privilege); err != nil {
		return
	}
	
	Router(UID, userName, privilege, reader, writer)
}

func Router(UID int, UserName, privilege string, reader *bufio.Reader, writer *bufio.Writer) {
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
			response.BadRequest(writer)
			continue
		}
		
		if req[0] == "connect" {
			if len(req) != 2 {
				response.BadRequest(writer)
				continue
			}
			connect.Handler(global.Serv, UID, UserName, privilege, req[1], reader, writer) 
			continue
		}
		
		if req[0] == "create" {
			create.Handler(privilege, req[1:], writer)
			continue
		}
		
		if req[0] == "grant" {
			grant.Handler(privilege, req[1:], writer)
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
