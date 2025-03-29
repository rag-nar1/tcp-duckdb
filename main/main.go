package main

import (
	"net"

	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rag-nar1/tcp-duckdb/server"
)

func init() {
	server.Init()
}

func main() {
	// start listing to tcp connections
	listener, err := net.Listen("tcp", server.Serv.Address)
	if err != nil {
		server.Serv.ErrorLog.Fatal(err)
	}
	defer listener.Close()

	server.Serv.InfoLog.Println("listening to " + server.Serv.Address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			server.Serv.ErrorLog.Fatal(err)
		}
		// starts a go routin to handle every connection
		go HandleConnection(conn)
	}
}
