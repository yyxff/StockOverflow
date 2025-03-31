package pool

import (
	"errors"
	"sync"
)

// stock pool with lru
// lru can be upgrade to young/old or lfu
type Pool struct {
	// map
	stockPool sync.Map

	// lru chain
	head *StockNode
	tail *StockNode

	// max size
	limit int

	// current size
	currentSize int
}

// ==============================public==============================

// new
func (pool *Pool) NewPool(limit int) *Pool {
	if limit < 10 {
		limit = 10
	}
	return &Pool{
		limit:       limit,
		currentSize: 0,
	}
}

// get a node from pool
func (pool *Pool) Get(sym string) (*StockNode, error) {
	value, exists := pool.stockPool.Load(sym)
	if !exists {
		return nil, errors.New("no " + sym + " in stock pool")
	}
	node := value.(*StockNode)
	pool.touch(node)
	return node, nil
}

// put a node into pool
func (pool *Pool) Put(node *StockNode) error {
	_, exists := pool.stockPool.Load(node.symbol)
	if exists {
		return errors.New(node.symbol + " already in stock pool")
	}
	pool.add(node)
	return nil
}

// ==============================private==============================

// add a node to front
func (pool *Pool) add(node *StockNode) {
	node.next = pool.head
	node.prev = nil

	if pool.head != nil {
		pool.head.prev = node
	}
	pool.head = node

	if pool.tail == nil {
		pool.tail = node
	}

	pool.stockPool.Store(node.symbol, node)
	pool.currentSize++
	pool.updateLRU()
}

// remove a sym
func (pool *Pool) removeSym(sym string) {
	value, exists := pool.stockPool.Load(sym)
	if !exists {
		return
	}

	node := value.(*StockNode)
	pool.removeNode(node)
}

// remove a node
func (pool *Pool) removeNode(node *StockNode) {
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

	pool.stockPool.Delete(node.symbol)
	pool.currentSize--
}

// touch node
// move it to front
func (pool *Pool) touch(node *StockNode) {
	pool.removeNode(node)
	pool.add(node)
}

// evict node
func (pool *Pool) evict() {
	if pool.tail != nil {
		pool.removeNode(pool.tail)
	}
}

// check limit
func (pool *Pool) updateLRU() {
	for pool.currentSize > pool.limit {
		pool.evict()
	}
}
