package login

import (
	"errors"

	response "github.com/rag-nar1/tcp-duckdb/response"
	global "github.com/rag-nar1/tcp-duckdb/server"
	utils "github.com/rag-nar1/tcp-duckdb/utils"

	"bufio"
	"strings"
)

func Handler(reader *bufio.Reader, writer *bufio.Writer, UID *int, userName, privilege *string) error {
	// read login request
	route := make([]byte, 1024)
	n, err := reader.Read(route)
	if err != nil {
		response.InternalError(writer)
		global.Serv.ErrorLog.Println(err)
		return err
	}
	// check for a valid request
	request := strings.Split(string(route[0:n]), " ")
	if request[0] != "login" || len(request) != 3 {
		response.BadRequest(writer)
		return errors.New(response.BadRequestMsg)
	}
	// validate the userName and password
	var password string
	*userName, password = utils.Trim(request[1]), utils.Trim(request[2])
	*UID, *privilege, err = Login(*userName, password, global.Serv.Dbstmt["login"])

	if err != nil {
		response.UnauthorizedError(writer)
		global.Serv.ErrorLog.Println(err)
		return err
	}
	response.Success(writer)
	return nil
}
