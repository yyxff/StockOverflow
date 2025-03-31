package pool

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	price decimal.Decimal
	time  time.Time
}

// new
func NewOrder(price decimal.Decimal, time time.Time) *Order {
	return &Order{
		price: price,
		time:  time,
	}
}

// get price
func (order *Order) GetPrice() decimal.Decimal {
	return order.price
}

// get order
func (order *Order) GetTime() time.Time {
	return order.time
}
