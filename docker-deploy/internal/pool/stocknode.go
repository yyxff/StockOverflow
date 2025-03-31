package pool

// a stock node is a trading room for a specific stock
type StockNode struct {

	// pointer in lru
	next *StockNode
	prev *StockNode

	// buyers and sellers heap
	sellers SellerHeap
	buyers  BuyerHeap
}
