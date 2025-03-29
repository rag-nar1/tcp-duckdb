package utils

import (
	"bufio"

	"github.com/rag-nar1/tcp-duckdb/pool"
	"github.com/rag-nar1/tcp-duckdb/request_handler"
)


func Write(writer *bufio.Writer, data []byte) error {
	data = append(data, '\n')
	if _, err := writer.Write(data); err != nil {
        return err
    }
    return writer.Flush()
}

func OpenDb(rh *request_handler.RequestHandler, dbname string) (pool.Connection, error){
	req := request_handler.NewRequest(dbname)
	rh.Push(req)

	Dbconn := <- req.Response
	err := <- req.Err
	return Dbconn, err
}