package main

import (
	"bufio"
	"net"
	"strings"
	
	create "github.com/rag-nar1/tcp-duckdb/commands/create"
	grant "github.com/rag-nar1/tcp-duckdb/commands/grant"
	link "github.com/rag-nar1/tcp-duckdb/commands/link"
	migrate "github.com/rag-nar1/tcp-duckdb/commands/migrate"
	update "github.com/rag-nar1/tcp-duckdb/commands/update"
	"github.com/rag-nar1/tcp-duckdb/connect"
	"github.com/rag-nar1/tcp-duckdb/login"
	"github.com/rag-nar1/tcp-duckdb/response"
	"github.com/rag-nar1/tcp-duckdb/utils"
	"github.com/rag-nar1/tcp-duckdb/server"
	
	_ "github.com/lib/pq"
	_ "github.com/marcboeker/go-duckdb"
)

type Handler interface {
	Handler(privilege string, req []string, writer *bufio.Writer)
}

type HandlerFunc func(privilege string, req []string, writer *bufio.Writer)

func (f HandlerFunc) Handler(privilege string, req []string, writer *bufio.Writer) {
	f(privilege, req, writer)
}

var Handlers = map[string]Handler{
	"create": HandlerFunc(create.Handler),
	"grant": HandlerFunc(grant.Handler),
	"link": HandlerFunc(link.Handler),
	"migrate": HandlerFunc(migrate.Handler),
	"update": HandlerFunc(update.Handler),
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	server.Serv.InfoLog.Println("Serving " + conn.RemoteAddr().String())
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
	server.Serv.InfoLog.Println("Serving: " + UserName)
	rawreq := make([]byte, 1024)
	for {
		n, err := reader.Read(rawreq)
		if err != nil {
			response.InternalError(writer)
			server.Serv.ErrorLog.Println(err)
			server.Serv.InfoLog.Println("Connection closed")
			return
		}

		req := strings.Split(string(rawreq[0:n]), " ")
		utils.TrimList(req)

		if _,ok := Handlers[req[0]]; !ok && req[0] != "connect"{
			response.BadRequest(writer)
			continue
		}

		if req[0] == "connect" {
			if len(req) != 2 {
				response.BadRequest(writer)
				continue
			}
			connect.Handler(UID, UserName, privilege, req[1], reader, writer)
			continue
		}

		Handlers[req[0]].Handler(privilege, req[1:], writer)
	}

}
