package pool

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	id     uint
	amount uint
	price  decimal.Decimal
	time   time.Time
}

// new
func NewOrder(id uint, amount uint, price decimal.Decimal, time time.Time) *Order {
	return &Order{
		id:     id,
		amount: amount,
		price:  price,
		time:   time,
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

// get ID
func (order *Order) GetID() uint {
	return order.id
}
