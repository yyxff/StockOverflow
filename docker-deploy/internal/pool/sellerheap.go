package pool

type SellerHeap struct {
	*LimitedHeap[Order]
}

func NewSellerHeap(maxSize uint, minSize uint) *SellerHeap {
	return &SellerHeap{
		&LimitedHeap[Order]{
			NewHeap(lessMin),
			maxSize,
			minSize,
			nil,
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
