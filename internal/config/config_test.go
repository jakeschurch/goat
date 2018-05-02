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
