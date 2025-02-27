package main

import (
	"TCP-Duckdb/internal"
	"TCP-Duckdb/utils"
	"bufio"
	"database/sql"
	"crypto/rand"
	// "fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/lib/pq"
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

func DBHandler(UID int, UserName, privilege string, reader *bufio.Reader, writer *bufio.Writer) {
	infoLog.Println("Serving: " + UserName)
	rawreq := make([]byte, 1024)
	for {
		n , err := reader.Read(rawreq)
		if err != nil {
			Write(writer, []byte("ERROR: while reading\n"))
			errorLog.Println(err)
			return
		}


		req := strings.Split(string(rawreq[0:n]) , " ")
		utils.TrimList(req)
		
		if req[0] != "connect" && req[0] != "create" && req[0] != "grant" && req[0] != "migrate" && req[0] != "link" {
			Write(writer, []byte("ERROR: BAD request\n"))
			continue
		}
		
		if req[0] == "connect" {
			if len(req) != 2 {
				Write(writer, []byte("ERROR: BAD request\n"))
				continue
			}
			DbConnectionHandler(UID, UserName, privilege, req[1], reader, writer) 
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

		MigrateHandler()

	}
	
}
// link dbname connStr
func LinkHandler(privilege string, req []string, writer *bufio.Writer) {
	if privilege != "super" {
		Write(writer, []byte("Unauthorized\n"))
		return
	}

	dbname, connStr := req[0], req[1]
	// check the existince of the database
	var DBID int
	err := server.dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err != nil {
		Write(writer, []byte("database: " + dbname + " does not exists\n"))
		errorLog.Println(err)
		return
	}
	// open duckdb
	duck, err := sql.Open("duckdb", os.Getenv("DBdir") + "/users/" + dbname + ".db")
	if err != nil {
		Write(writer, []byte("error while connecting to the duckdb database\n"))
		errorLog.Println(err)
		return
	}

	// check the connStr
	postgres, err := sql.Open("postgres", connStr)
	if err != nil {
		Write(writer, []byte("error while connecting to the postgresql database\n"))
		errorLog.Println(err)
		return
	}
	defer postgres.Close()

	err = postgres.Ping()
	if err != nil {
		Write(writer, []byte("error while connecting to the postgresql database\n"))
		errorLog.Println(err)
		return
	}

	// generete random 32 byte key for encryption
	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		Write(writer, []byte("Server Error\n"))
		errorLog.Println(err)
		return
	}

	encryptedConnStr , err := utils.Encrypt(connStr, key)
	if err != nil {
		Write(writer, []byte("Server Error\n"))
		errorLog.Println(err)
		return
	}

	// start a transaction to insert the key and the connstr
	transaction, err := server.Sqlitedb.Begin()
	if err != nil {
		Write(writer, []byte("Server Error\n"))
		errorLog.Println(err)
		return
	}
	defer transaction.Rollback()

	// insert the key
	_,err = transaction.Stmt(server.dbstmt["CreateKey"]).Exec(DBID, string(key))
	if err != nil {
		Write(writer, []byte("Server Error\n"))
		errorLog.Println(err)
		return
	}
	// insert the connstr
	_,err = transaction.Stmt(server.dbstmt["CreateLink"]).Exec(DBID, encryptedConnStr)
	if err != nil {
		Write(writer, []byte("Server Error\n"))
		errorLog.Println(err)
		return
	}
	Write(writer, []byte("successful Linking\n starting the schema migration....\n"))
	
	// migrate schema
	err = internal.Migrate(DBID, connStr, transaction, server.dbstmt["CreateTable"], postgres, duck)
	if err != nil {
		errorLog.Println(err)
		Write(writer, []byte("Error while migrating"))
		return
	}
	Write(writer, []byte("migration is successful"))
	transaction.Commit()
}




func MigrateHandler() { // todo
    
}
// create database [dbname],
// create user [dbname] [username] [password]
func CreatHandler(privilege string, req []string, writer *bufio.Writer) {
	if privilege != "super" {
		Write(writer, []byte("Unauthorized\n"))
		return
	}

	if (req[0] != "database" && req[0] != "user") || (req[0] == "database" && len(req) != 2) {
		Write(writer, []byte("ERROR: BAD request\n"))
		return
	}

	if req[1] == "database" {
		CreateDB(req[2], writer)
		return
	}

	if len(req) != 4 {
		Write(writer, []byte("ERROR: BAD request\n"))
		return
	}

	CreateUser(req[2], req[3], req[4], writer) 
}

func CreateDB(dbname string, writer *bufio.Writer) {
	var DBID int
	err := server.dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err == nil && DBID != 0 {
		Write(writer, []byte("database: " + dbname + " already exists\n"))
		return
	}
	
	// create file
	_ ,err = sql.Open("duckdb", os.Getenv("DBdir") + "/users/" + dbname + ".db")
	if err != nil {
		errorLog.Fatal(err)
	}

	_ , err = server.dbstmt["CreateDB"].Exec(dbname)
	if err != nil {
		Write(writer, []byte("ERROR: could not create databse\n"))

		err = os.Remove(os.Getenv("DBdir") + "/users/" + dbname + ".db")
		if err != nil {
			errorLog.Fatal(err)
		}
		return
	}
	Write(writer, []byte("success\n"))
}

func CreateUser(dbname, NewUser, password string, writer *bufio.Writer) {
	var DBID, UID int
	err := server.dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err != nil {
		Write(writer, []byte("database: " + dbname + " does not exists\n"))
		return
	}

	//start transaction
	transaction, err := server.Sqlitedb.Begin()
	defer transaction.Rollback()

	// create user
	_, err = transaction.Stmt(server.dbstmt["CreateUser"]).Exec(NewUser, utils.Hash(password) , "norm")

	if err != nil {
		Write(writer, []byte("user already exists\n"))
		errorLog.Println(err)
		return
	}

	var privilege string
	err = transaction.Stmt(server.dbstmt["login"]).QueryRow(NewUser , utils.Hash(password)).Scan(&UID , &privilege)
	if err != nil {
		Write(writer, []byte("server error\n"))
		errorLog.Println(err)
		return
	}

	_, err = transaction.Stmt(server.dbstmt["GrantDB"]).Exec(DBID, UID, "read")
	if err != nil {
		Write(writer, []byte("server error\n"))
		errorLog.Println(err)
		return
	}
	
	err = transaction.Commit()
	if err != nil {
		Write(writer, []byte("server error\n"))
		errorLog.Println(err)
		return
	}
	Write(writer, []byte("success\n"))
}

func DbConnectionHandler(UID int, UserName, privilege, dbname string, reader *bufio.Reader, writer *bufio.Writer) {
	// check for db existense
	var DBID int
	err := server.dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err != nil {
		Write(writer, []byte("database: " + dbname + " does not exists\n"))
		return
	}

	// check for authrization
	var access int 
	err = server.dbstmt["CheckDbAccess"].QueryRow(UID , DBID).Scan(&access)
	if err != nil {
		Write(writer, []byte("server error\n"))
		errorLog.Println(err)
		return
	}

	if access == 0 && privilege != "super" {
		Write(writer, []byte("user: " + UserName + " does not have access over database: " + dbname + "\n"))
		return
	}

	buffer := make([]byte, 4096)

    _ , err = sql.Open("duckdb" , os.Getenv("DBdir") + "/users/" + dbname + ".db")
    if err != nil {
        Write(writer, []byte("SERVER ERROR\n"))
        errorLog.Println(err)
        return
    }
    Write(writer, []byte("success\n"))
    
	for {
		n , err := reader.Read(buffer)
		if err != nil {
			Write(writer, []byte("ERROR: while reading\n"))
			errorLog.Println(err)
			return
		}

        query := strings.ToLower(strings.Split(string(buffer[0:n]) , " ")[0])

        if query == "start" {
            if strings.ToLower(utils.Trim(string(buffer[0:n]))) != "start transaction" {
                Write(writer, []byte("Bad Request\n"))
                continue
            }
            TransactionHandler(UID, DBID, dbname, privilege, reader, writer)
            continue
        }

        // single query
        QueryHandler(utils.Trim(string(buffer[0:n])), UserName, dbname, privilege, UID, DBID, writer)
        
	}

}

func TransactionHandler(UID, DBID int, dbname, privilege string, reader *bufio.Reader, writer *bufio.Writer) {
    buffer := make([]byte, 4096)
    db , err := sql.Open("duckdb", os.Getenv("DBdir") + "/users/" + dbname + ".db")
    if err != nil {
		Write(writer, []byte("SERVER ERROR\n"))
		errorLog.Println(err)
		return
	}

    transaction, err := db.Begin()
    if err != nil {
		Write(writer, []byte("SERVER ERROR\n"))
		errorLog.Println(err)
		return
	}
    defer transaction.Rollback()

    for {
        n , err := reader.Read(buffer)
		if err != nil {
			Write(writer, []byte("ERROR: while reading\n"))
			log.Println("ERROR" , err)
			return
		}
        query := strings.ToLower(utils.Trim(string(buffer[0:n])))
        if strings.HasPrefix(query, "rollback") {
            return
        }
        if strings.HasPrefix(query, "commit") {
            err = transaction.Commit()
            if err != nil {
                errorLog.Println(err)
                Write(writer, []byte("Error while commiting transaction\n"))
            }
            return
        }
        
        success, err := QueryExecuterTx(query, privilege, UID, DBID, transaction, writer)
        if !success{
            if err != nil {
                errorLog.Println(err)
                Write(writer, []byte("Error while executing query\n"))
                return
            }
            Write(writer, []byte("Error while executing query: ACCESS denied\n"))
            return
        }
    }

}

func QueryExecuterTx(query, privilege string, UID, DBID int, tx *sql.Tx, writer *bufio.Writer) (bool, error) {
    authraized, err := AccessHandler(query, privilege, UID, DBID)
    if err != nil {
        return false, err
	}
    if !authraized {
        return false, nil
    }

    if strings.HasPrefix(query, "select") {
        data, err := internal.SELECT(tx, query)
        if err != nil {
            return false, err
        }
		data = append(data, '\n')
        Write(writer, data)
        return true, nil
	}

	if strings.HasPrefix(query, "create") { 
		err = internal.CREATE(tx, server.Sqlitedb, server.dbstmt["CreateTable"], query, DBID)
        if err != nil {
            return false, err
        }
        return true, nil
	}

    // other statements
	LastInsertedID, RowsAffected, err := internal.EXEC(tx, query)
	if err != nil {
        return false, err
    }

	data := []byte(strconv.Itoa(int(LastInsertedID)) + " " + strconv.Itoa(int(RowsAffected)) + "\n")
	Write(writer, data)
    return true, nil
}

func AccessHandler(query, privilege string, UID, DBID int) (bool, error){
    query = strings.ToLower(query)
	if privilege != "super" {
		hasaccess , err := internal.CheckAccesOverTable(server.Sqlitedb, server.dbstmt["CheckTableAccess"], query, UID, DBID)
		hasDDL , err := internal.CheckDDLActions(query)
		if err != nil {
			return false, err
		}
        return (hasaccess && !hasDDL), nil
	}
    return true, nil
}

func QueryHandler(query, username, dbname, privilege string, UID, DBID int, writer *bufio.Writer) {
	query = strings.ToLower(query)
    authraized, err := AccessHandler(query, privilege, UID, DBID)
	if err != nil {
        Write(writer, []byte("SERVER ERROR\n"))
        errorLog.Println(err)
        return
	}
    if !authraized {
        Write(writer, []byte("Access denied\n"))
        return
    }
		
	db , err := sql.Open("duckdb", os.Getenv("DBdir") + "/users/" + dbname + ".db")
	if err != nil {
		Write(writer, []byte("SERVER ERROR\n"))
		errorLog.Println(err)
		return
	}

	if strings.HasPrefix(query, "select") {
        data, err := internal.SELECT(db, query)
        if err != nil {
            Write(writer, []byte("SERVER ERROR\n"))
            errorLog.Println(err)
            return
        }
		data = append(data, '\n')
        Write(writer, data)
        return
	}

	if strings.HasPrefix(query, "create") { 
		err = internal.CREATE(db, server.Sqlitedb, server.dbstmt["CreateTable"], query, DBID)
		if err != nil {
			Write(writer, []byte("SERVER ERROR\n"))
			errorLog.Println(err)
			return
		}
		Write(writer, []byte("success\n"))
        return
	}

    // other statements
	LastInsertedID, RowsAffected, err := internal.EXEC(db, query)
	if err != nil {
		Write(writer, []byte("SERVER ERROR\n"))
		errorLog.Println(err)
		return
	}

	data := []byte(strconv.Itoa(int(LastInsertedID)) + " " + strconv.Itoa(int(RowsAffected)) + "\n")
	Write(writer, data)
}

// grant database [dbname] [username] [accesstype] ,
// grant table [dbname] [tablename] [username] [accesstype] 
func GrantHandler(privilege string, req []string, writer *bufio.Writer) { 
	if privilege != "super" {
		Write(writer, []byte("Unauthorized\n"))
		return
	}
	if (req[0] != "database" && req[0] != "table") || (req[0] == "database" && len(req) != 4) || (req[0] == "table" && len(req) != 5) {
		Write(writer, []byte("ERROR: BAD request\n"))
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
		Write(writer, []byte("unsupported Access\n"))
		return
	}
	// get DBID , UID
	var DBID, UID int
	err := server.dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err != nil {
		Write(writer, []byte("database: " + dbname + " does not exists\n"))
		return
	}

	err = server.dbstmt["SelectUser"].QueryRow(username).Scan(&UID)
	if err != nil {
		Write(writer, []byte("user: " + username + " does not exists\n"))
		return
	}

	// grant access

	_, err = server.dbstmt["GrantDB"].Exec(DBID, UID, accesstype)

	if err != nil {
        Write(writer, []byte("SERVER ERROR\n"))
        errorLog.Println(err)
        return
    }

	Write(writer, []byte("success\n"))
}

func GrantTable(dbname, tablename, username, accesstype string, writer *bufio.Writer) {
	accesstype = strings.ToLower(accesstype)
	// check for DB access
	if accesstype != "select" && accesstype != "insert" && accesstype != "update" && accesstype != "delete"{
		Write(writer, []byte("unsupported Access\n"))
		return
	}

	var DBID, UID, TID int
	err := server.dbstmt["SelectDB"].QueryRow(dbname).Scan(&DBID)
	if err != nil {
		Write(writer, []byte("database: " + dbname + " does not exists\n"))
		return
	}

	err = server.dbstmt["SelectUser"].QueryRow(username).Scan(&UID)
	if err != nil {
		Write(writer, []byte("user: " + username + " does not exists\n"))
		return
	}

	err = server.dbstmt["SelectTable"].QueryRow(tablename, DBID).Scan(&TID)
	if err != nil {
		Write(writer, []byte("table: " + tablename + " does not exists in database: " + strconv.Itoa(DBID) + "\n"))
		return
	}
	
	_, err = server.dbstmt["GrantTable"].Exec(TID, UID, accesstype)
	if err != nil {
        Write(writer, []byte("SERVER ERROR\n"))
        errorLog.Println(err)
        return
    }

	Write(writer, []byte("success\n"))
}