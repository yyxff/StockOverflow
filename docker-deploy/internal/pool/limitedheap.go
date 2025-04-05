package pool

import (
	"container/heap"
	"database/sql"
	"errors"
)

type LimitedHeap[T any] struct {
	*Heap[T]
	// size
	maxSize uint
	minSize uint

	// db
	db *sql.DB
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

// set db
func (h *LimitedHeap[T]) SetDB(db *sql.DB) {
	h.db = db
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
	} else if uint(h.Len()) < h.minSize {
		h.pullFromDB()
	}
}

// pull enough data from db
func (h *LimitedHeap[T]) pullFromDB() error {
	if h.db == nil {
		return errors.New("no db connected now!")
	}

	// do sql

	return nil
}
