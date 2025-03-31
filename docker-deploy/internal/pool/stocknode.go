package pool

// a stock node is a trading room for a specific stock
type StockNode struct {
	// name
	symbol string

	// pointer in lru
	next *StockNode
	prev *StockNode

	// buyers and sellers heap
	sellers SellerHeap
	buyers  BuyerHeap
}

// new
func NewStockNode(symbol string) *StockNode {
	return &StockNode{
		symbol:  symbol,
		next:    nil,
		prev:    nil,
		sellers: SellerHeap{},
		buyers:  BuyerHeap{},
	}
}
