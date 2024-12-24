package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"strconv"
	"log"
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
)


var (
	currentdir string
	dbpath string
	db *sql.DB
	connectionString string = "postgres://%s:%s@%s:%s/%s?sslmode=disable" // user , password , hostname , port , database
	PostgresHost string
	PostgresPort string
	PostgresUser string
	PostgresPass string
	PostgresDbName string
	ServerHost = "localhost"
	ServerPort = "2003"
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

/**
 * function handles connections while it's open and respones to requets
 * requests format
	* AUTH < Username > < Password > < Database >
	* QUERY < SqlQuery >
*/


func HandleConnection(connection net.Conn) {
	defer connection.Close()

	
	
	reader := bufio.NewReader(connection)
	authenticated := false
	var Username , Password , Database string
	var UID , DBID int
outer:
	for {
		message, err := reader.ReadString('\n')
		if err != nil { // send error
			connection.Write([]byte("Error while reading from the connection\nConnection Closed"))
			return 
		}
		
		message = strings.TrimSpace(message)
		if !strings.HasPrefix(message , "AUTH") && strings.HasPrefix(message , "QUERY") {
			connection.Write([]byte("Error invalid request\n"))
			continue
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
				connection.Write([]byte("Error while executing query user\n" + err.Error()))
				continue
			}

			var correctPass *string = new(string)
			var uid *int = new(int)
			userAuth.Next()
			err = userAuth.Scan(uid , correctPass)
			if err != nil {
				connection.Write([]byte("Error invalid Username\n" + err.Error()))
				continue
			}

			//check the password
			if *correctPass != Password {
				connection.Write([]byte("Error wrong password\n"))
				continue
			}

			UID = *uid

			dbAuth , err := db.Query(fmt.Sprintf("SELECT id FROM databases WHERE name LIKE '%s' AND userid = %d" , Database , UID))
			if err != nil {
				connection.Write([]byte("Error while executing query db\n" + err.Error()))
				continue
			}

			var dbid *int = new(int)
			dbAuth.Next()
			err = dbAuth.Scan(dbid)
			if err != nil {
				connection.Write([]byte("Error this user does not own this database \n" + err.Error()))
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

		// get the databse name 
		dbnameResult , err := db.Query(fmt.Sprintf("SELECT name FROM databases WHERE id = %d" , DBID))
		if err != nil {
			connection.Write([]byte("Error while executing query db\n" + err.Error()))
			continue
		}

		var dbname *string = new(string)
		dbnameResult.Next()
		err = dbnameResult.Scan(dbname)
		if err != nil {
			connection.Write([]byte("Error while scanning db name\n" + err.Error()))
			continue
		}

		userDbPath := dbpath + strconv.Itoa(DBID) + "_" + *dbname

		userDb , err := sql.Open("duckdb" , userDbPath)
		if err != nil {
			connection.Write([]byte("Error while connecting to database\n" + err.Error()))
			continue
		}

		//check if the query is select or not
		if strings.HasPrefix(strings.ToUpper(query) , "SELECT") {
			queryResult , err := userDb.Query(query)
			if err != nil {
				connection.Write([]byte("Error while executing query\n" + err.Error()))
				continue
			}

			var respones string
			for queryResult.Next() {
				var row string
				err = queryResult.Scan(&row)
				if err != nil {
					connection.Write([]byte("Error while scanning row\n" + err.Error()))
					continue outer
				}
				respones += row + "\n"
			}
			respones = "success\n" + respones
			connection.Write([]byte(respones))
			continue
		}

		//if the query is not select we use exec 
		queryResult , err := userDb.Exec(query)
		if err != nil {
			connection.Write([]byte("Error while executing query\n" + err.Error()))
			continue
		}
		EffectedRows , err := queryResult.RowsAffected()
		if err != nil {
			connection.Write([]byte("Error while getting effected rows\n" + err.Error()))
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
	err := godotenv.Load()
	PostgresHost = os.Getenv("PostgresHost")
	PostgresPort = os.Getenv("PostgresPort")
	PostgresUser = os.Getenv("PostgresUser")
	PostgresPass = os.Getenv("PostgresPass")
	PostgresDbName = os.Getenv("PostgresDbName")
	checkFetal(err)
	connection := fmt.Sprintf(connectionString , PostgresUser , PostgresPass , PostgresHost , PostgresPort , PostgresDbName)
 	db , err = sql.Open("postgres" , connection)
	checkFetal(err)
	fmt.Println("connected to postgresql")

	currentdir , err = os.Getwd()
	checkFetal(err)
	dbpath = currentdir + "/db/"
}

func main() {
	defer db.Close()
	listener , err := CreateListener(ServerHost + ":" + ServerPort)
	checkFetal(err)
	err = Start(listener)
	checkLog(err)
}