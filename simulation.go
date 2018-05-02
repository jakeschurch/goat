package goat

import (
	"io"
	"os"
	"sync"

	"github.com/jakeschurch/goat/internal/config"
	"github.com/jakeschurch/goat/internal/worker"
	"github.com/jakeschurch/instruments"
)

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
	return sim
}

// checkBuys from quote information.
// Buy Orders handled by Simulation;
// sells by Portfolios.
func (sim *Simulation) checkBuys(quote instruments.Quote) *instruments.Order {
	if _, ok := sim.ignore.Load(quote.Name); ok {
		return nil
	}
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

	quoteChan := make(chan *instruments.Quote)
	done := make(chan struct{})
	go func(inChan <-chan *instruments.Quote) {
		var quote *instruments.Quote
		var ok bool
	loop:
		for {
			if quote, ok = <-quoteChan; !ok {
				break loop
			}
			if quote != nil {
				sim.process(quote)
			}
			continue
		}
		close(done)
	}(quoteChan)

	if file, err = os.Open(fname); err != nil {
		return err
	}
	worker.Run(quoteChan, file)
	<-done
	return nil
}

func (sim *Simulation) process(quote *instruments.Quote) {

}
