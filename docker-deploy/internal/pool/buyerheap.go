package pool

import (
	"container/heap"
	"errors"
	"fmt"
)

type BuyerHeap []Order

func (buyers BuyerHeap) Len() int { return len(buyers) }
func (buyers BuyerHeap) Less(i, j int) bool {
	diff := buyers[i].price.Cmp(buyers[j].price)
	if diff > 0 {
		return true
	} else if diff == 0 {
		return buyers[i].time.Before(buyers[j].time)
	}
	return false
}
func (buyers BuyerHeap) Swap(i, j int) { buyers[i], buyers[j] = buyers[j], buyers[i] }

// push
func (buyers *BuyerHeap) Push(ele interface{}) {
	order, ok := ele.(Order)
	if !ok {
		fmt.Println("invalid argument in Push")
	}
	*buyers = append(*buyers, order)
}

// safe push
func (buyers *BuyerHeap) SafePush(ele *Order) error {
	heap.Push(buyers, *ele)
	return nil
}

// pop
func (buyers *BuyerHeap) Pop() interface{} {
	old := *buyers
	ele := old[old.Len()-1]
	*buyers = old[0 : old.Len()-1]
	return ele
}

// safe Pop
func (buyers *BuyerHeap) SafePop() (interface{}, error) {
	if buyers.Len() == 0 {
		return nil, errors.New("pop from empty heap")
	}
	return heap.Pop(buyers), nil
}
