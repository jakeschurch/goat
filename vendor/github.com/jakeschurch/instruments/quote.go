// Copyright (c) 2017 Jake Schurch
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

package instruments

import (
	"errors"
	"time"
)

var (
	ErrZeroValue = errors.New("zero value found")
	ErrNilValue  = errors.New("nil value found")
)

// ----------------------------------------------------------------------------

// Quote reflects a static state of a security.
type Quote struct {
	Name      string
	Bid, Ask  *QuotedMetric
	Timestamp time.Time
}

// TEMP: come back and change vol to int or something idk
func (q *Quote) FillOrder(price Price, vol Volume, buy bool, logic Logic) *Order {
	return NewOrder(q.Name, buy, logic, price, vol, q.Timestamp)
}

// TotalAsk returns a Amount representation of the total Ask amount of a quote.
func (q *Quote) TotalAsk() (Amount, error) {
	if q.Ask == nil {
		return 0, ErrNilValue
	}
	return q.Ask.Total()
}

// TotalBid returns a Amount representation of the total Bid amount of a quote.
func (q *Quote) TotalBid() (Amount, error) {
	if q.Bid == nil {
		return 0, ErrNilValue
	}
	return q.Bid.Total()
}

// ----------------------------------------------------------------------------

// A QuotedMetric is a representation of a Price with an associated Volume.
type QuotedMetric struct {
	Price
	Volume
}

func NewQuotedMetric(price, volume float64) *QuotedMetric {
	return &QuotedMetric{Price: NewPrice(price), Volume: NewVolume(volume)}
}

// Total returns the product of a Price and a Volume.
func (q *QuotedMetric) Total() (a Amount, err error) {
	if a = Amount(q.Price * Price(q.Volume)); a == 0 {
		return 0, ErrZeroValue
	}
	return a, err
}
