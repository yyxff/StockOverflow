package pool

type BuyerHeap struct {
	*LimitedHeap[Order]
}

func NewBuyerHeap() *BuyerHeap {
	return &BuyerHeap{&LimitedHeap[Order]{NewHeap(lessMax)}}
}

// implement compare
func lessMax(i, j Order) bool {
	diff := i.price.Cmp(j.price)
	if diff > 0 {
		return true
	} else if diff == 0 {
		return i.time.Before(j.time)
	}
	return false
}

func NewSellerHeap() *Heap[Order] {
	return NewHeap(lessMin)
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
