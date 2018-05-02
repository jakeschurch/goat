package goat

import (
	"encoding/json"
	"os"

	"github.com/jakeschurch/goat/internal/config"
	"github.com/jakeschurch/instruments"
)

// ReadConfig
func ReadConfig(file *os.File) config.Config {
	var conf = &config.Config{}
	decoder := json.NewDecoder(file)
	decoder.Decode(conf)
	return *conf
}

// Algorithm is an interface that needs to be implemented in the pipeline by a user to fill orders based on the conditions that they specify.
type Algorithm interface {
	Buy(instruments.Quote) (*instruments.Order, bool)
	Sell(instruments.Quote, ...*instruments.Holding) (*instruments.Order, error)
}
type Simulation struct {
	conf  config.Config
	algos []Algorithm
}

func NewSim(c config.Config, algos ...Algorithm) *Simulation {
	return &Simulation{
		conf:  c,
		algos: algos,
	}
}

// checkBuys from quote information.
// Buy Orders handled by Simulation;
// sells by Portfolios.
func (sim Simulation) checkBuys(quote instruments.Quote) *instruments.Order {
	for _, algo := range sim.algos {
		if order, ok := algo.Buy(quote); ok {
			return order
		}
	}
	return nil
}

func (sim Simulation) checkSells(quote instruments.Quote) *instruments.Order {
	for _, algo := range sim.algos {
		if order, ok := algo.Buy(quote); ok {
			return order
		}
	}
	return nil
}

func (sim *Simulation) Run() {

}
