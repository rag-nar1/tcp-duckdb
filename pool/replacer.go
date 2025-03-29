/*
	Lrureplacer is an LRU-K Lrureplacer
*/

package pool

import (
	list "container/list"
	atomic "sync/atomic"

	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/rag-nar1/tcp-duckdb/globals"
)

const ( 
	InvalidDbId uint = 0
	InvalidDbidStmt = "invalid dbid %d"
	LruReplacerFullErrorStmt = "replacer is full use Evict"
)

type Node struct {
	Dbid		uint
	Evictable	bool
	Access 		*list.List 	// accesslist for each database hold the time stamp
}

func NewNode(dbid uint) *Node {
	return &Node{
		Dbid: dbid,
		Evictable: false,
		Access: &list.List{},
	}
}

func (n *Node) RecordAccess(time *atomic.Uint64) {
	n.Access.PushBack(time.Load())
	time.Add(1)
}

func (n *Node) Len() uint {
	return uint(n.Access.Len())
}

func (n *Node) Remove() {
	n.Access.Remove(n.Access.Front())
}

func (n *Node) GetKDistance(time *atomic.Uint64, K uint) uint64 {
	if n.Access.Len() < int(K) {
		return math.MaxUint64
	}

	return time.Load() - n.Access.Front().Value.(uint64)
}

func (n *Node) GetMostRecentAccess() uint64 {
	return n.Access.Back().Value.(uint64)
}

type LruReplacer struct {
	Nodes		map[uint]*Node
	Size		uint				// curr size of the Lrureplacer
	K 			uint				// the size of the time window which the Lrureplacer will use to evict db
	Latch		sync.Mutex		// latch to control concurrent use
	CurrTime	atomic.Uint64		// atomic value for shared access to record the time where a database is accessed
}


func NewLruReplacer(k uint) *LruReplacer {
	r := LruReplacer {
		Nodes: make(map[uint]*Node),
		Size: 0,
		K: k,
		Latch: sync.Mutex{},
		CurrTime: atomic.Uint64{},
	}

	return &r
}

func (r *LruReplacer) RecordAccess(dbid uint) error {
	r.Latch.Lock() // lock the latch to record the access
	defer r.Latch.Unlock()
	
	if dbid > globals.DbPoolSize {
		return fmt.Errorf(InvalidDbidStmt, dbid)
	}

	if r.Size == uint(globals.DbPoolSize) {
		return errors.New(LruReplacerFullErrorStmt)
	}

	if _,ok := r.Nodes[dbid]; !ok {
		r.Nodes[dbid] = NewNode(dbid)
	}

	// record access 
	r.Nodes[dbid].RecordAccess(&r.CurrTime)
	if r.Nodes[dbid].Access.Len() > int(r.K) {
		r.Nodes[dbid].Remove()
	}

	return nil
}

func (r *LruReplacer) Evict() (uint) {
	r.Latch.Lock() // lock the latch to evict
	defer r.Latch.Unlock()

	var victim uint = InvalidDbId
	var maxKDistans uint64 = 0
	var mostRecentAccess uint64 = 0

	for dbid, node := range r.Nodes {
		if !node.Evictable {
			continue
		}

		currKDistance := node.GetKDistance(&r.CurrTime, r.K)
		if currKDistance == math.MaxUint64 && mostRecentAccess < node.GetMostRecentAccess() {
			victim = uint(dbid)
			maxKDistans = currKDistance
			mostRecentAccess = node.GetMostRecentAccess()
			continue
		}
		
		if currKDistance > maxKDistans {
			victim = uint(dbid)
			maxKDistans = currKDistance
			mostRecentAccess = node.GetMostRecentAccess()
		}
	}

	if victim != InvalidDbId {
		r.Size --
		delete(r.Nodes, victim)
	}

	return victim
}

func (r *LruReplacer) SetEvictable(dbid uint, evictable bool) {
	r.Latch.Lock() // lock the latch
	defer r.Latch.Unlock()
	if _,ok := r.Nodes[dbid]; !ok {
		return
	}

	if r.Nodes[dbid].Evictable == evictable {
		return
	}

	r.Nodes[dbid].Evictable = evictable
}