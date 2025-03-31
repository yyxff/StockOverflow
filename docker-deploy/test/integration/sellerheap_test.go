package test

import (
	. "StockOverflow/internal/pool"
	"container/heap"
	"testing"

	"github.com/shopspring/decimal"
)

func TestPush(t *testing.T) {
	sellers := &SellerHeap{}
	heap.Push(sellers, 1.0)
	heap.Push(sellers, 0.5)
	heap.Push(sellers, 2.0)

	x := heap.Pop(sellers)
	d := x.(decimal.Decimal)
	if d.Cmp(decimal.NewFromFloat(0.5)) != 0 {
		t.Errorf("should get 0.5 but %d", x)
	}

}

func TestPop(t *testing.T) {
	sellers := SellerHeap{}
	// heap.Push(1.0)
	// heap.Push(2.0)
	// heap.Push(0.5)

	_, err := sellers.SafePop()
	if err.Error() != "pop from empty heap" {
		t.Errorf("should output error")
	}
}
