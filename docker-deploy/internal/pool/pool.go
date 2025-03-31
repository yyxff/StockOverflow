package pool

// stock pool with lru
// lru can be upgrade to young/old or lfu
type Pool struct {
	pool map[string]*StockNode
	head *StockNode
	tail *StockNode
}
