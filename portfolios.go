package goat

import (
	"sync"

	"github.com/jakeschurch/instruments"

	"github.com/jakeschurch/collections"
)

type Portfolio struct {
	*collections.Portfolio
	cash instruments.Amount
	sync.RWMutex
}

func NewPortfolio(cash instruments.Amount) *Portfolio {
	return &Portfolio{
		Portfolio: collections.NewPortfolio(),
		cash:      cash,
	}
}

// checkSells to see if we can create an orders.
// if holdings are empty GetSlice will return error.
func (p *Portfolio) checkSells(quote instruments.Quote, algos ...Algorithm) ([]*instruments.Order, error) {
	var sells = make([]*instruments.Order, 0)
	var holdings, err = p.GetSlice(quote.Name)
	if err != nil {
		return sells, err
	}

	var i, j = 0, 0
	for j <= len(algos) {
		if order, ok := algos[j].Sell(quote, holdings[i]); ok {
			sells = append(sells, order)
		}
		if i < len(holdings) {
			i++
		}
		if j < len(algos) {
			j++
			i = 0
		} else {
			break
		}
	}
	return sells, nil
}

type Benchmark struct {
	*collections.Portfolio
	sync.RWMutex
}
