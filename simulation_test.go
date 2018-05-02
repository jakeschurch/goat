package goat

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jakeschurch/goat/internal/config"
)

func TestReadConfig(t *testing.T) {
	// SETUP
	name, _ := filepath.Abs("../example/config.json")
	fmt.Println(name)
	file, _ := os.Open(name)
	fmt.Println(file)
	// END SETUP

	type args struct {
		file *os.File
	}
	tests := []struct {
		name string
		args args
		want config.Config
	}{
		{"base case", args{file}, ReadConfig(file)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadConfig(tt.args.file); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
