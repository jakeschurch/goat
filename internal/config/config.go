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

package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

func ReadConfig(filename string) Config {
	var conf Config
	var file, _ = ioutil.ReadFile(filename)
	_ = json.Unmarshal(file, &conf)
	return conf
}

type Config struct {
	File struct {
		Glob          string `json:"glob,omitempty"`
		Headers       bool   `json:"headers,omitempty"`
		Delim         string `json:"delim,omitempty"`
		ExampleDate   string `json:"exampleDate,omitempty"`
		TimestampUnit string `json:"timestampUnit,omitempty"`

		Columns struct {
			Ticker    uint8 `json:"ticker,omitempty"`
			Timestamp uint8 `json:"timestamp,omitempty"`
			Bid       uint8 `json:"bid,omitempty"`
			BidSize   uint8 `json:"bidSize,omitempty"`
			Ask       uint8 `json:"ask,omitempty"`
			AskSize   uint8 `json:"askSize,omitempty"`
		} `json:"columns,omitempty"`
	} `json:"file,omitempty"`

	Backtest struct {
		StartCashAmt     float64  `json:"startCashAmt,omitempty"`
		IgnoreSecurities []string `json:"ignoreSecurities,omitempty"`
		Slippage         float64  `json:"slippage,omitempty"`
		Commission       float64  `json:"commission,omitempty"`
	} `json:"backtest,omitempty"`

	Simulation struct {
		StartDate    string        `json:"startDate,omitempty"`
		EndDate      string        `json:"endDate,omitempty"`
		BarRate      time.Duration `json:"barRate,omitempty"`
		OutputFormat string        `json:"outFmt,omitempty"`
		//  IngestRate measures how many bars to skip
		// IngestRate BarDuration `json:"ingestRate"`
	} `json:"simulation,omitempty"`

	Benchmark struct {
		Use    bool `json:"use,omitempty"`
		Update bool `json:"update,omitempty"`
	} `json:"benchmark,omitempty"`
}

func (c Config) FileInfo() (fname string, date time.Time, err error) {
	var fileGlob []string

	// read file glob and get corresponding files.
	if fileGlob, err = filepath.Glob(c.File.Glob); err != nil || len(fileGlob) == 0 {
		return fname, date, err
	}
	// get file name
	fname, _ = filepath.Abs(fileGlob[0])

	// parse date from file string
	fdate := fname[strings.LastIndex(fname, "_")+1:]
	date, err = time.Parse(c.File.ExampleDate, fdate)

	return fname, date, err
}
