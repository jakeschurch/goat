package goat

import (
	"encoding/json"
	"os"

	"github.com/jakeschurch/goat/internal/config"
)

// ReadConfig
func ReadConfig(file *os.File) config.Config {
	var conf = &config.Config{}
	decoder := json.NewDecoder(file)
	decoder.Decode(conf)
	return *conf
}
