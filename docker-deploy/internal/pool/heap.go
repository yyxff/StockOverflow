package pool

import (
	"container/heap"
)

type Heap[T any] struct {
	data []T
	cmp  func(a, b T) bool
}

// implement heap.Interface
func (h *Heap[T]) Len() int           { return len(h.data) }
func (h *Heap[T]) Less(i, j int) bool { return h.cmp(h.data[i], h.data[j]) }
func (h *Heap[T]) Swap(i, j int)      { h.data[i], h.data[j] = h.data[j], h.data[i] }

func (h *Heap[T]) Push(x any) {
	h.data = append(h.data, x.(T))
}

func (h *Heap[T]) Pop() any {
	old := h.data
	n := len(old)
	x := old[n-1]
	h.data = old[:n-1]
	return x
}

// new a heap by cmp
func NewHeap[T any](cmp func(a, b T) bool) *Heap[T] {
	h := &Heap[T]{cmp: cmp}
	heap.Init(h)
	return h
}
