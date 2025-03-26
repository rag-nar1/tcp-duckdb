package request_handler

import "github.com/rag-nar1/TCP-Duckdb/pool"

type Request struct {
	Dbname   string
	Response chan pool.Connection
	Err      chan error
}

func NewRequest(dbname string) *Request {
	return &Request{
		Dbname:   dbname,
		Response: make(chan pool.Connection, 1),
		Err:      make(chan error, 1),
	}
}

type RequestHandler struct {
	Requests   chan *Request
	PoolManger *pool.Pool
}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{
		Requests:   make(chan *Request),
		PoolManger: pool.NewPool(),
	}
}

func HandleRequest(rh *RequestHandler, curr *Request) {
	for {
		connection, err := rh.PoolManger.Get(curr.Dbname)
		if err != nil && err.Error() == pool.LruReplacerFullErrorStmt {
			continue
		}

		if err != nil {
			close(curr.Response)
			curr.Err <- err
			break
		}

		curr.Response <- connection
		close(curr.Err)
		break
	}
}

func (rh *RequestHandler) Spin() {
	for curr := range rh.Requests {
		go HandleRequest(rh, curr)
	}
}

func (rh *RequestHandler) Push(req *Request) {
	rh.Requests <- req
}
