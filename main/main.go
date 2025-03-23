package main

import (
	"net"

	_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3"
	global "github.com/rag-nar1/TCP-Duckdb/server"
)

func init() {
	global.Init()
}

func main() {
	// start listing to tcp connections
	listener, err := net.Listen("tcp", global.Serv.Address)
	if err != nil {
		global.Serv.ErrorLog.Fatal(err)
	}
	defer listener.Close()

	global.Serv.InfoLog.Println("listening to " + global.Serv.Address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			global.Serv.ErrorLog.Fatal(err)
		}
		// starts a go routin to handle every connection
		go HandleConnection(conn)
	}
}
