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
// FITNESS FOR k PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package collections

import (
	"github.com/jakeschurch/instruments"

	"github.com/jakeschurch/collections/internal/linkedlist/holdings"
)

type Portfolio struct {
	Holdings *holdings.Holdings
}

// NewPortfolio constructs a new Portfolio instance.
func NewPortfolio() *Portfolio {
	return &Portfolio{
		Holdings: holdings.New(),
	}
}

// Insert a new holding into a portfolio instance.
func (p *Portfolio) Insert(holding instruments.Holding) {
	p.Holdings.Insert(holding)
}

// Update aggregate holding data from quoted data.
// Returns error if associated instrument not held,
// otherwise returns nil.
func (p *Portfolio) Update(quote instruments.Quote) error {
	return p.Holdings.Update(quote)
}

// Remove a set of holdings from a portfolio.
func (p *Portfolio) Remove(key string) error {
	return p.Holdings.Remove(key)
}
