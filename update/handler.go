package update

import (
	"bufio"

	"github.com/rag-nar1/TCP-Duckdb/response"
)

// update databse [dbname] newdbname
// update user username [username] [newusername]
// update user password [username] [password]
func Handler(privilege string, req []string, writer *bufio.Writer) {
	if privilege != "super" {
		response.UnauthorizedError(writer)
		return
	}

	if req[0] != "database" && req[0] != "user" {
		response.BadRequest(writer)
		return
	}

	if req[0] == "database" && len(req) != 3{
		response.BadRequest(writer)
		return
	}

	if req[0] == "user" && len(req) != 4 {
		response.BadRequest(writer)
		return
	}

	if req[0] == "database" {
		UpdateDatabase(writer, req[1], req[2])
		return
	}

	if req[1] != "username" && req[1] != "password" {
		response.BadRequest(writer)
		return
	}

	if req[1] == "username" {
		UpdateUserUsername(writer, req[2], req[3])
		return
	}

	UpdateUserPassword(writer, req[2], req[3])
}