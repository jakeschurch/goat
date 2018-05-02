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
	"fmt"
	"time"

	"github.com/jakeschurch/instruments/internal/ordering"
)

// Order stores logic for transacting a stock.
type Order struct {
	Name string
	quotedMetric
	filled Volume

	Buy       bool
	Status    Status
	Logic     Logic
	timestamp time.Time
	ticker    *ordering.OrderTicker
}

func (o *Order) String() string {
	return fmt.Sprintf("\nName: %v\nPrice: %d\nVolume: %d\ntimestamp:%s", o.Name, o.Price, o.Volume, o.timestamp)
}

// NewOrder instantiates a new order struct.
func NewOrder(name string, buy bool, logic Logic, price Price, volume Volume, timestamp time.Time) *Order {
	return &Order{
		Name:         name,
		Buy:          buy,
		quotedMetric: quotedMetric{Price: price, Volume: volume},
		Logic:        logic,
		Status:       Open,
		timestamp:    timestamp,

		ticker: ordering.NewOrderTicker(),
		filled: 0,
	}
}

func (o *Order) timestampTx() time.Time {
	return o.timestamp.Add(o.ticker.Duration())
}

// Transact a fulfillment of an order; yielding a new transaction struct.
func (o *Order) Transact(price Price, volume Volume) *Transaction {
	o.filled += volume
	o.Volume -= volume
	return &Transaction{
		Name:         o.Name,
		Buy:          o.Buy,
		quotedMetric: quotedMetric{price, volume},
		Timestamp:    o.timestampTx(),
	}
}

// ----------------------------------------------------------------------------

// Transaction represents a fulfillment of a financial order.
type Transaction struct {
	Name string
	Buy  bool
	quotedMetric
	Timestamp time.Time
}

// Status variables refer to a status of an order's execution.
type Status int

const (
	// Open indicates that an order has not been transacted.
	Open Status = iota // 0
	// Closed indicates that an order has been transacted.
	Closed
	// Cancelled indicates than an order was closed, but order was not transacted.
	Cancelled
)

// Logic is used to identify when the order should be executed.
type Logic int

const (
	Market Logic = iota // 0
	Limit
)
