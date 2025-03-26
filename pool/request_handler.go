package pool

type Request struct {
	Dbname   string
	Response chan Connection
	Err      chan error
}

func NewRequest(dbname string) *Request {
	return &Request{
		Dbname:   dbname,
		Response: make(chan Connection, 1),
		Err:      make(chan error, 1),
	}
}

type RequestHandler struct {
	Requests   chan *Request
	PoolManger *Pool
}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{
		Requests:   make(chan *Request),
		PoolManger: NewPool(),
	}
}

func HandleRequest(rh *RequestHandler, curr *Request) {
	for {
		connection, err := rh.PoolManger.Get(curr.Dbname)
		if err != nil && err.Error() == LruReplacerFullErrorStmt {
			continue
		}

		if err != nil {
			curr.Err <- err
			break
		}

		curr.Response <- connection
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
