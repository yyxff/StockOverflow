package exchange

import (
	"StockOverflow/internal/pool"
	"sync"
)

// SymbolManager manages order books and locks for each symbol
type SymbolManager struct {
	orderBooks  map[string]*OrderBook
	symbolLocks map[string]*sync.RWMutex
	mapMutex    sync.RWMutex
}

// OrderBook contains buy and sell heaps for a symbol
type OrderBook struct {
	BuyOrders  *pool.BuyerHeap
	SellOrders *pool.SellerHeap
}

// NewSymbolManager creates a new symbol manager
func NewSymbolManager() *SymbolManager {
	return &SymbolManager{
		orderBooks:  make(map[string]*OrderBook),
		symbolLocks: make(map[string]*sync.RWMutex),
	}
}

// GetOrCreateSymbolLock gets or creates a lock for the specified symbol
func (sm *SymbolManager) GetOrCreateSymbolLock(symbol string) *sync.RWMutex {
	// Check if lock exists - read lock is sufficient here
	sm.mapMutex.RLock()
	lock, exists := sm.symbolLocks[symbol]
	sm.mapMutex.RUnlock()

	if !exists {
		// Lock not found, need to create one - use write lock
		sm.mapMutex.Lock()
		defer sm.mapMutex.Unlock()

		// Check again in case it was created between our checks
		lock, exists = sm.symbolLocks[symbol]
		if !exists {
			lock = &sync.RWMutex{}
			sm.symbolLocks[symbol] = lock
		}
	}

	return lock
}

// GetOrCreateOrderBook gets or creates an order book for the specified symbol
func (sm *SymbolManager) GetOrCreateOrderBook(symbol string) *OrderBook {
	// Check if order book exists
	sm.mapMutex.RLock()
	book, exists := sm.orderBooks[symbol]
	sm.mapMutex.RUnlock()

	if !exists {
		// Create new order book with empty heaps
		sm.mapMutex.Lock()
		defer sm.mapMutex.Unlock()

		// Check again in case it was created between our checks
		book, exists = sm.orderBooks[symbol]
		if !exists {
			book = &OrderBook{
				BuyOrders:  &pool.BuyerHeap{},
				SellOrders: &pool.SellerHeap{},
			}
			sm.orderBooks[symbol] = book
		}
	}

	return book
}

// LockSymbol locks a symbol for exclusive access
func (sm *SymbolManager) LockSymbol(symbol string) {
	lock := sm.GetOrCreateSymbolLock(symbol)
	lock.Lock()
}

// UnlockSymbol unlocks a symbol
func (sm *SymbolManager) UnlockSymbol(symbol string) {
	sm.mapMutex.RLock()
	lock, exists := sm.symbolLocks[symbol]
	sm.mapMutex.RUnlock()

	if exists {
		lock.Unlock()
	}
}

// RLockSymbol locks a symbol for shared access
func (sm *SymbolManager) RLockSymbol(symbol string) {
	lock := sm.GetOrCreateSymbolLock(symbol)
	lock.RLock()
}

// RUnlockSymbol unlocks a symbol from shared access
func (sm *SymbolManager) RUnlockSymbol(symbol string) {
	sm.mapMutex.RLock()
	lock, exists := sm.symbolLocks[symbol]
	sm.mapMutex.RUnlock()

	if exists {
		lock.RUnlock()
	}
}
