package main

import (
	"database/sql"
)

func LoginHandler(UserName, password string, dbstmt *sql.Stmt) (int , error){
	
	var UID int
	err := dbstmt.QueryRow(UserName , password).Scan(&UID)
	if err != nil {
		return -1 , err
	}

	return UID , nil
}