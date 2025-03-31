package test

import (
	. "StockOverflow/internal/pool"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestBuyerPush(t *testing.T) {
	buyers := &BuyerHeap{}
	buyers.SafePush(NewOrder(decimal.NewFromFloat(1.0), time.Now()))
	buyers.SafePush(NewOrder(decimal.NewFromFloat(2.0), time.Now()))
	buyers.SafePush(NewOrder(decimal.NewFromFloat(0.5), time.Now()))

	x, _ := buyers.SafePop()
	d := x.(Order)
	if d.GetPrice().Cmp(decimal.NewFromFloat(2.0)) != 0 {
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
