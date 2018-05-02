package ordering

import (
	"time"
)

// orderTicker is an internal struct that allows
// for simulation of order/transaction datetimes.
type OrderTicker struct {
	start  time.Time
	ticker *time.Ticker
}

func NewOrderTicker() *OrderTicker {
	return &OrderTicker{
		start:  time.Now(),
		ticker: time.NewTicker(time.Millisecond),
	}
}
func (oT *OrderTicker) Duration() time.Duration {
	t := <-oT.ticker.C
	return t.Sub(oT.start)
}
