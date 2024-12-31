package main

import (
	"bufio"
	"database/sql"
	"net"
	"os"
	"strconv"
	"strings"
)

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
			if !strings.HasPrefix(command , "DATABASE") && !strings.HasPrefix(command , "ADMIN") && !strings.HasPrefix(command , "USER") && !strings.HasPrefix(command , "TABLE"){
				connection.Write([]byte("Error invalid request\n"))
				continue
			}

			args := strings.Split(command , " ")
			if len(args) != 3 {
				connection.Write([]byte("Error invalid request\n"))
				continue
			}

			if strings.HasPrefix(command , "ADMIN") { // create a new server admin
				
				NewAdmin , Password := addSingleQuete(args[1]) , addSingleQuete(args[2])
				err := INSERT(sqlitedb , "admins" , strconv.Itoa(NextID["admins"]) ,NewAdmin , Password)
				if err != nil {
					connection.Write([]byte("Error while inserting \n" + err.Error() + "\n"))
					continue
				}
				NextID["admins"] ++
				connection.Write([]byte("success\n"))
				continue
			}

			if strings.HasPrefix(command , "USER") {

				NewUser , Password := addSingleQuete(args[1]) , addSingleQuete(args[2])

				err := INSERT(sqlitedb , "users" , strconv.Itoa(NextID["users"]) , NewUser , Password)
				if err != nil {
					connection.Write([]byte("Error while inserting \n" + err.Error() + "\n"))
					continue
				}
				NextID["users"] ++
				connection.Write([]byte("success\n"))
				continue
			}

			if strings.HasPrefix(command , "DATABASE") {
				NewDB , SuperUser := addSingleQuete(args[1]) , addSingleQuete(args[2])

				DataPtr , err := SELECT[int](sqlitedb , "users" , "userid" , "username like " + SuperUser)
				if err != nil {
					connection.Write([]byte("Error while inserting \n" + err.Error() + "\n"))
					continue
				}

				UID := *DataPtr
				DBID := NextID["database"]
				err = INSERT(sqlitedb , "database" , strconv.Itoa(DBID) , NewDB , strconv.Itoa(UID))

				if err != nil {
					connection.Write([]byte("Error while inserting \n" + err.Error() + "\n"))
					continue
				}

				err = INSERT(sqlitedb , "useraccess" ,  strconv.Itoa(DBID) , strconv.Itoa(UID))
				if err != nil {
					connection.Write([]byte("Error while inserting \n" + err.Error() + "\n"))
					continue
				}
				
				
				_ , err = sql.Open("duckdb" , dbpath + strconv.Itoa(UID) + "_" + NewDB) // creates the db file
				if err != nil {
					connection.Write([]byte("Error while Creating file \n" + err.Error() + "\n"))
					// TODO: Delete the entry
					continue
				}

				NextID["database"] ++
				connection.Write([]byte("success\n"))
				continue
			}

			tableName , dbName := args[1] , args[2]
			
			DataPtr , err := SELECT[int](sqlitedb , "database" , "dbid" , "dbname like " + dbName)

			if err != nil {
				connection.Write([]byte("Error while inserting \n" + err.Error() + "\n"))
				continue
			}

			DBID := *DataPtr
			TableID := NextID["tables"]
			err = INSERT(sqlitedb , "tables" , strconv.Itoa(TableID) , addSingleQuete(tableName) , strconv.Itoa(DBID))

			if err != nil {
				connection.Write([]byte("Error while inserting \n" + err.Error() + "\n"))
				continue
			}

			NextID["tables"] ++
			connection.Write([]byte("success\n"))
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