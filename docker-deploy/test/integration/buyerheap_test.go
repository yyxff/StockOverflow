package test

import (
	. "StockOverflow/internal/pool"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestBuyerPush(t *testing.T) {
	buyers := NewBuyerHeap(10)
	buyers.SafePush(NewOrder(1, 1, decimal.NewFromFloat(1.0), time.Now()))
	buyers.SafePush(NewOrder(1, 1, decimal.NewFromFloat(2.0), time.Now()))
	buyers.SafePush(NewOrder(1, 1, decimal.NewFromFloat(0.5), time.Now()))

	x, _ := buyers.SafePop()
	d := x.(Order)
	if d.GetPrice().Cmp(decimal.NewFromFloat(2.0)) != 0 {
		t.Errorf("should get 2.0 but %d", x)
	}

}

func TestBuyerPop(t *testing.T) {
	buyers := NewBuyerHeap(10)
	// heap.Push(1.0)
	// heap.Push(2.0)
	// heap.Push(0.5)

	_, err := buyers.SafePop()
	if err.Error() != "pop from empty heap" {
		t.Errorf("should output error")
	}
}

func TestBuyerTime(t *testing.T) {
	buyers := NewBuyerHeap(10)

	t1 := time.Now()
	t2 := t1.Add(2 * time.Second)
	buyers.SafePush(NewOrder(1, 1, decimal.NewFromFloat(3.0), t1))
	buyers.SafePush(NewOrder(1, 1, decimal.NewFromFloat(3.0), t2))
	buyers.SafePush(NewOrder(1, 1, decimal.NewFromFloat(2.0), time.Now()))

	x, _ := buyers.SafePop()
	d := x.(Order)
	if d.GetPrice().Cmp(decimal.NewFromFloat(3.0)) != 0 {
		t.Errorf("should get 3.0 but %d", x)
	}

	if d.GetTime() != t1 {
		t.Errorf("should get first t1")
	}
}

func TestBuyerTime2(t *testing.T) {
	buyers := NewBuyerHeap(10)

	t1 := time.Now()
	t2 := t1.Add(2 * time.Second)
	buyers.SafePush(NewOrder(1, 1, decimal.NewFromFloat(3.0), t2))
	buyers.SafePush(NewOrder(1, 1, decimal.NewFromFloat(3.0), t1))
	buyers.SafePush(NewOrder(1, 1, decimal.NewFromFloat(2.0), time.Now()))

	x, _ := buyers.SafePop()
	d := x.(Order)
	if d.GetPrice().Cmp(decimal.NewFromFloat(3.0)) != 0 {
		t.Errorf("should get 3.0 but %d", x)
	}

	if d.GetTime() != t1 {
		t.Errorf("should get first t1")
	}
}
