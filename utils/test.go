package utils

import (
	// "bufio"
	"TCP-Duckdb/response"
	"fmt"
	"net"
	"strings"
	"time"
	"os"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func StartUp() (*net.TCPConn){
	if err := godotenv.Load("../.env"); err != nil {
		panic(err)
	}
	conn := Connection()
	return conn
}

func Connection() *net.TCPConn {
	tcpAddr, err := net.ResolveTCPAddr("tcp","localhost:4000")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp",nil,tcpAddr)
	if err != nil {
		panic(err)
	}
	// Set a deadline for the operation (optional, for timeout)
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	return conn
}

func LoginAsAdmin(conn *net.TCPConn) error {
	_, err := conn.Write([]byte("login duck duck"))
	if err != nil {
		return err
	}
	res := Read(conn)
	if res != response.SuccessMsg {
		return fmt.Errorf("unauth: %s", res)
	}
	return nil
}

func Login(conn *net.TCPConn, username, password string) error {
	_, err := conn.Write([]byte(fmt.Sprintf("login %s %s", username, password)))
	if err != nil {
		return err
	}
	res := Read(conn)
	if res != response.SuccessMsg {
		return fmt.Errorf("unauth: %s", res)
	}
	return nil
}

func Read(conn *net.TCPConn) string {
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}
	return strings.Trim(string(buffer[:n])," \n\t")
} 

func CreateDB(conn *net.TCPConn, dbname string) error {
	if _, err := conn.Write([]byte("create database " + dbname)); err != nil {
		return err
	}

	if res := Read(conn); res != response.SuccessMsg {
		return fmt.Errorf("%s", res)
	}
	return nil
}

func CreateUser(conn *net.TCPConn, userName, password string) error {
	if _, err := conn.Write([]byte("create user " + userName + " " + password)); err != nil {
		return err
	}

	if res := Read(conn); res != response.SuccessMsg {
		return fmt.Errorf("%s", res)
	}
	return nil
}

func CreateTable(conn *net.TCPConn, tablename string) error {
	if _, err := conn.Write([]byte(fmt.Sprintf("CREATE TABLE %s(id int, name text);", tablename))); err != nil {
		return err
	}

	if res := Read(conn); res != response.SuccessMsg {
		return fmt.Errorf("%s", res)
	}
	return nil
}


func ConnectDb(conn *net.TCPConn, dbname string) error {
	_, err := conn.Write([]byte("connect " + dbname))
	if err != nil {
		return err
	}
	res := Read(conn)
	if res != response.SuccessMsg {
		return fmt.Errorf("%s", res)
	}
	return nil
}

func Grant(conn *net.TCPConn, username, dbname, privilege string) error {
	_, err := conn.Write([]byte(fmt.Sprintf("grant database %s %s %s", dbname, username, privilege)))
	if err != nil {
		return err
	}
	res := Read(conn)
	if res != response.SuccessMsg {
		return fmt.Errorf("%s", res)
	}
	return nil
}

func CleanUpDb(db *sql.DB, dbname string) error {
	if err := os.Remove(os.Getenv("DBdir") + "/users/" + dbname + ".db"); err != nil {
		return err
	}

	if _, err := db.Exec("DELETE FROM DB;"); err != nil {
		return err
	}

	if _, err := db.Exec("DELETE FROM dbprivilege;"); err != nil {
		return err
	}

	return nil
}




func CleanUpUsers(db *sql.DB) error {
	
	if _, err := db.Exec("DELETE FROM user WHERE usertype not like 'super';"); err != nil {
		return err
	}
	
	return nil
}

func CleanUpTables(db *sql.DB) error {

	if _, err := db.Exec("DELETE FROM tables;"); err != nil {
		return err
	}

	if _, err := db.Exec("DELETE FROM tableprivilege;"); err != nil {
		return err
	}

	return nil
}

func CleanUp(dbname string) error {
	db , err := sql.Open("sqlite3",os.Getenv("DBdir") + os.Getenv("ServerDbFile"))
	if err != nil {
		return err
	}

	if dbname == "" {
		goto SkipDb
	}

	if err := CleanUpDb(db, dbname); err != nil {
		return err
	}

SkipDb:
	if err := CleanUpUsers(db); err != nil {
		return err
	}

	if err := CleanUpTables(db); err != nil {
		return err
	}
	return nil
}