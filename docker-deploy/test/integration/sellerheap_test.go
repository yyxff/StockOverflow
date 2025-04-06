package test

import (
	. "StockOverflow/internal/pool"
	"container/heap"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestPush(t *testing.T) {
	sellers := NewSellerHeap("SPY", 10, 3)
	sellers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(1.0), time.Now()))
	sellers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(0.5), time.Now()))
	sellers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(2.0), time.Now()))

	x := heap.Pop(sellers)
	d := x.(Order)
	if d.GetPrice().Cmp(decimal.NewFromFloat(0.5)) != 0 {
		t.Errorf("should get 0.5 but %d", x)
	}

}

func TestPop(t *testing.T) {
	sellers := NewSellerHeap("SPY", 10, 3)
	// heap.Push(1.0)
	// heap.Push(2.0)
	// heap.Push(0.5)

	_, err := sellers.SafePop()
	if err.Error() != "pop from empty heap" {
		t.Errorf("should output error")
	}
}

func TestSellerTime(t *testing.T) {
	sellers := NewSellerHeap("SPY", 10, 3)

	t1 := time.Now()
	t2 := t1.Add(2 * time.Second)
	sellers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(1.0), t1))
	sellers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(1.0), t2))
	sellers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(2.0), time.Now()))

	x := heap.Pop(sellers)
	d := x.(Order)
	if d.GetPrice().Cmp(decimal.NewFromFloat(1.0)) != 0 {
		t.Errorf("should get 0.5 but %d", x)
	}

	if d.GetTime() != t1 {
		t.Errorf("should get first t1")
	}
}

func TestSellerTime2(t *testing.T) {
	sellers := NewSellerHeap("SPY", 10, 3)

	t1 := time.Now()
	t2 := t1.Add(2 * time.Second)
	sellers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(1.0), t2))
	sellers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(1.0), t1))
	sellers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(2.0), time.Now()))

	x := heap.Pop(sellers)
	d := x.(Order)
	if d.GetPrice().Cmp(decimal.NewFromFloat(1.0)) != 0 {
		t.Errorf("should get 0.5 but %d", x)
	}

	if d.GetTime() != t1 {
		t.Errorf("should get first t1")
	}
}
