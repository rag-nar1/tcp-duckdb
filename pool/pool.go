package pool

import (
	"container/list"
	"database/sql"
	"errors"

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
	PinCount	atomic.Int32
	Replacer   	*LruReplacer
	PoolLatch	*sync.Mutex
}

type Connection interface {
	DB() *sql.DB
	Destroy()

	GetPinCount() int // for testing only
}

func (d *DBPool) DB() *sql.DB {
	return d.ConnPool
}

func (d *DBPool) GetPinCount() int {
	return int(d.PinCount.Load())
}

func (d *DBPool) Destroy() {
	if d.ConnPool == nil {
		return
	}

	d.PoolLatch.Lock()
	
	d.PinCount.Add(-1)
	if d.PinCount.Load() == 0 {
		d.Replacer.SetEvictable(d.Dbid, true)
	}
	d.PoolLatch.Unlock()

	d = nil
}

func (d *DBPool) Delete() {
	if d == nil || d.ConnPool == nil {
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

func NewPool() *Pool {
	pool := &Pool{
		DB: make([]*DBPool, server.DbPoolSize + 1),
		Ids: make(map[string]uint),
		Free: list.New(),
		Size: 0,
		Replacer: NewLruReplacer(server.ReplacerK),
		Latch: sync.Mutex{},
	}

	for i := 1; i <= int(server.DbPoolSize); i ++ {
		pool.Free.PushBack(uint(i))
	}

	return pool
}


func (p *Pool) Get(dbname string) (Connection, error) {
	p.Latch.Lock()
	defer p.Latch.Unlock()

	dbid, ok := p.Ids[dbname]
	if ok {
		p.DB[dbid].PinCount.Add(1)
		return p.DB[dbid], nil
	}

	var connPool *sql.DB

	if p.Size != server.DbPoolSize { // no evaction needed
		dbid = p.Free.Front().Value.(uint)
		p.Free.Remove(p.Free.Front())
	} else {
		// try to evict
		dbid = p.Replacer.Evict()
		if dbid == InvalidDbId {
			return nil, errors.New(LruReplacerFullErrorStmt)
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
		PinCount: atomic.Int32{},
		Replacer: p.Replacer,
		PoolLatch: &p.Latch,
	}

	if err := p.Replacer.RecordAccess(dbPool.Dbid); err != nil {
		dbPool.Delete()
		return nil, err
	}

	dbPool.PinCount.Add(1)
	p.DB[dbid] = &dbPool
	p.Ids[dbname] = dbid
	p.Size ++

	return &dbPool, nil
}