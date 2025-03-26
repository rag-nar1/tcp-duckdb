package pool_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/joho/godotenv"
	"github.com/rag-nar1/TCP-Duckdb/pool"
	"github.com/stretchr/testify/assert"
)

func TestRHBasic(t *testing.T) {
	t.Cleanup(CleanUp)
	if err := godotenv.Load("../.env"); err != nil {
		panic(err)
	}

	rh := pool.NewRequestHandler()
	go rh.Spin()

	req := pool.NewRequest("db1")
	rh.Push(req)

	connection := <- req.Response
	assert.NotNil(t, connection)
	_,err := connection.DB().Exec("create table t1(id int);")
	assert.Nil(t, err)
}

func TestRHConcurruncy(t *testing.T) {
	t.Cleanup(CleanUp)
	if err := godotenv.Load("../.env"); err != nil {
		panic(err)
	}
	rh := pool.NewRequestHandler()
	go rh.Spin()

	var wg sync.WaitGroup

	for i := 1; i <= 5; i ++ {
		wg.Add(1)
		go func(t *testing.T, rh *pool.RequestHandler, dbname string) {
			req := pool.NewRequest(dbname)
			rh.Push(req)

			connection := <- req.Response
			assert.NotNil(t, connection)
			defer connection.Destroy()

			_,err := connection.DB().Exec("create table t1(id int);")
			assert.Nil(t, err)
			wg.Done()
		}(t, rh, fmt.Sprintf("db%d", i))
	}
	wg.Wait()

	for i := 1; i <= 5; i ++ {
		wg.Add(1)
		go func(t *testing.T, rh *pool.RequestHandler, dbname string) {
			req := pool.NewRequest(dbname)
			rh.Push(req)

			connection := <- req.Response
			assert.NotNil(t, connection)
			defer connection.Destroy()

			_,err := connection.DB().Exec("insert into t1(id) values(1);")
			assert.Nil(t, err)
			wg.Done()
		}(t, rh, fmt.Sprintf("db%d", i))
	}
	wg.Wait()
}