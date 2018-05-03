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

package goat

import (
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"github.com/jakeschurch/instruments"

	"github.com/jakeschurch/goat/internal/config"
)

func TestReadConfig(t *testing.T) {
	filename, err := filepath.Abs("../../example/config.json")
	if err != nil {
		panic("could not read json")
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want config.Config
	}{
		{"base case", args{filename}, ReadConfig(filename)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadConfig(tt.args.filename); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

type Algorithm_Example struct{}

func (algo Algorithm_Example) Buy(quote instruments.Quote) (*instruments.Order, bool) {
	newOrder := quote.FillOrder(quote.Ask.Price, instruments.NewVolume(20.00), true, instruments.Market)
	return newOrder, true
}

func (algo Algorithm_Example) Sell(quote instruments.Quote, holding *instruments.Holding) (*instruments.Order, bool) {
	if quote.Name == holding.Name && quote.Bid.Price > holding.Buy.Price {
		newOrder := quote.FillOrder(quote.Bid.Price, holding.Volume, false, instruments.Market)
		return newOrder, true
	}
	return nil, false

}

func TestNewSim(t *testing.T) {
	// SETUP
	filename, err := filepath.Abs("../../example/config.json")
	if err != nil {
		panic("could not read json")
	}

	wanted := &Simulation{
		conf:   ReadConfig(filename),
		algos:  []Algorithm{Algorithm_Example{}},
		ignore: sync.Map{},
	}
	for _, name := range wanted.conf.Backtest.IgnoreSecurities {
		wanted.ignore.Store(name, struct{}{})
	}
	// END SETUP

	type args struct {
		c     config.Config
		algos []Algorithm
	}
	tests := []struct {
		name string
		args args
		want *Simulation
	}{
		{"base case", args{ReadConfig(filename), []Algorithm{Algorithm_Example{}}}, wanted},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSim(tt.args.c, tt.args.algos...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSim() = %v, want %v", got, tt.want)
			}
		})
	}
}
