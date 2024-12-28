package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"github.com/joho/godotenv"
	"database/sql"
	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
	"github.com/mattn/go-sqlite3"
	"sync"
	"time"
)

type SQLitedb struct {
	db *sql.DB
	latch sync.Mutex
}

var (
	currentdir string
	dbpath string

	// postgresql
	Postgresdb *sql.DB
	connectionString string = "postgres://%s:%s@%s:%s/%s?sslmode=disable" // user , password , hostname , port , database
	PostgresHost string
	PostgresPort string
	PostgresUser string
	PostgresPass string
	PostgresDbName string

	ServerHost = "localhost"
	ServerPort = "2003"

	// sqlite3 
	sqlitedb *SQLitedb

	NextID map[string]int
	ServerDbTables = []string {
		"admins",
		"users",
		"database",
		"useraccess",
		"privilege",
		"privilege_type",
		"tables",
	}
)

const (

	serverdb = "main.sqlite3"

)


func NewSqliteDB (dbpath string) (*SQLitedb , error) {
	db , err := sql.Open("sqlite3" , dbpath + "?_journal=WAL&_busy_timeout=5000")
	if err != nil {
        return nil, fmt.Errorf("failed to open database: %v", err)
    }

	// Configure connection pool
	db.SetMaxOpenConns(1)             // Limit to 1 connection to prevent "database is locked" errors
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)
	return &SQLitedb{
        db: db,
    }, nil
}


/**
 * function handles connections while it's open and respones to requets
 * requests format
	* AUTH < Username > < Password > < Database >
	* QUERY < SqlQuery >
	* ADMIN < Password > 
*/


func HandleConnection(connection net.Conn) {
	defer connection.Close()

	reader := bufio.NewReader(connection)

	for {
		message, err := reader.ReadString('\n')
		if err != nil { // send error
			connection.Write([]byte("Error while reading from the connection\nConnection Closed"))
			return 
		}
		
		message = strings.TrimSpace(message)
		if !strings.HasPrefix(message , "AUTH")  && !strings.HasPrefix(message , "ADMIN") && !strings.HasPrefix(message , "EXIT") {
			connection.Write([]byte("Error invalid request\n"))
			continue
		}

		if strings.HasPrefix(message , "EXIT") {
			connection.Write([]byte("success\n"))
			return
		}

		if strings.HasPrefix(message , "ADMIN") {
			// check if the password is correct for the admin
			request := strings.Split(message , " ")
			if len(request) != 3 {
				connection.Write([]byte("Error invalid request\n"))
				continue
			}

			Username , Password := request[1] , request[2]
			PasswordPtr , err := SELECT[string](sqlitedb , "admins" , "password" , "username LIKE " + addSingleQuete(Username))
			if err != nil {
				connection.Write([]byte("Error while executing query admin\n" + err.Error() + "\n"))
				continue
			}

			//check the password
			if *PasswordPtr != Password {
				connection.Write([]byte("Error wrong password\n"))
				continue
			}
			
			connection.Write([]byte("success\n"))
			HandleAdmin(connection)
			return
		}
		// TODO: ADD user interface

	}
}

/**
 * Function to start Accepting connections 
 * @pram -> listner net.Listener
 * @return -> error
*/

func Start(listner net.Listener) error {
	defer listner.Close()
	for {
		connection , err := listner.Accept()
		if err != nil {
			return err
		}
		go HandleConnection(connection)
	}
}




func init() {
	// load the env variables
	err := godotenv.Load()
	checkFetal(err)
	// connect to the postgresql
	PostgresHost = os.Getenv("PostgresHost")
	PostgresPort = os.Getenv("PostgresPort")
	PostgresUser = os.Getenv("PostgresUser")
	PostgresPass = os.Getenv("PostgresPass")
	PostgresDbName = os.Getenv("PostgresDbName")
	checkFetal(err)
	connection := fmt.Sprintf(connectionString , PostgresUser , PostgresPass , PostgresHost , PostgresPort , PostgresDbName)
	Postgresdb , err = sql.Open("postgres" , connection)
	checkFetal(err)
	fmt.Println("connected to postgresql")

	// create the db directory
	currentdir , err = os.Getwd()
	checkFetal(err)
	dbpath = currentdir + "/DB/"

	// register the sqlite3 with extensions
	sql.Register("sqlite3_with_extensions",
	&sqlite3.SQLiteDriver{
		Extensions: []string{
			"sqlite3_mod_regexp",
		},
	})

	// create the main database
	sqlitedb , err = NewSqliteDB(dbpath + serverdb)
	checkFetal(err)
	fmt.Println("connected to sqlite3")


	// initilaize the next id for each table
	NextID = make(map[string]int , 0)
	for _ , table := range ServerDbTables { 
		
		cntPtr , err := SELECT[int](sqlitedb , table , "COUNT(*)" , "TRUE")
		checkFetal(err)
		NextID[table] = *cntPtr + 1
	}
	

}

func main() {
	defer Postgresdb.Close()
	listener , err := CreateListener(ServerHost + ":" + ServerPort)
	checkFetal(err)
	err = Start(listener)
	checkLog(err)
}