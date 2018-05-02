package goat

import (
	"sync"

	"github.com/jakeschurch/instruments"

	"github.com/jakeschurch/collections"
)

type Portfolio struct {
	*collections.Portfolio
	cash instruments.Amount
	sync.RWMutex
}
