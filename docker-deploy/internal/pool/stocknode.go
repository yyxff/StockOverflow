package pool

// a stock node is a trading room for a specific stock
type StockNode struct {

	// buyers and sellers heap
	sellers SellerHeap
	buyers  BuyerHeap
}

// new
func NewStockNode(symbol string) *LruNode[*StockNode] {
	return &LruNode[*StockNode]{
		symbol: symbol,
		value:  &StockNode{},
	}
}

// get buyerheap
func (node *StockNode) GetBuyers() *BuyerHeap {
	return &node.buyers
}

// get sellerheap
func (node *StockNode) GetSellers() *SellerHeap {
	return &node.sellers
}
