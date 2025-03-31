package pool

type SellerHeap struct {
	*LimitedHeap[Order]
}

func NewSellerHeap() *SellerHeap {
	return &SellerHeap{&LimitedHeap[Order]{NewHeap(lessMin)}}
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
