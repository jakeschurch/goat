package collections

import (
	"github.com/jakeschurch/collections/internal/linkedlist/orders"
	"github.com/jakeschurch/instruments"
)

type OrderBook struct {
	Buys  *orders.Orders
	Sells *orders.Orders
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Buys:  orders.New(),
		Sells: orders.New(),
	}
}

func (ob *OrderBook) Insert(o *instruments.Order) {
	switch o.Buy {
	case true:
		ob.Buys.Insert(o)
	case false:
		ob.Sells.Insert(o)
	}
}

func (ob *OrderBook) Remove(o *instruments.Order) (ok bool) {

	if o.Buy {
		if data := ob.Buys.Search(o); data != nil {
			return ob.Buys.Remove(data)
		}
		return false
	}
	if data := ob.Sells.Search(o); data != nil {
		return ob.Sells.Remove(data)
	}
	return false
}

func (ob *OrderBook) GetBuys(key string) ([]*instruments.Order, error) {
	return ob.Buys.GetSlice(key)
}

func (ob *OrderBook) GetSells(key string) ([]*instruments.Order, error) {
	return ob.Sells.GetSlice(key)
}
