package pool

import (
	"StockOverflow/internal/database"
	"database/sql"
	"time"
)

type OrderHeap struct {
	*LimitedHeap[Order]
}

func NewOrderHeap(symbol string, maxSize uint, minSize uint, cmp func(a, b Order) bool, heapType string) *OrderHeap {
	return &OrderHeap{
		&LimitedHeap[Order]{
			NewHeap(cmp),
			maxSize,
			minSize,
			nil,
			symbol,
			refillFn,
			heapType,
		},
	}
}

func refillFn(db *sql.DB, symbol string, heapType string, size int) []Order {

	if db == nil {
		return nil
	}

	// do sql
	orders, err := database.GetOpenOrdersBySymbolForHeap(db, symbol, heapType, size)
	if err != nil {
		return nil
	}

	// new heap data
	var data []Order
	for _, order := range orders {
		neworder := NewOrder(order.ID, uint(order.Amount.IntPart()), order.Price, time.Unix(order.Timestamp, 0))
		data = append(data, *neworder)
	}
	return data
}
