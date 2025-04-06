package test

import (
	"StockOverflow/internal/database"
	"StockOverflow/internal/server"

	. "StockOverflow/internal/pool"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestBuyerPush(t *testing.T) {
	buyers := NewBuyerHeap("SPY", 10, 3)
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(1.0), time.Now()))
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(2.0), time.Now()))
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(0.5), time.Now()))

	x, _ := buyers.SafePop()
	d := x.(Order)
	if d.GetPrice().Cmp(decimal.NewFromFloat(2.0)) != 0 {
		t.Errorf("should get 2.0 but %d", x)
	}

}

func TestBuyerPop(t *testing.T) {
	buyers := NewBuyerHeap("SPY", 10, 3)
	// heap.Push(1.0)
	// heap.Push(2.0)
	// heap.Push(0.5)

	_, err := buyers.SafePop()
	if err.Error() != "pop from empty heap" {
		t.Errorf("should output error")
	}
}

func TestBuyerTime(t *testing.T) {
	buyers := NewBuyerHeap("SPY", 10, 3)

	t1 := time.Now()
	t2 := t1.Add(2 * time.Second)
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(3.0), t1))
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(3.0), t2))
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(2.0), time.Now()))

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
	buyers := NewBuyerHeap("SPY", 10, 3)

	t1 := time.Now()
	t2 := t1.Add(2 * time.Second)
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(3.0), t2))
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(3.0), t1))
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(2.0), time.Now()))

	x, _ := buyers.SafePop()
	d := x.(Order)
	if d.GetPrice().Cmp(decimal.NewFromFloat(3.0)) != 0 {
		t.Errorf("should get 3.0 but %d", x)
	}

	if d.GetTime() != t1 {
		t.Errorf("should get first t1")
	}
}

func TestUpdate(t *testing.T) {
	buyers := NewBuyerHeap("SPY", 10, 3)

	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(3.0), time.Now()))
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(3.0), time.Now()))
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(2.0), time.Now()))
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(2.0), time.Now()))
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(2.0), time.Now()))
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(2.0), time.Now()))
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(2.0), time.Now()))
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(2.0), time.Now()))
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(2.0), time.Now()))

	if buyers.Len() != 9 {
		t.Errorf("should get 9 but %d", buyers.Len())
	}

	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(2.0), time.Now()))
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(2.0), time.Now()))

	if buyers.Len() != 5 {
		t.Errorf("should get 5 but %d", buyers.Len())
	}

}

func TestRefill(t *testing.T) {
	dbm := database.DatabaseMaster{
		ConnStr: server.GetDBConnStr(),
		DbName:  "stockoverflowtest",
	}

	// Connect to database
	// logger.Println("Connecting to database...")
	dbm.Connect()
	dbm.CreateDB()
	dbm.Init()

	buyers := NewBuyerHeap("SPY", 10, 3)
	buyers.SetDB(dbm.Db)

	t1 := time.Now()
	t2 := t1.Add(2 * time.Second)
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(3.0), t2))
	buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(3.0), t1))
	// buyers.SafePush(NewOrder("1", 1, decimal.NewFromFloat(2.0), time.Now()))

	x, _ := buyers.SafePop()
	d := x.(Order)
	if d.GetPrice().Cmp(decimal.NewFromFloat(300.0)) != 0 {
		t.Errorf("should get 300.0 but %d", x)
	}

	// if d.GetTime() != t1 {
	// 	t.Errorf("should get first t1")
	// }
}
