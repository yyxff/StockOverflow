package pool

import (
	"sync"
	"sync/atomic"
)

// generic lru node
type LruNode[T any] struct {
	symbol string
	value  T
	prev   *LruNode[T]
	next   *LruNode[T]
	// rw lock
	rw sync.RWMutex
	// pin
	pin atomic.Int64
}

// get read lock
func (node *LruNode[T]) RLock() {
	node.rw.RLock()
	node.pin.Add(1)
}

// get write lock
func (node *LruNode[T]) Lock() {
	node.rw.Lock()
	node.pin.Add(1)
}

// try write lock
func (node *LruNode[T]) TryLock() bool {
	if !node.pin.CompareAndSwap(0, 1) {
		return false
	}
	node.rw.Lock()
	return true
}

// unlock read
func (node *LruNode[T]) RUnlock() {
	node.rw.RUnlock()
	node.pin.Add(-1)
}

// unlock write
func (node *LruNode[T]) Unlock() {
	node.rw.Unlock()
	node.pin.Add(-1)
}

func (node *LruNode[T]) IsEvictable() bool {
	return node.pin.Load() == 0
}
