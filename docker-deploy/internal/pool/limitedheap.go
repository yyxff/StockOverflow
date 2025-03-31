package pool

import (
	"container/heap"
	"errors"
)

type LimitedHeap[T any] struct {
	*Heap[T]
	maxSize uint
	// minSize uint
}

// safe Pop
func (h *LimitedHeap[T]) SafePop() (interface{}, error) {
	if h.Len() == 0 {
		return nil, errors.New("pop from empty heap")
	}
	return heap.Pop(h), nil
}

// safe push
func (h *LimitedHeap[T]) SafePush(ele *Order) error {
	heap.Push(h, *ele)
	h.updateHeap()
	return nil
}

func (h *LimitedHeap[T]) updateHeap() {
	if uint(h.Len()) > h.maxSize {
		//  todo set minsize
		h.data = h.data[:h.maxSize/2]
		heap.Init(h)
	}
}
