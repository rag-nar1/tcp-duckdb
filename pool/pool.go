package pool

import (
	"container/list"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/rag-nar1/TCP-Duckdb/server"
)

type DBPool struct {
	ConnPool	*sql.DB
	Dbname		string
	Dbid		uint
	PinCount	atomic.Uint32
} 

func (d *DBPool) Delete() {
	if d.ConnPool == nil {
		return
	}
	d.ConnPool.Close()
}

type Pool struct {
	DB			[]*DBPool
	Ids 		map[string]uint
	Free		*list.List
	Size		uint
	Replacer 	*LruReplacer
	Latch		sync.Mutex
}

func NewPool(size uint) *Pool {
	pool := &Pool{
		DB: make([]*DBPool, size),
		Ids: make(map[string]uint),
		Free: list.New(),
		Size: 0,
		Replacer: NewLruReplacer(server.ReplacerK),
		Latch: sync.Mutex{},
	}

	for i := 1; i <= int(size); i ++ {
		pool.Free.PushBack(i)
	}

	return pool
}


func (p *Pool) Get(dbname string) (*sql.DB, error) {
	p.Latch.Lock()
	defer p.Latch.Unlock()

	dbid, ok := p.Ids[dbname]
	if ok {
		return p.DB[dbid].ConnPool, nil
	}

	var connPool *sql.DB

	if p.Size != server.DbPoolSize { // no evaction needed
		dbid = p.Free.Front().Value.(uint)
		p.Free.Remove(p.Free.Front())
	} else {
		// try to evict
		dbid = p.Replacer.Evict()
		if dbid == InvalidDbId {
			return nil, fmt.Errorf(InvalidDbidStmt, dbid)
		}
	}

	connPool, err := sql.Open("duckdb", os.Getenv("DBdir") + "users/"+ dbname + ".db")
	if err != nil {
		return nil, err
	}
	p.DB[dbid].Delete()

	dbPool := DBPool{
		ConnPool: connPool,
		Dbname: dbname,
		Dbid: dbid,
		PinCount: atomic.Uint32{},
	}

	if err := p.Replacer.RecordAccess(dbPool.Dbid); err != nil {
		dbPool.Delete()
		return nil, err
	}

	dbPool.PinCount.Add(1)
	p.DB[dbid] = &dbPool
	p.Ids[dbname] = dbid
	p.Size ++

	return dbPool.ConnPool, nil
}