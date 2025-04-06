package pool

type SellerHeap struct {
	*OrderHeap
}

func NewSellerHeap(symbol string, maxSize uint, minSize uint) *SellerHeap {
	return &SellerHeap{
		NewOrderHeap(symbol, maxSize, minSize, lessMin, "seller"),
	}
}

// implement compare
func lessMin(i, j Order) bool {
	diff := i.price.Cmp(j.price)
	if diff < 0 {
		return true
	} else if diff == 0 {
		return i.time.Before(j.time)
	}
	return false
}
