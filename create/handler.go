package create

import (
	response "github.com/rag-nar1/TCP-Duckdb/response"

	"bufio"

	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
)

// database [dbname],
// user [username] [password]
func Handler(privilege string, req []string, writer *bufio.Writer) {
	if privilege != "super" {
		response.UnauthorizedError(writer)
		return
	}

	if (req[0] != "database" && req[0] != "user") || (req[0] == "database" && len(req) != 2) || (req[0] == "user" && len(req) != 3) {
		response.BadRequest(writer)
		return
	}

	if req[0] == "database" {
		CreateDB(req[1], writer)
		return
	}
	CreateUser(req[1], req[2], writer)
}
