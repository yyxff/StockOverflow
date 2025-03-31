package test

import (
	. "StockOverflow/internal/pool"
	"container/heap"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestPush(t *testing.T) {
	sellers := &SellerHeap{}
	sellers.SafePush(NewOrder(decimal.NewFromFloat(1.0), time.Now()))
	sellers.SafePush(NewOrder(decimal.NewFromFloat(0.5), time.Now()))
	sellers.SafePush(NewOrder(decimal.NewFromFloat(2.0), time.Now()))

	x := heap.Pop(sellers)
	d := x.(Order)
	if d.GetPrice().Cmp(decimal.NewFromFloat(0.5)) != 0 {
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
