package main

import (
	connect "github.com/rag-nar1/TCP-Duckdb/connect"
	create "github.com/rag-nar1/TCP-Duckdb/create"
	grant "github.com/rag-nar1/TCP-Duckdb/grant"
	link "github.com/rag-nar1/TCP-Duckdb/link"
	login "github.com/rag-nar1/TCP-Duckdb/login"
	migrate "github.com/rag-nar1/TCP-Duckdb/migrate"
	response "github.com/rag-nar1/TCP-Duckdb/response"
	global "github.com/rag-nar1/TCP-Duckdb/server"
	"github.com/rag-nar1/TCP-Duckdb/update"
	utils "github.com/rag-nar1/TCP-Duckdb/utils"

	"bufio"
	"net"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	global.Serv.InfoLog.Println("Serving " + conn.RemoteAddr().String())
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	var (
		UID       int
		userName  string
		privilege string
	)

	if err := login.Handler(reader, writer, &UID, &userName, &privilege); err != nil {
		return
	}

	Router(UID, userName, privilege, reader, writer)
}

func Router(UID int, UserName, privilege string, reader *bufio.Reader, writer *bufio.Writer) {
	global.Serv.InfoLog.Println("Serving: " + UserName)
	rawreq := make([]byte, 1024)
	for {
		n, err := reader.Read(rawreq)
		if err != nil {
			response.InternalError(writer)
			global.Serv.ErrorLog.Println(err)
			global.Serv.InfoLog.Println("Connection closed")
			return
		}

		req := strings.Split(string(rawreq[0:n]), " ")
		utils.TrimList(req)

		if req[0] != "connect" && req[0] != "create" && req[0] != "update" && req[0] != "grant" && req[0] != "migrate" && req[0] != "link" {
			response.BadRequest(writer)
			continue
		}

		if req[0] == "connect" {
			if len(req) != 2 {
				response.BadRequest(writer)
				continue
			}
			connect.Handler(global.Serv, UID, UserName, privilege, req[1], reader, writer)
			continue
		}

		if req[0] == "create" {
			create.Handler(privilege, req[1:], writer)
			continue
		}

		if req[0] == "grant" {
			grant.Handler(privilege, req[1:], writer)
			continue
		}

		if req[0] == "link" {
			link.Handler(privilege, req[1:], writer)
		}

		if req[0] == "migrate" {
			migrate.Handler(privilege, req[1:], writer)
		}

		if req[0] == "update" {
			update.Handler(privilege, req[1:], writer)
		}
	}

}
