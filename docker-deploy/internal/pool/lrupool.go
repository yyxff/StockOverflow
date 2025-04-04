package pool

import (
	"errors"
	"sync"
)

// generic lru chain
// lru can be upgrade to young/old or lfu
type LruPool[T any] struct {
	// map
	nodePool sync.Map

	// mutex
	mu sync.Mutex

	// lru chain
	head *LruNode[T]
	tail *LruNode[T]

	// max size
	limit int

	// current size
	currentSize int
}

// ==============================public==============================

// new
func newLruPool[T any](limit int) *LruPool[T] {
	if limit < 10 {
		limit = 10
	}
	return &LruPool[T]{
		limit:       limit,
		currentSize: 0,
	}
}

// get a node from pool
func (pool *LruPool[T]) Get(sym string) (*LruNode[T], error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	value, exists := pool.nodePool.Load(sym)
	if !exists {
		return nil, errors.New("no " + sym + " in stock pool")
	}
	node := value.(*LruNode[T])
	pool.touch(node)

	return node, nil
}

// put a node into pool
func (pool *LruPool[T]) Put(node *LruNode[T]) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	_, exists := pool.nodePool.Load(node.symbol)
	if exists {
		return errors.New(node.symbol + " already in stock pool")
	}
	pool.add(node)

	return nil
}

func (pool *LruPool[T]) GetSize() int {
	return pool.currentSize
}

// ==============================private==============================

// add a node to front
func (pool *LruPool[T]) add(node *LruNode[T]) {
	node.next = pool.head
	node.prev = nil

	if pool.head != nil {
		pool.head.prev = node
	}
	pool.head = node

	if pool.tail == nil {
		pool.tail = node
	}

	pool.nodePool.Store(node.symbol, node)
	pool.currentSize++
	pool.updateLRU()
}

// remove a sym
func (pool *LruPool[T]) removeSym(sym string) {
	value, exists := pool.nodePool.Load(sym)
	if !exists {
		return
	}

	node := value.(*LruNode[T])
	pool.removeNode(node)
}

// remove a node
func (pool *LruPool[T]) removeNode(node *LruNode[T]) {

	// if was read or written by other
	// cancel remove

	// if no other is using this node
	// get a write lock then remove it
	if !node.TryLock() {
		return
	}
	prev := node.prev
	next := node.next

	if node == pool.head {
		pool.head = next
	}
	if node == pool.tail {
		pool.tail = prev
	}
	if prev != nil {
		prev.next = next
	}
	if next != nil {
		next.prev = prev
	}

	pool.nodePool.Delete(node.symbol)

	node.Unlock()
	pool.currentSize--
}

// touch node
// move it to front
func (pool *LruPool[T]) touch(node *LruNode[T]) {
	if node == pool.head {
		return
	}
	pool.removeNode(node)
	pool.add(node)
}

// evict node
func (pool *LruPool[T]) evict() {
	if pool.tail != nil {
		pool.removeNode(pool.tail)
	}
}

// check limit
func (pool *LruPool[T]) updateLRU() {
	for pool.currentSize > pool.limit {
		pool.evict()
	}
}
