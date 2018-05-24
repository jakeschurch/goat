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

package worker

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jakeschurch/instruments"
)

type Config struct {
	Name, Bid, BidSz, Ask, AskSz, Timestamp uint8
	Timeunit                                string
	Date                                    time.Time
}
type Worker struct {
	DataChan chan *instruments.Quote
	config   Config
}

func New(c Config) *Worker {
	return &Worker{
		config:   c,
		DataChan: make(chan *instruments.Quote),
	}
}

// IDEA: Have Produce take custom callback type as input

func (w *Worker) Produce(r io.ReadSeeker, headers bool) {
	var done = make(chan struct{})
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	if headers {
		scanner.Scan() // for headers...
	}
	go func(ch chan<- *instruments.Quote) {
		for scanner.Scan() {
			line := scanner.Text()

			// Check to see if error has been thrown or
			if err := scanner.Err(); err != nil {
				if err == io.EOF {
					break
				} else {
					log.Fatalln(err)
				}
			}

			var data, err = w.consume(strings.Split(line, "|"))

			if err == nil {
				ch <- data
			} else {
				fmt.Println(err)
			}
		}
		close(ch)
		close(done)
	}(w.DataChan)

	<-done
}

var ErrParseRecord = errors.New("record could not be parsed correctly")

func (w *Worker) consume(record []string) (*instruments.Quote, error) {
	var quote = new(instruments.Quote)
	quote.Ask, quote.Bid = new(instruments.QuotedMetric), new(instruments.QuotedMetric)

	quote.Name = record[w.config.Name]

	qbid, bidErr := strconv.ParseFloat(record[w.config.Bid], 64)
	if qbid == 0 || bidErr != nil {
		return quote, ErrParseRecord
	}
	quote.Bid.Price = instruments.NewPrice(qbid)

	qbidSz, bidSzErr := strconv.ParseFloat(record[w.config.BidSz], 64)
	if qbidSz == 0 || bidSzErr != nil {
		return quote, ErrParseRecord
	}
	quote.Bid.Volume = instruments.NewVolume(qbidSz)

	qask, askErr := strconv.ParseFloat(record[w.config.Ask], 64)
	if qask == 0 || askErr != nil {
		return quote, ErrParseRecord
	}
	quote.Ask.Price = instruments.NewPrice(qask)

	qaskSz, askSzErr := strconv.ParseFloat(record[w.config.AskSz], 64)
	if qaskSz == 0 || askSzErr != nil {
		return quote, ErrParseRecord
	}
	quote.Ask.Volume = instruments.NewVolume(qaskSz)

	tickDuration, timeErr := time.ParseDuration(record[w.config.Timestamp] + w.config.Timeunit)
	if timeErr != nil {
		return quote, ErrParseRecord
	}
	quote.Timestamp = w.config.Date.Add(tickDuration)

	return quote, nil
}
