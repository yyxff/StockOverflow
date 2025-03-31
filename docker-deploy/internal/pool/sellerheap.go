package pool

import (
	"github.com/shopspring/decimal"
)

type SellerHeap []decimal.Decimal

func (heap SellerHeap) Len() int           { return len(heap) }
func (heap SellerHeap) Less(i, j int) bool { return heap[i].Cmp(heap[j]) < 0 }
func (heap SellerHeap) Swap(i, j int)      { heap[i], heap[j] = heap[j], heap[i] }

// push
func (heap *SellerHeap) Push(ele interface{}) {
	*heap = append(*heap, ele.(decimal.Decimal))
}

// pop
func (heap *SellerHeap) Pop() interface{} {
	old := *heap
	ele := old[old.Len()-1]
	*heap = old[0 : old.Len()-1]
	return ele
}
