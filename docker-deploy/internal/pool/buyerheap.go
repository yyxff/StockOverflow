package pool

import (
	"github.com/shopspring/decimal"
)

type BuyerHeap []decimal.Decimal

func (heap BuyerHeap) Len() int           { return len(heap) }
func (heap BuyerHeap) Less(i, j int) bool { return heap[i].Cmp(heap[j]) > 0 }
func (heap BuyerHeap) Swap(i, j int)      { heap[i], heap[j] = heap[j], heap[i] }

// push
func (heap *BuyerHeap) Push(ele interface{}) {
	*heap = append(*heap, ele.(decimal.Decimal))
}

// pop
func (heap *BuyerHeap) Pop() interface{} {
	old := *heap
	ele := old[old.Len()-1]
	*heap = old[0 : old.Len()-1]
	return ele
}
