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

package datafeed

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/jakeschurch/instruments"
	"github.com/pkg/errors"
	"github.com/timpalpant/go-iex"
	"github.com/timpalpant/go-iex/iextp/tops"
)

const (
	bufferSize uint = 150 // at the moment: arbitrary number.
	histFormat      = "20060102"
)

var (
	// initialStart is first day of data given from IEX.
	initialStart = time.Date(2016, 12, 12, 0, 0, 0, 0, time.UTC)
	// ErrInvalidDate throws error on invalid date argument given.
	ErrInvalidDate = errors.New("Invalid date argument given")
)

type producerSlice []*dataProducer

func (p producerSlice) Len() int {
	return len(p)
}

func (p producerSlice) Less(i, j int) bool {
	return p[i].date.Before(p[j].date)
}

func (p producerSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func historicalData(start, end time.Time) producerSlice {
	var result = make(producerSlice, 0)
	// startStr, endStr := start.Format(histFormat), end.Format(histFormat)

	client := iex.NewClient(&http.Client{})
	// Get historical data dumps available for input date argument.
	histData, err := client.GetAllAvailableHIST()
	if err != nil {
		panic(err)
	} else if len(histData) == 0 {
		panic(fmt.Errorf("Found %v available data feeds", len(histData)))
	}
	for key, val := range histData {
		dateKey, _ := time.Parse(histFormat, key)
		// reduce data down to selected range.
		if (dateKey.After(start) || dateKey.Equal(start)) && (dateKey.Before(end) || dateKey.Equal(end)) {
			result = append(result, newProducer(dateKey, val))
		}
	}
	sort.Sort(result)
	return result
}

// Produce populates a channel given as input to allow for a pipeline of *instruments.Quote structs.
func Produce(outCh chan<- *instruments.Quote, start, end time.Time) error {
	// check if start before inital data start date
	if !start.After(initialStart) {
		return errors.Wrapf(ErrInvalidDate, "Invalid Start Date: %v, must be %v or after", start, initialStart)
	}
	// check if invalid date range given
	if !start.Before(end) {
		return errors.Wrapf(ErrInvalidDate, "Start Date: %v cannot be before End Date: %v", start, end)
	}
	// check if end date greater than today
	if !end.Before(time.Now()) {
		return errors.Wrapf(ErrInvalidDate, "End Date: %v cannot be after today's date: %v", end, time.Now())
	}
	var wg = sync.WaitGroup{}

	workers := historicalData(start, end)

	go func(outCh chan<- *instruments.Quote, producers producerSlice) {
		i := 0
		for i >= len(producers) {
			val, ok := <-producers[i].ch
			if !ok {
				close(producers[i].ch)
				wg.Done()
				i++
				continue
			}
			outCh <- val
		}
		close(outCh)
	}(outCh, workers)

	for i := range workers {

		wg.Add(1)
		go func(producer *dataProducer) {
			producer.collect()
		}(workers[i])
	}
	wg.Wait()
	return nil
}

type dataProducer struct {
	date time.Time
	hist []*iex.HIST
	ch   chan *instruments.Quote
}

func newProducer(date time.Time, data []*iex.HIST) *dataProducer {
	return &dataProducer{
		date: date,
		hist: data,
		ch:   make(chan *instruments.Quote, bufferSize),
	}
}

func (d *dataProducer) collect() {

	// Fetch the pcap dump for that date and iterate through its messages.
	resp, err := http.Get(d.hist[0].Link)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	packetDataSource, err := iex.NewPacketDataSource(resp.Body)
	if err != nil {
		panic(err)
	}
	pcapScanner := iex.NewPcapScanner(packetDataSource)

	// Write each quote update message to dataProducer's channel.
	for {
		msg, err := pcapScanner.NextMessage()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		switch msg := msg.(type) {
		case *tops.QuoteUpdateMessage:

			if msg.AskPrice == 0 || msg.AskSize == 0 ||
				msg.BidPrice == 0 || msg.BidSize == 0 {
				continue
			}
			d.ch <- &instruments.Quote{
				Name: msg.Symbol,
				// TEMP / TODO: change instruments.Volume as int32 to comply with go-iex
				Bid:       instruments.NewQuotedMetric(msg.BidPrice, float64(msg.BidSize)),
				Ask:       instruments.NewQuotedMetric(msg.AskPrice, float64(msg.AskSize)),
				Timestamp: msg.Timestamp,
			}
		}
	}
}
