package utils

import (
	// "bufio"
	"TCP-Duckdb/response"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
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

	if res := Read(conn); strings.HasPrefix(res, "ERROR") {
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

func Query(conn *net.TCPConn, query string) error {
	_, err := conn.Write([]byte(query))
	if err != nil {
		return err
	}
	res := Read(conn)
	if strings.HasPrefix(res, "Error") {
		return fmt.Errorf("%s", res)
	}
	return nil
}

func QueryData(conn *net.TCPConn, query string) (string, error) {
	_, err := conn.Write([]byte(query))
	if err != nil {
		return "", err
	}
	res := Read(conn)
	if strings.HasPrefix(res, "Error") {
		return "", fmt.Errorf("%s", res)
	}
	return res, nil
}

func GrantDb(conn *net.TCPConn, username, dbname, privilege string) error {
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

func GrantTable(conn *net.TCPConn, username, dbname, tablename, privilege string) error {
	_, err := conn.Write([]byte(fmt.Sprintf("grant table %s %s %s %s", dbname, tablename, username, privilege)))
	if err != nil {

		return err
	}
	res := Read(conn)
	if res != response.SuccessMsg {
		return fmt.Errorf("%s", res)
	}
	return nil
}

func CleanUpDb(db *sql.DB) error {
	files, err := filepath.Glob("../storge/users/*")
	if err != nil {
		return err
	}

	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			return err
		}
	}

	if _, err := db.Exec("DELETE FROM DB;"); err != nil {
		return err
	}

	if _, err := db.Exec("DELETE FROM dbprivilege;"); err != nil {
		return err
	}
	if _, err := db.Exec("DELETE FROM postgres;"); err != nil {
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

func CleanUp() {
	db , err := sql.Open("sqlite3","../storge/server/db.sqlite3")
	if err != nil {
		log.Fatal(err)
		 
	}
	defer db.Close()

	if err := CleanUpDb(db); err != nil {
		log.Fatal(err)
		 
	}

	if err := CleanUpUsers(db); err != nil {
		log.Fatal(err)
		 
	}

	if err := CleanUpTables(db); err != nil {
		log.Fatal(err)
		 
	}
	if err := CleanUpLink(); err != nil {
		log.Fatal(err)
		
	}
}

func Link(conn *net.TCPConn, dbname, connStr string) error {
	if _, err := conn.Write([]byte(fmt.Sprintf("link %s %s", dbname, connStr))); err != nil {
		return err
	}
	res := Read(conn)
	if strings.HasPrefix(res, "Error") {
		return fmt.Errorf("%s", res)
	}
	return nil
}

func CleanUpLink() error {
	pq , err := sql.Open("postgres", "postgresql://postgres:1242003@localhost:5432")
	if err != nil {
		return err
	}
	defer pq.Close()

	if _, err := pq.Exec("Drop database testdb;"); err != nil {
		return err
	}
	if _, err := pq.Exec("create database testdb;"); err != nil {
		return err
	}
	pq, err = sql.Open("postgres", "postgresql://postgres:1242003@localhost:5432/testdb")
	if err != nil {
		return err
	}
	defer pq.Close()

	for _,t := range []string{"t1", "t2", "t3"} {
		if _, err := pq.Exec(fmt.Sprintf("create table %s(id int primary key);", t)); err != nil {
			return err
		}
		for i := 1; i <= 3; i ++ {
			if _, err := pq.Exec(fmt.Sprintf("insert into %s(id) values(%d);", t,i)); err != nil {
				return err
			}
		}
	}
	return nil
}

func Migrate(conn *net.TCPConn, dbname string) error {
	if _, err := conn.Write([]byte(fmt.Sprintf("migrate %s", dbname))); err != nil {
		return err
	}
	res := Read(conn)
	if strings.HasPrefix(res, "Error") {
		return fmt.Errorf("%s", res)
	}
	return nil
}