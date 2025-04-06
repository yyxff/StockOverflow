package pool

// a stock node is a trading room for a specific stock
type StockNode struct {

	// buyers and sellers heap
	sellers SellerHeap
	buyers  BuyerHeap
}

// new
func NewStockNode(symbol string, limit uint) *LruNode[*StockNode] {
	return &LruNode[*StockNode]{
		symbol: symbol,
		value: &StockNode{
			sellers: *NewSellerHeap(symbol, limit, 1),
			buyers:  *NewBuyerHeap(symbol, limit, 1),
		},
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
