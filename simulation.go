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
	"io"
	"os"
	"sync"

	"github.com/jakeschurch/goat/internal/output"

	"github.com/jakeschurch/goat/internal/config"
	"github.com/jakeschurch/goat/internal/worker"
	"github.com/jakeschurch/instruments"
)

var (
	orderManager   *OrderManager
	Port           *Portfolio
	performanceLog *output.PerformanceLog
	// benchmark      *Benchmark
)

func init() {
	orderManager = NewOrderManager()
	Port = NewPortfolio(instruments.Amount(0))
	performanceLog = output.NewPerformanceLog()
}

// ReadConfig
func ReadConfig(filename string) config.Config {
	return config.ReadConfig(filename)
}

// Algorithm is an interface that needs to be implemented in the pipeline by a user to fill orders based on the conditions that they specify.
type Algorithm interface {
	Buy(instruments.Quote) (*instruments.Order, bool)
	Sell(instruments.Quote, *instruments.Holding) (*instruments.Order, bool)
}
type Simulation struct {
	conf   config.Config
	algos  []Algorithm
	ignore sync.Map
}

func NewSim(c config.Config, algos ...Algorithm) *Simulation {
	var sim = &Simulation{
		conf:   c,
		algos:  algos,
		ignore: sync.Map{},
	}
	for _, name := range c.Backtest.IgnoreSecurities {
		sim.ignore.Store(name, struct{}{})
	}
	Port.cash += instruments.NewAmount(instruments.NewPrice(1.00), instruments.NewVolume(c.Backtest.StartCashAmt))
	return sim
}

// checkBuys from quote information.
// Buy Orders handled by Simulation;
// sells by Portfolios.
func (sim *Simulation) checkBuys(quote instruments.Quote) *instruments.Order {
	for _, algo := range sim.algos {
		if order, ok := algo.Buy(quote); ok {
			return order
		}
	}
	return nil
}

func (sim *Simulation) Run() error {
	var file io.ReadSeeker

	// Setup file reader
	var fname, date, err = sim.conf.FileInfo()
	if err != nil {
		return err
	}
	// Setup Worker & WorkerConfig
	wc := worker.Config{
		Name: sim.conf.File.Columns.Ticker,
		Bid:  sim.conf.File.Columns.Bid, BidSz: sim.conf.File.Columns.BidSize,
		Ask: sim.conf.File.Columns.Ask, AskSz: sim.conf.File.Columns.AskSize,
		Timestamp: sim.conf.File.Columns.Timestamp, Date: date,
		Timeunit: sim.conf.File.TimestampUnit,
	}
	worker := worker.New(wc)

	if file, err = os.Open(fname); err != nil {
		return err
	}

	var wg = sync.WaitGroup{}

	wg.Add(1)

	var done = make(chan struct{})

	go func(ch <-chan *instruments.Quote) {
		for {
			if elem, ok := <-ch; ok {
				if _, ok := sim.ignore.Load(elem.Name); !ok {
					sim.process(elem)
				}
			} else {
				break
			}
		}
		wg.Done()
	}(worker.DataChan)

	go worker.Produce(file, true)

	wg.Wait()
	<-done

	Port.CloseAll()
	performanceLog.OutputResults(output.CSV, "/home/jake/Desktop/simResults.csv")
	return nil
}

func (sim *Simulation) process(quote *instruments.Quote) {
	// Check if we can buy new holding
	if newBuy := sim.checkBuys(*quote); newBuy != nil {
		orderManager.Add(newBuy)
	}
	Port.Update(*quote, sim.algos...)
}
