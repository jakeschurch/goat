// Copyright (c) 2018 Jake Schurch
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
