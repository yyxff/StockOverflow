package pool

import (
	"container/heap"
	"errors"

	"github.com/shopspring/decimal"
)

type BuyerHeap []decimal.Decimal

func (buyers BuyerHeap) Len() int           { return len(buyers) }
func (buyers BuyerHeap) Less(i, j int) bool { return buyers[i].Cmp(buyers[j]) > 0 }
func (buyers BuyerHeap) Swap(i, j int)      { buyers[i], buyers[j] = buyers[j], buyers[i] }

// push
func (buyers *BuyerHeap) Push(ele interface{}) {
	var d decimal.Decimal
	switch t := ele.(type) {
	case decimal.Decimal:
		d = t
	case float64:
		d = decimal.NewFromFloat(t)
	}
	*buyers = append(*buyers, d)
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
