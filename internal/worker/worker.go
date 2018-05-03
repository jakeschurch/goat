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
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jakeschurch/instruments"
)

type Config struct {
	Name, Bid, BidSz, Ask, AskSz, Timestamp uint8
	Timeunit                                string
	Date                                    time.Time
}
type Worker struct {
	dataChan chan []string
	config   Config
}

func New(wc Config) *Worker {
	return &Worker{
		config: wc,
	}
}

func (worker *Worker) Run(outChan chan<- *instruments.Quote, r io.ReadSeeker) {
	var lineCount int
	var wg sync.WaitGroup
	wg.Add(2)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lineCount++
	}
	r.Seek(0, 0)

	worker.dataChan = make(chan []string, lineCount)

	go func() {
		for {
			data, ok := <-worker.dataChan
			if !ok {
				if len(worker.dataChan) == 0 {
					close(outChan)
					break
				}
			}
			quote, err := worker.consume(data)
			if quote != nil && err == nil {
				outChan <- quote
			}
		}
		defer wg.Done()
	}()
	go worker.produce(r, &wg)

	wg.Wait()
}

func (worker *Worker) produce(r io.ReadSeeker, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	scanner.Scan() // for headers...
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
		record := strings.Split(line, "|")
		if len(record) >= 4 {
			worker.dataChan <- record
		}
	}
	close(worker.dataChan)
	log.Println("done reading from file")
}

var ErrParseRecord = errors.New("record could not be parsed correctly")

func (worker *Worker) consume(record []string) (*instruments.Quote, error) {
	var quote = &instruments.Quote{}

	quote.Name = record[worker.config.Name]

	qbid, bidErr := strconv.ParseFloat(record[worker.config.Bid], 64)
	if qbid == 0 || bidErr != nil {
		return quote, ErrParseRecord
	}
	quote.Bid.Price = instruments.NewPrice(qbid)

	qbidSz, bidSzErr := strconv.ParseFloat(record[worker.config.BidSz], 64)
	if qbidSz == 0 || bidSzErr != nil {
		return quote, ErrParseRecord
	}
	quote.Bid.Volume = instruments.NewVolume(qbidSz)

	qask, askErr := strconv.ParseFloat(record[worker.config.Ask], 64)
	if qask == 0 || askErr != nil {
		return quote, ErrParseRecord
	}
	quote.Ask.Price = instruments.NewPrice(qask)

	qaskSz, askSzErr := strconv.ParseFloat(record[worker.config.AskSz], 64)
	if qaskSz == 0 || askSzErr != nil {
		return quote, ErrParseRecord
	}
	quote.Ask.Volume = instruments.NewVolume(qaskSz)

	tickDuration, timeErr := time.ParseDuration(record[worker.config.Timestamp] + worker.config.Timeunit)
	if timeErr != nil {
		return quote, ErrParseRecord
	}
	quote.Timestamp = worker.config.Date.Add(tickDuration)
	return quote, nil
}
