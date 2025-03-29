package grant

import (
	"bufio"

	response "github.com/rag-nar1/tcp-duckdb/response"
)

// database [dbname] [username] [accesstype] ,
// table [dbname] [tablename] [username] [accesstype]
func Handler(privilege string, req []string, writer *bufio.Writer) {
	if privilege != "super" {
		response.UnauthorizedError(writer)
		return
	}
	if (req[0] != "database" && req[0] != "table") || (req[0] == "database" && len(req) != 4) || (req[0] == "table" && len(req) < 5) {
		response.BadRequest(writer)
		return
	}

	if req[0] == "database" {
		GrantDB(writer, req[1], req[2], req[3])
		return
	}

	GrantTable(writer, req[1], req[2], req[3], req[4:]...)
}
