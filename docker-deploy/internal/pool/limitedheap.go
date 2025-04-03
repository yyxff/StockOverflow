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

// update heap to keep it small
func (h *LimitedHeap[T]) updateHeap() {
	if uint(h.Len()) > h.maxSize {
		//  todo set minsize
		newSize := h.maxSize / 2
		var i uint
		var data []T
		for i = 0; i < newSize; i++ {
			ele, err := h.SafePop()
			x := ele.(T)
			if err != nil {

			} else {
				data = append(data, x)
			}

		}
		h.data = data
		heap.Init(h)
	}
}
