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
	"time"

	"github.com/pkg/errors"
)

var ErrInvalidTx = errors.New("invalid transaction type given")

type Holding struct {
	Name   string
	Volume Volume
	Buy    TxMetric
	Sell   TxMetric
}

// Buy creates a new Holding from transaction data.
func Buy(tx Transaction) (*Holding, error) {
	if !tx.Buy {
		return nil, errors.Wrap(ErrInvalidTx, "wanted buy, got sell")
	}
	return &Holding{
		Name:   tx.Name,
		Volume: tx.Volume,
		Buy:    TxMetric{Price: tx.Price, Date: tx.Timestamp},
	}, nil
}

// SellOff a number of securities from transaction data.
func (h *Holding) SellOff(tx Transaction) (*Holding, error) {
	if tx.Buy || h.Volume < tx.Volume {
		return nil, errors.Wrap(ErrInvalidTx, "wanted sell, got buy")
	}
	h.Volume -= tx.Volume
	return h, nil
}

// TxMetric is an associated price-date metric pair.
type TxMetric struct {
	Price Price
	Date  time.Time
}

// ----------------------------------------------------------------------------

// Summary is a metric summary for a particular financial instrument.
type Summary struct {
	Name             string
	N                uint
	Volume           Volume
	AvgBid, AvgAsk   *Price
	LastBid, LastAsk *SummaryMetric
	MaxBid, MaxAsk   *SummaryMetric
	MinBid, MinAsk   *SummaryMetric
}

// NewSummary returns a new summary instance.
// Summary is constructed from fields supplied by a Holding instance.
func NewSummary(h Holding) *Summary {
	metric := &SummaryMetric{Price: h.Buy.Price, Date: h.Buy.Date}
	return &Summary{
		Name: h.Name, N: 0, Volume: h.Volume, AvgBid: &h.Buy.Price, AvgAsk: &h.Buy.Price,
		MaxBid: metric, MaxAsk: metric, MinBid: metric, MinAsk: metric, LastAsk: metric, LastBid: metric,
	}
}

func (s *Summary) UpdateMetrics(qBid, qAsk Price, t time.Time) {
	if qBid == 0 || qAsk == 0 {
		return
	}
	s.LastBid = &SummaryMetric{Price: qBid, Date: t}
	s.LastAsk = &SummaryMetric{Price: qAsk, Date: t}

	s.AvgBid.Avg(s.N, qBid)
	s.AvgAsk.Avg(s.N, qAsk)
	s.N++

	s.MaxAsk.Max(qAsk, t)
	s.MinAsk.Min(qAsk, t)

	s.MaxBid.Max(qBid, t)
	s.MinBid.Min(qBid, t)
}

// ----------------------------------------------------------------------------

// SummaryMetric records an associated price-date pair.
type SummaryMetric struct {
	Price Price
	Date  time.Time
}

// Avg is calculated from an old price value,
// the number of times the average has been calculated, and the new quote price.
func (p *Price) Avg(n uint, quotePrice Price) Price {
	numerator := *p*Price(n) + quotePrice
	newAvg := numerator / (Price(n) + 1)
	p = &newAvg

	return newAvg
}

// Max is calculated from a SummaryMetric's price field, and a new quoted price.
func (s *SummaryMetric) Max(quotePrice Price, timestamp time.Time) Price {
	if s.Price <= quotePrice {
		s.Price = quotePrice
	}
	return s.Price
}

// Min is calculated from a SummaryMetric's price field, and a new quoted price.
func (s *SummaryMetric) Min(quotePrice Price, timestamp time.Time) Price {
	if s.Price >= quotePrice {
		s.Price = quotePrice
	}
	return s.Price
}
