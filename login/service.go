package login

import (
	utils "TCP-Duckdb/utils"

	"database/sql"
)

func Login(UserName, password string, Dbstmt *sql.Stmt) (int , string , error){
	
	var UID int
	var privilige string
	err := Dbstmt.QueryRow(UserName , utils.Hash(password)).Scan(&UID, &privilige)
	if err != nil {
		return -1 , "", err
	}

	return UID, privilige, nil
}