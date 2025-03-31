package pool

import (
	"container/heap"
	"errors"

	"github.com/shopspring/decimal"
)

type SellerHeap []decimal.Decimal

func (heap SellerHeap) Len() int           { return len(heap) }
func (heap SellerHeap) Less(i, j int) bool { return heap[i].Cmp(heap[j]) < 0 }
func (heap SellerHeap) Swap(i, j int)      { heap[i], heap[j] = heap[j], heap[i] }

// push
func (heap *SellerHeap) Push(ele interface{}) {
	var d decimal.Decimal
	switch t := ele.(type) {
	case decimal.Decimal:
		d = t
	case float64:
		d = decimal.NewFromFloat(t)
	}
	*heap = append(*heap, d)
}

// pop
func (heap *SellerHeap) Pop() interface{} {
	old := *heap
	ele := old[old.Len()-1]
	*heap = old[0 : old.Len()-1]
	return ele
}

// safe Pop
func (sellers *SellerHeap) SafePop() (interface{}, error) {
	if sellers.Len() == 0 {
		return nil, errors.New("pop from empty heap")
	}
	return heap.Pop(sellers), nil
}
