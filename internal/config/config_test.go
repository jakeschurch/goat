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
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestConfig_FileInfo(t *testing.T) {
	var fname, _ = filepath.Abs("../../example/config.json")
	var wantedName, _ = filepath.Abs("../../example/testQuotes_20170814")
	var wantedDate, _ = time.Parse("20060102", "20170814")
	var conf = ReadConfig(fname)

	type fields struct {
		conf Config
	}
	tests := []struct {
		name      string
		fields    fields
		wantFname string
		wantDate  time.Time
		wantErr   bool
	}{
		{"base case", fields{conf}, wantedName, wantedDate, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields.conf
			gotFname, gotDate, err := c.FileInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.FileInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFname != tt.wantFname {
				t.Errorf("Config.FileInfo() gotFname = %v, want %v", gotFname, tt.wantFname)
			}
			if !reflect.DeepEqual(gotDate, tt.wantDate) {
				t.Errorf("Config.FileInfo() gotDate = %v, want %v", gotDate, tt.wantDate)
			}
		})
	}
}
