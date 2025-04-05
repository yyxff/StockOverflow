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
	// rw sync.RWMutex

	// mutex lock
	mu sync.Mutex

	// pin
	pin atomic.Int64
}

// get mutex lock
func (node *LruNode[T]) Lock() {
	node.mu.Lock()
	node.pin.Add(1)
}

// try write lock
func (node *LruNode[T]) TryLock() bool {
	if !node.pin.CompareAndSwap(0, 1) {
		return false
	}
	node.mu.Lock()
	return true
}

// unlock write
func (node *LruNode[T]) Unlock() {
	node.mu.Unlock()
	node.pin.Add(-1)
}

func (node *LruNode[T]) IsEvictable() bool {
	return node.pin.Load() == 0
}

// getter for value
func (node *LruNode[T]) GetValue() T {
	return node.value
}
