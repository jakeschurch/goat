package goat

import (
	"errors"

	"github.com/jakeschurch/collections"
	"github.com/jakeschurch/instruments"
)

var ErrLowVolume = errors.New("not enough volueme to fill order")

type OrderManager struct {
	*collections.OrderBook
}

func NewOrderManager() *OrderManager {
	return &OrderManager{
		OrderBook: collections.NewOrderBook(),
	}
}

func (o *OrderManager) Add(order *instruments.Order) {
	o.Insert(order)
}

func (o *OrderManager) Sell(order *instruments.Order, port *Portfolio) ([]*instruments.Transaction, error) {
	var TXs = make([]*instruments.Transaction, 0)
	var sellVol instruments.Volume

	list, err := port.Holdings.Get(order.Name)
	if err != nil {
		return TXs, err
	}

	// If order volume cannot be filled, return error.
	if list.Volume < order.Volume {
		return TXs, ErrLowVolume
	}

	list.Lock()
	// Check to see if we still have holdings
	x, err := list.Peek()
	if err != nil {
		list.Unlock()
		return TXs, err
	}
	// Loop over list elements until
	for x != nil {
		switch x.Volume < order.Volume {
		case true:
			sellVol = x.Volume
		case false:
			sellVol = order.Volume
		}
		// Create new transaction from order.
		tx := order.Transact(order.Price, sellVol)

		// Update portfolio's cash value if appropriate.
		if amt, err := tx.Total(); err == nil {
			port.cash += amt
		}

		// Apply transaction logic to x's Holding.
		x.Holding.SellOff(*tx)

		// Append new tx to TXs slice.
		TXs = append(TXs, tx)

		if x.Volume == 0 {
			list.Unlock()
			list.Pop()
		}
		x = x.Prev()
	}
	list.Unlock()

	return TXs, nil
}

func (o *OrderManager) Buy(order *instruments.Order, port *Portfolio) ([]*instruments.Transaction, error) {
	var TXs = make([]*instruments.Transaction, 0)
	var orderAmt, err = order.Total()
	if err != nil {
		return TXs, err
	}
	// If order volume cannot be filled, return error.
	if port.cash > orderAmt {
		return TXs, ErrLowVolume
	}

	// Create new transaction from order.
	tx := order.Transact(order.Price, order.Volume)

	// Update portfolio's cash value if appropriate.
	if amt, err := tx.Total(); err == nil {
		port.cash -= amt
	}
	// Apply transaction logic to buy NewHolding.
	h, _ := instruments.Buy(*tx)
	port.Insert(*h)
	// Append new tx to TXs slice.
	TXs = append(TXs, tx)

	return TXs, nil
}
