package config

import (
	"log"
	"path/filepath"
	"strings"
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

func (c Config) FileInfo() (string, time.Time) {
	fileGlob, err := filepath.Glob(c.File.Glob)
	if err != nil || len(fileGlob) == 0 {
		log.Println(err)
		return "", time.Time{}
	}
	filename := fileGlob[0]
	lastUnderscore := strings.LastIndex(filename, "_")
	fileDate := filename[lastUnderscore+1:]

	lastDate, dateErr := time.Parse(c.File.ExampleDate, fileDate)
	if dateErr != nil {
		log.Fatal("Date cannot be parsed")
	}
	return filename, lastDate
}
