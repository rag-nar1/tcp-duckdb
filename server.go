package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"github.com/joho/godotenv"
	"database/sql"
	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
	"github.com/mattn/go-sqlite3"
)


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
	sqlitedb *sql.DB

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

func checkFetal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func checkLog(err error) {
	if err != nil {
		log.Println(err)
	}
}


/**
 * Function to create a tcp listener
 * @pram -> address string : address to listen to 
 * @return -> net.Listener , error
*/
func CreateListener(address string) (net.Listener , error) { 
	listner , err := net.Listen("tcp", address)
	if err != nil {
		return nil , err
	}
	return listner , nil
}


func ProccessQuery(queryResult *sql.Rows) (string , error) {
	columns , err := queryResult.Columns()
	columnsTypes , err := queryResult.ColumnTypes()
	if err != nil {
		return "" , err
	}

	var respones string = ""
	for _ , column := range columns {
		respones += column + " , "
	}
	respones = strings.TrimSuffix(respones , " , ")
	respones += "\n"

	for queryResult.Next() {
		rowData := make([]interface{} , len(columns))
		rowDataPtr := make([]interface{} , len(columns))
		for i := range rowData {
			rowDataPtr[i] = &rowData[i]
		}

		err = queryResult.Scan(rowDataPtr...)

		values := make(map[string]interface{})
		for i , column := range columns {
			values[column] = rowData[i]
		}

		for i , column := range columns {
			// interprate the value based on the type
			switch columnsTypes[i].DatabaseTypeName() {
				case "INT":
					respones += fmt.Sprintf("%d" , values[column].(int))
				case "BIGINT":
					respones += fmt.Sprintf("%d" , values[column].(int64))
				case "FLOAT":
					respones += fmt.Sprintf("%f" , values[column].(float64))
				case "VARCHAR":
					respones += fmt.Sprintf("%s" , values[column].(string))
				case "BOOL":
					respones += fmt.Sprintf("%t" , values[column].(bool))
				default:
					respones += fmt.Sprintf("%v" , values[column])
			}
			respones += " , "
		}
		respones = strings.TrimSuffix(respones , " , ")
		respones += "\n"
	}
	return respones , nil
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
			PasswordPtr , err := SELECT[string](sqlitedb , "admin" , "password" , "username LIKE " + addSingleQuete(Username))
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
	sqlitedb , err = sql.Open("sqlite3" , dbpath + serverdb)
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
	defer sqlitedb.Close()
	listener , err := CreateListener(ServerHost + ":" + ServerPort)
	checkFetal(err)
	err = Start(listener)
	checkLog(err)
}