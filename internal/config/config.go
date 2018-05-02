package config

import (
	"time"
)

type Config struct {
	File struct {
		Glob          string `json:"glob"`
		Headers       bool   `json:"headers"`
		Delim         string `json:"delim"`
		ExampleDate   string `json:"exampleDate"`
		TimestampUnit string `json:"timestampUnit"`

		Columns struct {
			Ticker    uint8 `json:"ticker"`
			Timestamp uint8 `json:"timestamp"`
			Bid       uint8 `json:"bid"`
			BidSize   uint8 `json:"bidSize"`
			Ask       uint8 `json:"ask"`
			AskSize   uint8 `json:"askSize"`
		} `json:"columns"`
	} `json:"file"`

	Backtest struct {
		StartCashAmt     float64  `json:"startCashAmt"`
		IgnoreSecurities []string `json:"ignoreSecurities"`
		Slippage         float64  `json:"slippage"`
		Commission       float64  `json:"commission"`
	} `json:"backtest"`

	Simulation struct {
		StartDate    string        `json:"startDate"`
		EndDate      string        `json:"endDate"`
		BarRate      time.Duration `json:"barRate"`
		OutputFormat string        `json:"outFmt"`
		//  IngestRate measures how many bars to skip
		// IngestRate BarDuration `json:"ingestRate"`
	} `json:"simulation"`

	Benchmark struct {
		Use    bool `json:"use"`
		Update bool `json:"update"`
	} `json:"benchmark"`
}
