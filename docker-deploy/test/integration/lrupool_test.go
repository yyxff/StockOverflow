package test

import (
	. "StockOverflow/internal/pool"
	"testing"
)

func TestPut(t *testing.T) {
	pool := NewPool(10)

	if pool.GetSize() != 00 {
		t.Errorf("size should be 0 but %d", pool.GetSize())
	}
	pool.Put(NewStockNode("abc", 10))
	pool.Put(NewStockNode("1", 10))
	pool.Put(NewStockNode("2", 10))
	pool.Put(NewStockNode("3", 10))
	pool.Put(NewStockNode("4", 10))
	pool.Put(NewStockNode("5", 10))
	pool.Put(NewStockNode("6", 10))
	pool.Put(NewStockNode("7", 10))
	pool.Put(NewStockNode("8", 10))
	pool.Put(NewStockNode("9", 10))

	if pool.GetSize() != 10 {
		t.Errorf("size should be 10 but %d", pool.GetSize())
	}
}

func TestEvict(t *testing.T) {
	pool := NewPool(10)

	pool.Put(NewStockNode("abc", 10))
	pool.Put(NewStockNode("1", 10))
	pool.Put(NewStockNode("2", 10))
	pool.Put(NewStockNode("3", 10))
	pool.Put(NewStockNode("4", 10))
	pool.Put(NewStockNode("5", 10))
	pool.Put(NewStockNode("6", 10))
	pool.Put(NewStockNode("7", 10))
	pool.Put(NewStockNode("8", 10))
	pool.Put(NewStockNode("9", 10))

	if pool.GetSize() != 10 {
		t.Errorf("size should be 10 but %d", pool.GetSize())
	}

	pool.Put(NewStockNode("10", 10))

	if pool.GetSize() != 10 {
		t.Errorf("size should be 10 but %d", pool.GetSize())
	}

}

func TestGet(t *testing.T) {
	pool := NewPool(10)

	node := NewStockNode("abc", 10)
	pool.Put(node)
	pool.Put(NewStockNode("1", 10))
	pool.Put(NewStockNode("2", 10))
	pool.Put(NewStockNode("3", 10))
	pool.Put(NewStockNode("4", 10))
	pool.Put(NewStockNode("5", 10))
	pool.Put(NewStockNode("6", 10))
	pool.Put(NewStockNode("7", 10))
	pool.Put(NewStockNode("8", 10))

	getnode, err := pool.Get("abc")
	if err != nil {
		t.Errorf("failed to get ele")
	}
	if getnode != node {
		t.Errorf("shoud get %p but %p", node, getnode)
	}
}
