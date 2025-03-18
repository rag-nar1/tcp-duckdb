package login

import (
	response	"TCP-Duckdb/response"
	global 		"TCP-Duckdb/server"
	utils 		"TCP-Duckdb/utils"

	"bufio"
	"strings"
)

func Handler(reader *bufio.Reader, writer *bufio.Writer, UID *int, userName, privilege *string) (error){
	// read login request
	route := make([]byte, 1024)
	n , err := reader.Read(route)
	if err != nil {
		global.Serv.ErrorLog.Println(err)
		return err
	}
	// check for a valid request
	request := strings.Split(string(route[0 : n]) , " ")
	if request[0] != "login" || len(request) != 3 {
		response.BadRequest(writer)
		return  err
	}
	// validate the userName and password
	var password  string
	*userName, password = utils.Trim(request[1]), utils.Trim(request[2]) 
	*UID, *privilege, err = Login(*userName, password, global.Serv.Dbstmt["login"])

	if err != nil {
		utils.Write(writer, []byte("Unauthorized\n"))
		global.Serv.ErrorLog.Println(err)
		return  err
	}
	response.Success(writer)
	return  nil
}

