package pool

import "github.com/shopspring/decimal"

type AccountNode struct {
	ID      string
	Balance decimal.Decimal
}

// new
func NewAccountNode(symbol string) *LruNode[*AccountNode] {
	return &LruNode[*AccountNode]{
		symbol: symbol,
		value:  &AccountNode{},
	}
}

// get id
func (node *AccountNode) GetId() string {
	return node.ID
}

// get new
func (node *AccountNode) GetBalance() decimal.Decimal {
	return node.Balance
}
