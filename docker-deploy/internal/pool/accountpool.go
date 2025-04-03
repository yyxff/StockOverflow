package pool

// stock pool with lru
// lru can be upgrade to young/old or lfu
type AccountPool struct {
	LruPool[*AccountNode]
}

// ==============================public==============================

// new
func NewAccountPool(limit int) *AccountPool {
	if limit < 10 {
		limit = 10
	}
	return &AccountPool{
		*newLruPool[*AccountNode](limit),
	}
}
