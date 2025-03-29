package pool_test

import (
	"sync"
	"testing"

	"github.com/rag-nar1/tcp-duckdb/pool"
	"github.com/rag-nar1/tcp-duckdb/globals"
	"github.com/stretchr/testify/assert"
)


func TestReaplacerBasic(t *testing.T) {
	var dbid1 uint = 1
	var dbid2 uint = 2
	replacer := pool.NewLruReplacer(2)
	replacer.RecordAccess(dbid1)
	replacer.RecordAccess(dbid1)
	replacer.RecordAccess(dbid2)
	replacer.SetEvictable(dbid1, true)
	replacer.SetEvictable(dbid2, true)
	evicted := replacer.Evict()
	assert.Equal(t, dbid2, evicted)
}

func TestReaplacerMain(t *testing.T) {
	var dbid1 uint = 1
	var dbid2 uint = 2
	replacer := pool.NewLruReplacer(2)
	replacer.RecordAccess(dbid1)
	replacer.RecordAccess(dbid1)
	replacer.RecordAccess(dbid2)
	replacer.RecordAccess(dbid2)
	replacer.SetEvictable(dbid1, true)
	replacer.SetEvictable(dbid2, true)
	evicted := replacer.Evict()
	assert.Equal(t, dbid1, evicted)

	replacer.SetEvictable(dbid2, false)
	replacer.RecordAccess(dbid1)
	replacer.RecordAccess(dbid1)
	evicted = replacer.Evict()
	assert.Equal(t, pool.InvalidDbId, evicted)

	replacer.SetEvictable(dbid1, true)
	evicted = replacer.Evict()
	assert.Equal(t, dbid1, evicted)
}

func TestReaplacerConcurruncy(t *testing.T) {
	replacer := pool.NewLruReplacer(3)
	var wg sync.WaitGroup          // WaitGroup to synchronize goroutines
	for i := 1; i <= int(globals.DbPoolSize); i ++ {
		wg.Add(1)
		go func(dbid uint, replacer *pool.LruReplacer) {
			for time := 0; time < 3; time ++ {
				if err := replacer.RecordAccess(uint(dbid)); err != nil {
					panic(err)
				}
			}
			if dbid % 2 == 1 {
				replacer.SetEvictable(dbid, true)
			}
			wg.Done()
		}(uint(i), replacer)
	}
	wg.Wait()
	
	replacer.SetEvictable(1, false)
	err := replacer.RecordAccess(1)
	assert.Nil(t, err)

	evicted := replacer.Evict()
	assert.Equal(t, uint(1), evicted % 2)
	assert.NotEqual(t, uint(1) , evicted)

	for i := 1; i <= int(globals.DbPoolSize); i ++ {
		wg.Add(1)
		go func(dbid uint, replacer *pool.LruReplacer) {
			replacer.SetEvictable(dbid, false)
			wg.Done()
		}(uint(i), replacer)
	}
	wg.Wait()
	
	replacer.SetEvictable(1, true)
	evicted = replacer.Evict()
	assert.Equal(t, uint(1) , evicted)
}

