package main

import "sync"

type Account struct {
	ID        string               `xml:"id"`
	Balance   float64              `xml:"balance"`
	Positions map[string]*Position `xml:"positions"`
	Mutex     sync.Mutex
}

type Position struct {
	Symbol string  `xml:"symbol"`
	Amount float64 `xml:"amount"`
}

type Symbol struct {
	Symbol string `xml:"symbol"`
}

type Order struct {
	ID         string  `xml:"id"`
	AccountID  string  `xml:"account_id"`
	Symbol     string  `xml:"symbol"`
	Amount     float64 `xml:"amount"`
	LimitPrice float64 `xml:"limit_price"`
	Status     string  `xml:"status"`
	Remaining  float64 `xml:"remaining"`
	Timestamp  int64   `xml:"timestamp"`
}
