package main

import (
	"bufio"
	"fmt"
	"strconv"
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

	AdminNextID int
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



func HandleAdmin(connection net.Conn) {
	defer connection.Close()

	reader := bufio.NewReader(connection)
	for {
		command , err := reader.ReadString('\n')
		if err != nil {
			connection.Write([]byte("Error while reading from the connection\nConnection Closed"))
			return 
		}

		command = strings.TrimSpace(command)

		if !strings.HasPrefix(command , "CREATE") && !strings.HasPrefix(command , "DELETE") && !strings.HasPrefix(command , "EXIT") {
			connection.Write([]byte("Error invalid request\n"))
			continue
		}
		
		if strings.HasPrefix(command , "EXIT") { // exit the admin mode and close the connection
			connection.Write([]byte("success\n"))
			return
		}

		if strings.HasPrefix(command , "CREATE") { // create a new database or admin user

			command = strings.TrimPrefix(command , "CREATE ")
			if !strings.HasPrefix(command , "DATABASE") && !strings.HasPrefix(command , "ADMIN") && !strings.HasPrefix(command , "USER"){
				connection.Write([]byte("Error invalid request\n"))
				continue
			}

			if strings.HasPrefix(command , "ADMIN") { // create a new server admin
				args := strings.Split(command , " ")
				if len(args) != 3 {
					connection.Write([]byte("Error invalid request\n"))
					continue
				}

				NewAdmin , Password := args[1] , args[2]

				_ , err = sqlitedb.Exec("INSERT INTO admins VALUES(%d , '%s' , '%s')" , AdminNextID , NewAdmin , Password)
				if err != nil {
					connection.Write([]byte("Error while inserting \n" + err.Error() + "\n"))
					continue
				}
				AdminNextID ++
				connection.Write([]byte("success\n"))
				continue
			}

			





			continue
		}

		request := strings.Split(command , " ")
		if len(request) != 2 {
			connection.Write([]byte("Error invalid request\n"))
			continue
		}

		dbname := request[1]

		// delete the database file in the db directory
		err = os.Remove(dbpath + dbname + ".db")
		
		if err != nil {
			connection.Write([]byte("Error while deleting database\n" + err.Error() + "\n"))
			continue
		}
		connection.Write([]byte("success\n"))
		continue
	}

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
	authenticated := false
	var Username , Password , Database string
	var UID , DBID int

	for {
		message, err := reader.ReadString('\n')
		if err != nil { // send error
			connection.Write([]byte("Error while reading from the connection\nConnection Closed"))
			return 
		}
		
		message = strings.TrimSpace(message)
		if !strings.HasPrefix(message , "AUTH") && !strings.HasPrefix(message , "QUERY") && !strings.HasPrefix(message , "ADMIN") && !strings.HasPrefix(message , "EXIT") {
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

			Username , Password = request[1] , request[2]
			adminAuth , err := sqlitedb.Query("SELECT password FROM admin WHERE username LIKE '%s'" , Username)
			if err != nil {
				connection.Write([]byte("Error while executing query admin\n" + err.Error() + "\n"))
				continue
			}

			var correctPass *string = new(string)
			adminAuth.Next()
			err = adminAuth.Scan(correctPass)
			if err != nil {
				connection.Write([]byte("Error invalid Username\n" + err.Error() + "\n"))
				continue
			}

			//check the password
			if *correctPass != Password {
				connection.Write([]byte("Error wrong password\n"))
				continue
			}
			
			connection.Write([]byte("success\n"))
			HandleAdmin(connection)
			return
		}
		
		if strings.HasPrefix(message , "AUTH") {
			request := strings.Split(message , " ")
			if len(request) != 4 {
				connection.Write([]byte("Error invalid request\n"))
				continue
			}

			Username , Password , Database = request[1] , request[2] , request[3]
			userAuth , err := db.Query(fmt.Sprintf("SELECT id , password FROM users WHERE username LIKE '%s'" , Username))
			if err != nil {
				connection.Write([]byte("Error while executing query user\n" + err.Error() + "\n"))
				continue
			}

			var correctPass *string = new(string)
			var uid *int = new(int)
			userAuth.Next()
			err = userAuth.Scan(uid , correctPass)
			if err != nil {
				connection.Write([]byte("Error invalid Username\n" + err.Error() + "\n"))
				continue
			}

			//check the password
			if *correctPass != Password {
				connection.Write([]byte("Error wrong password\n"))
				continue
			}

			UID = *uid

			dbAuth , err := db.Query(fmt.Sprintf("SELECT dbid FROM can_access join database using(dbid) where dbname like '%s' and userid = '%d'" , Database , UID))
			if err != nil {
				connection.Write([]byte("Error while executing query db\n" + err.Error() + "\n"))
				continue
			}

			var dbid *int = new(int)
			dbAuth.Next()
			err = dbAuth.Scan(dbid)
			if err != nil {
				connection.Write([]byte("Error this user does not own this database \n" + err.Error() + "\n"))
				continue
			}

			DBID = *dbid

			authenticated = true

			connection.Write([]byte("success\n"))
			continue
		}

		if !authenticated {
			connection.Write([]byte("Error you can't excute queries with authentication\n"))
			continue
		}

		query := strings.TrimPrefix(message , "QUERY ")


		userDbPath := dbpath + strconv.Itoa(DBID) + "_" + Database + ".db"
		userDb , err := sql.Open("duckdb" , userDbPath)
		if err != nil {
			connection.Write([]byte("Error while connecting to database\n" + err.Error() + "\n"))
			continue
		}
		defer userDb.Close()

		//check if the query is select or not
		if strings.HasPrefix(strings.ToUpper(query) , "SELECT") {
			queryResult , err := userDb.Query(query)
			if err != nil {
				connection.Write([]byte("Error while executing query\n" + err.Error() + "\n"))
				continue
			}
			respones , err := ProccessQuery(queryResult)
			if err != nil {
				connection.Write([]byte("Error while processing query\n" + err.Error() + "\n"))
				continue
			}

			respones = "success\n" + respones
			connection.Write([]byte(respones))
			continue 
		}

		//if the query is not select we use exec 
		queryResult , err := userDb.Exec(query)
		if err != nil {
			connection.Write([]byte("Error while executing query\n" + err.Error() + "\n"))
			continue
		}
		EffectedRows , err := queryResult.RowsAffected()
		if err != nil {
			connection.Write([]byte("Error while getting effected rows\n" + err.Error() + "\n"))
			continue
		}
		respones := fmt.Sprintf("success\nEffected Rows : %d\n" , EffectedRows)
		connection.Write([]byte(respones))

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
	sqlitedb , err = sql.Open("sqlite3_with_extensions" , dbpath + serverdb)
	checkFetal(err)
	fmt.Println("connected to sqlite3")

	admincnt , err := sqlitedb.Query("SELECT count(*) FROM admins;")
	checkFetal(err)

	var ADMINCNT *int
	admincnt.Next()
	
	err = admincnt.Scan(ADMINCNT)
	checkFetal(err)
	AdminNextID = *ADMINCNT + 1

}

func main() {
	defer Postgresdb.Close()
	defer sqlitedb.Close()
	listener , err := CreateListener(ServerHost + ":" + ServerPort)
	checkFetal(err)
	err = Start(listener)
	checkLog(err)
}