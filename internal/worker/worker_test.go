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
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/jakeschurch/goat/internal/config"
	"github.com/jakeschurch/instruments"
)

func mockConfig() Config {
	return Config{
		Timestamp: 0, Name: 1, Bid: 2, BidSz: 3, Ask: 4, AskSz: 5,
		Timeunit: "ns", Date: time.Date(2017, time.August, 14, 0, 0, 0, 0, time.UTC),
	}
}

func mockWc() Config {
	conf := config.ReadConfig("././example/config.json")
	var _, date, _ = conf.FileInfo()

	// Setup Worker & WorkerConfig
	return Config{
		Name: conf.File.Columns.Ticker,
		Bid:  conf.File.Columns.Bid, BidSz: conf.File.Columns.BidSize,
		Ask: conf.File.Columns.Ask, AskSz: conf.File.Columns.AskSize,
		Timestamp: conf.File.Columns.Timestamp, Date: date,
		Timeunit: conf.File.TimestampUnit,
	}
}

func testConfigEquality(t1, t2 *Worker) bool {
	if t1.config.Timestamp != t2.config.Timestamp ||
		t1.config.Name != t2.config.Name ||
		t1.config.Bid != t2.config.Bid ||
		t1.config.BidSz != t2.config.BidSz ||
		t1.config.Ask != t2.config.Ask ||
		t1.config.AskSz != t2.config.AskSz {
		return false
	}
	return true
}

func TestNew(t *testing.T) {
	type args struct {
		c Config
	}
	tests := []struct {
		name string
		args args
		want *Worker
	}{
		{"base case", args{mockConfig()}, New(mockConfig())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.c); !testConfigEquality(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorker_consume(t *testing.T) {
	type args struct {
		record []string
	}
	tests := []struct {
		name    string
		worker  *Worker
		args    args
		want    *instruments.Quote
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.worker.consume(tt.args.record)
			if (err != nil) != tt.wantErr {
				t.Errorf("Worker.consume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Worker.consume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorker_Produce(t *testing.T) {
	// SETUP
	var file, _ = os.Open("../../example/mockedQuoteData")
	var w = New(mockConfig())
	var i = 0
	var wg = sync.WaitGroup{}
	// END SETUP
	wg.Add(1)

	go func(ch <-chan *instruments.Quote) {
		for {
			if elem, ok := <-ch; ok {
				fmt.Println(elem)
				i++
			} else {
				break
			}
		}
		wg.Done()
	}(w.DataChan)

	go w.Produce(file, true)

	wg.Wait()
	if i != 1 {
		t.Error()
	}
}
