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
	"io"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/jakeschurch/goat/internal/config"

	"github.com/jakeschurch/instruments"
)

func mockWc() Config {
	conf := config.ReadConfig("../../example/config.json")
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
func TestNew(t *testing.T) {
	type args struct {
		wc Config
	}
	tests := []struct {
		name string
		args args
		want *Worker
	}{
		{"base case", args{mockWc()}, New(mockWc())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.wc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorker_Run(t *testing.T) {
	quoteChan := make(chan *instruments.Quote)
	done := make(chan struct{})

	go func() {
		for {
			if _, ok := <-quoteChan; !ok {
				break
			}
		}
		close(done)
	}()

	file, _ := os.Open("../../example/config.json")
	type args struct {
		outChan chan<- *instruments.Quote
		r       io.ReadSeeker
	}
	tests := []struct {
		name   string
		worker *Worker
		args   args
	}{
		{"base case", New(mockWc()), args{quoteChan, file}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.worker.Run(tt.args.outChan, tt.args.r)
			<-done
		})
	}
}

func TestWorker_produce(t *testing.T) {
	type args struct {
		r  io.ReadSeeker
		wg *sync.WaitGroup
	}
	tests := []struct {
		name   string
		worker *Worker
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.worker.produce(tt.args.r, tt.args.wg)
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
