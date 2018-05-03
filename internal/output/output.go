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
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

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

func (plog *PerformanceLog) OutputResults(format Format, pathName string) {
	var holdingResults = make([][]string, 0)

	for key, _ := range plog.holdings.Holdings.Keys() {
		holdingSlice, _ := plog.holdings.GetSlice(key)
		holdingResults = append(holdingResults, NewHoldingSummary(holdingSlice...).ToSlice())
	}
	switch format {
	case JSON:
		ToJSON(holdingResults, pathName)
	default:
		ToCSV(holdingResults, pathName)
	}
}

// ----------------------------------------------------------------------------

func ToCSV(data [][]string, pathName string) {
	data = append([][]string{GetHeaders()}, data...)

	var file, err = os.Create(pathName)
	if err != nil {
		log.Fatalln(err)
	}

	w := csv.NewWriter(file)
	for _, row := range data {
		w.Write(row)
	}
	w.Flush()
	file.Close()
}

func ToJSON(data [][]string, pathName string) {
	var file, err = os.Create(pathName)
	if err != nil {
		log.Fatalln(err)
	}

	w := io.Writer(file)
	marshalled, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}
	w.Write(marshalled)
	file.Close()
}

type Format uint

const (
	CSV Format = iota
	JSON
)

type holdingSummary struct {
	Name           string
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

func (hs *holdingSummary) ToSlice() []string {
	return []string{
		hs.Name,
		hs.AvgVolume.String(),
		hs.AvgAsk.String(),
		hs.AvgBid.String(),
		hs.MaxAsk.Price.String(),
		hs.MaxAsk.Date.Format(time.RFC1123),
		hs.MinAsk.Price.String(),
		hs.MinAsk.Date.Format(time.RFC1123),
		hs.MaxBid.Price.String(),
		hs.MaxBid.Date.Format(time.RFC1123),
		hs.MinBid.Price.String(),
		hs.MinBid.Date.Format(time.RFC1123),
		string(hs.NumOrderFilled),
		hs.PctReturn.ToPercent(),
		hs.Alpha.ToPercent(),
	}
}
func GetHeaders() []string {
	return []string{
		"Name",
		"Avg. Volume",
		"Avg. Ask Price",
		"Avg. Bid Price",
		"Max. Ask Price",
		"Max. Ask Date",
		"Min. Ask Price",
		"Min. Ask Date",
		"Max. Bid Price",
		"Max. Bid Date",
		"Min. Bid Price",
		"Min. Bid Date",
		"Number of Orders Filled",
		"Percent Return",
		"Alpha",
	}
}

func (hs *holdingSummary) update(h *instruments.Holding) {
	if hs.NumOrderFilled == 0 {
		hs.Name = h.Name
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
