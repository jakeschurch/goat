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

package output

import (
	"github.com/jakeschurch/collections"
	"github.com/jakeschurch/instruments"
)

// PerformanceLog tracks closed orders and holdings.
type PerformanceLog struct {
	orders   *collections.OrderBook
	holdings *collections.Portfolio
}

func (plog *PerformanceLog) AddOrders(orders ...*instruments.Order) {
	for i := range orders {
		plog.orders.Insert(orders[i])
	}
}

func (plog *PerformanceLog) AddHoldings(holdings ...*instruments.Holding) {
	for i := range holdings {
		plog.holdings.Insert(*holdings[i])
	}
}

// ----------------------------------------------------------------------------

type Format uint

const (
	CSV Format = iota
	JSON
)

type holdingSummary struct {
	AvgVolume      instruments.Volume   `json:"AvgVolume,omitempty"`
	AvgAsk         instruments.Price    `json:"AvgAsk,omitempty"`
	AvgBid         instruments.Price    `json:"AvgBid,omitempty"`
	MaxAsk         instruments.TxMetric `json:"MaxAsk,omitempty"`
	MinAsk         instruments.TxMetric `json:"MinAsk,omitempty"`
	MaxBid         instruments.TxMetric `json:"MaxBid,omitempty"`
	MinBid         instruments.TxMetric `json:"MinBid,omitempty"`
	NumOrderFilled uint                 `json:"NumOrderFilled,omitempty"`
	PctReturn      instruments.Amount   `json:"pctReturn,omitempty"`
	Alpha          instruments.Amount   `json:"alpha,omitempty"`
}

func (hs *holdingSummary) update(h *instruments.Holding) {
	if hs.NumOrderFilled == 0 {
		hs.AvgVolume = h.Volume
		hs.AvgAsk = h.Buy.Price
		hs.AvgBid = h.Sell.Price
		hs.MaxAsk = h.Buy
		hs.MinAsk = h.Buy
		hs.MaxBid = h.Sell
		hs.MinBid = h.Sell
	} else {
		hs.AvgVolume = (hs.AvgVolume*instruments.Volume(hs.NumOrderFilled) + h.Volume) / instruments.Volume(hs.NumOrderFilled+1)
		hs.AvgAsk = (hs.AvgAsk*instruments.Price(hs.NumOrderFilled) + h.Buy.Price) / instruments.Price(hs.NumOrderFilled+1)
		hs.AvgBid = (hs.AvgAsk*instruments.Price(hs.NumOrderFilled) + h.Buy.Price) / instruments.Price(hs.NumOrderFilled+1)

		if h.Buy.Price >= hs.MaxAsk.Price {
			hs.MaxAsk = h.Buy
		}
		if h.Buy.Price < hs.MinAsk.Price {
			hs.MinAsk = h.Buy
		}
		if h.Sell.Price >= hs.MaxBid.Price {
			hs.MaxBid = h.Sell
		}
		if h.Sell.Price < hs.MinBid.Price {
			hs.MinBid = h.Sell
		}
	}
	hs.NumOrderFilled++
}

func NewHoldingSummary(holdings ...*instruments.Holding) *holdingSummary {
	var hs = new(holdingSummary)
	var pctReturns = make([]instruments.Amount, len(holdings))

	divideAmts := func(top, bottom instruments.Price) instruments.Amount {
		return instruments.Amount((top*200 + bottom) / (bottom * 2))
	}

	for i := range holdings {
		hs.update(holdings[i])
		pctReturns = append(
			pctReturns, divideAmts(holdings[i].Sell.Price-holdings[i].Buy.Price, holdings[i].Buy.Price))
	}
	var holdingReturn instruments.Amount
	for _, pctReturn := range pctReturns {
		holdingReturn += pctReturn
	}
	hs.PctReturn = holdingReturn / instruments.Amount(len(pctReturns)+1)

	return hs
}
