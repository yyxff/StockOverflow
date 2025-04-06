package pool

type BuyerHeap struct {
	*OrderHeap
}

func NewBuyerHeap(symbol string, maxSize uint, minSize uint) *BuyerHeap {
	return &BuyerHeap{
		NewOrderHeap(symbol, maxSize, minSize, lessMax, "buyer"),
	}
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
