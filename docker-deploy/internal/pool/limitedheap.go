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

	// symbol
	symbol string

	// refill fn
	refillFn func(db *sql.DB, symbol string, heapType string, size int) []T

	// heap type
	heapType string
}

// safe Pop
func (h *LimitedHeap[T]) SafePop() (interface{}, error) {
	if h.Len() == 0 {
		return nil, errors.New("pop from empty heap")
	}
	h.CheckMin()
	result := heap.Pop(h)
	return result, nil
}

// safe push
func (h *LimitedHeap[T]) SafePush(ele *Order) error {
	heap.Push(h, *ele)
	h.checkMax()
	return nil
}

// set db
func (h *LimitedHeap[T]) SetDB(db *sql.DB) {
	h.db = db
}

// update heap to keep it small
func (h *LimitedHeap[T]) checkMax() {
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

// update heap to keep it big
func (h *LimitedHeap[T]) CheckMin() {
	if uint(h.Len()) < h.minSize {
		h.pullFromDB()
		heap.Init(h)
	}
}

// pull enough data from db
func (h *LimitedHeap[T]) pullFromDB() error {
	if h.db == nil {
		return errors.New("no db connected now")
	}

	h.data = h.refillFn(h.db, h.symbol, "buyer", int((h.maxSize+h.minSize)/2))

	// update minsize
	if h.Len() < int(h.minSize) {
		h.minSize = uint(h.Len() / 2)
	}
	return nil
}
