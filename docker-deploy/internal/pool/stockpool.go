package pool

// stock pool with lru
// lru can be upgrade to young/old or lfu
type StockPool struct {
	LruPool[*StockNode]
}

// ==============================public==============================

// new
func NewPool(limit int) *StockPool {
	if limit < 10 {
		limit = 10
	}
	return &StockPool{
		*newLruPool[*StockNode](limit),
	}
}
