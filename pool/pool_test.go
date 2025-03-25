package pool_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/joho/godotenv"
	"github.com/rag-nar1/TCP-Duckdb/pool"
	"github.com/rag-nar1/TCP-Duckdb/server"
	"github.com/stretchr/testify/assert"
)
func CleanUp() {
	files, err := filepath.Glob("../storge/users/*")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			panic(err)
		}
	}
}

func TestCleanUp(t *testing.T) {
	files, err := filepath.Glob("../storge/users/*")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			panic(err)
		}
	}
}
func TestPoolBasic(t *testing.T) {
	t.Cleanup(CleanUp)
	if err := godotenv.Load("../.env"); err != nil {
		panic(err)
	}

	connPool := pool.NewPool()

	db1, err := connPool.Get("db1")
	assert.Nil(t, err)
	_, err = connPool.Get("db2")
	assert.Nil(t, err)
	_, err = connPool.Get("db3")
	assert.Nil(t, err)
	_, err = connPool.Get("db4")
	assert.Nil(t, err)
	_, err = connPool.Get("db5")
	assert.NotNil(t, err)

	db1.Destroy()
	_, err = connPool.Get("db5")
	assert.Nil(t, err)
}

func TestPoolMain(t *testing.T) {
	t.Cleanup(CleanUp)
	if err := godotenv.Load("../.env"); err != nil {
		panic(err)
	}

	connPool := pool.NewPool()

	db1, err := connPool.Get("db1")
	assert.Nil(t, err)
	_, err = connPool.Get("db2")
	assert.Nil(t, err)
	_, err = connPool.Get("db3")
	assert.Nil(t, err)
	_, err = connPool.Get("db4")
	assert.Nil(t, err)
	_, err = connPool.Get("db5")
	assert.NotNil(t, err)

	db1.Destroy()
	_, err = connPool.Get("db5")
	assert.Nil(t, err)
	db5, err := connPool.Get("db5")
	assert.Nil(t, err)
	assert.Equal(t, 2, db5.GetPinCount())
}

func TestPoolConcurrunct(t *testing.T) {
	t.Cleanup(CleanUp)
	if err := godotenv.Load("../.env"); err != nil {
		panic(err)
	}

	connPool := pool.NewPool()
	var wg sync.WaitGroup

	for i := 1; i <= int(server.DbPoolSize); i ++ {
		wg.Add(1)
		go func (t *testing.T, dbid int, connPool *pool.Pool) {
			db, err := connPool.Get(fmt.Sprintf("db%d", dbid))
			assert.Nil(t, err)
			assert.Equal(t, 1, db.GetPinCount())
			defer func() {
				db.Destroy()
				wg.Done()
			}()
		}(t, i, connPool)
	}

	wg.Wait()
}