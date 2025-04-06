package pool

import (
	"time"

	"github.com/shopspring/decimal"
)

type Orderer interface {
	GetID() int
}
type Order struct {
	id     string
	amount uint
	price  decimal.Decimal
	time   time.Time
}

// new
func NewOrder(id string, amount uint, price decimal.Decimal, time time.Time) *Order {
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
func (order *Order) GetID() string {
	return order.id
}
