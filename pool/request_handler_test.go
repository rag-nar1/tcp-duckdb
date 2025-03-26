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

	threads := 10

	var wg sync.WaitGroup

	for i := 1; i <= threads; i ++ {
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

	for i := 1; i <= threads; i ++ {
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
	for i := 1; i <= threads; i ++ {
		wg.Add(1)
		go func(t *testing.T, rh *pool.RequestHandler, dbname string) {
			req := pool.NewRequest(dbname)
			rh.Push(req)

			connection := <- req.Response
			assert.NotNil(t, connection)
			defer connection.Destroy()

			res := connection.DB().QueryRow("select * from t1 limit 1;")
			assert.NotNil(t, res)
			var id int
			err := res.Scan(&id)
			assert.Nil(t, err)
			assert.Equal(t, 1, id)
			wg.Done()
		}(t, rh, fmt.Sprintf("db%d", i))
	}
	wg.Wait()
}