package pool

import (
	"container/heap"
	"errors"
	"fmt"
)

type SellerHeap []Order

func (heap SellerHeap) Len() int { return len(heap) }
func (heap SellerHeap) Less(i, j int) bool {
	diff := heap[i].price.Cmp(heap[j].price)
	if diff < 0 {
		return true
	} else if diff == 0 {
		return heap[i].time.Before(heap[j].time)
	}
	return false
}
func (heap SellerHeap) Swap(i, j int) { heap[i], heap[j] = heap[j], heap[i] }

// push
func (sellers *SellerHeap) Push(ele interface{}) {
	order, ok := ele.(Order)
	if !ok {
		fmt.Println("invalid argument in Push")
	}
	*sellers = append(*sellers, order)
}

// safe push
func (sellers *SellerHeap) SafePush(ele *Order) error {
	heap.Push(sellers, *ele)
	return nil
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
