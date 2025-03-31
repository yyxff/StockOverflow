package test

import (
	. "StockOverflow/internal/pool"
	"container/heap"
	"testing"

	"github.com/shopspring/decimal"
)

func TestBuyerPush(t *testing.T) {
	buyers := &BuyerHeap{}
	heap.Push(buyers, 1.0)
	heap.Push(buyers, 2.0)
	heap.Push(buyers, 0.5)

	x, _ := buyers.SafePop()
	d := x.(decimal.Decimal)
	if d.Cmp(decimal.NewFromFloat(2.0)) != 0 {
		t.Errorf("should get 2.0 but %d", x)
	}

}

func TestBuyerPop(t *testing.T) {
	buyers := BuyerHeap{}
	// heap.Push(1.0)
	// heap.Push(2.0)
	// heap.Push(0.5)

	_, err := buyers.SafePop()
	if err.Error() != "pop from empty heap" {
		t.Errorf("should output error")
	}
}
