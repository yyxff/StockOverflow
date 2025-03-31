package pool

import "sync"

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

	// rw lock
	rw sync.RWMutex
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

// get read lock
func (node *StockNode) RLock() {
	node.rw.RLock()
}

// get write lock
func (node *StockNode) WLock() {
	node.rw.Lock()
}

// unlock read
func (node *StockNode) RUnlock() {
	node.rw.RUnlock()
}

// unlock write
func (node *StockNode) Unlock() {
	node.rw.Unlock()
}

// get buyerheap
func (node *StockNode) GetBuyers() *BuyerHeap {
	return &node.buyers
}

// get sellerheap
func (node *StockNode) GetSellers() *SellerHeap {
	return &node.sellers
}
