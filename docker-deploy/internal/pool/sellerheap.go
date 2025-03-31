package pool

type SellerHeap struct {
	*LimitedHeap[Order]
}

func NewSellerHeap(limit uint) *SellerHeap {
	return &SellerHeap{
		&LimitedHeap[Order]{
			NewHeap(lessMin),
			limit,
		},
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
