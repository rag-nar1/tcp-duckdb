package main

import (
	connect 	"TCP-Duckdb/connect"
	create 		"TCP-Duckdb/create"
	grant		"TCP-Duckdb/grant"
	link		"TCP-Duckdb/link"
	internal 	"TCP-Duckdb/internal"
	response 	"TCP-Duckdb/response"
	global 		"TCP-Duckdb/server"
	utils 		"TCP-Duckdb/utils"
	login		"TCP-Duckdb/login"

	"bufio"
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
			response.InternalError(writer)
			global.Serv.ErrorLog.Println(err)
			global.Serv.InfoLog.Println("Connection closed")
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
			link.Handler(privilege, req[1:], writer)
		}

		if req[0] == "migrate" {
			MigrateHandler(privilege, req[1:], writer)
		}
	}
	
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
