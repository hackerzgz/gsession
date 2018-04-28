package ticker

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	// StampTail keep last 6 number in timestamp at most
	StampTail = 1000000
)

// Method of refresh timestamp each different timing
type Method int32

const (
	// EachSecTicker refresh timestamp each 1 seconds
	EachSecTicker Method = 1
	// TenSecTicker refresh timestamp each 10 seconds
	TenSecTicker Method = 10
	// HunSecTicker refresh timestamp each 100 seconds
	HunSecTicker = 100
)

// TimeTicker
type TimeTicker struct {
	start     time.Time
	timestamp int32
	bits      int32

	m          *sync.RWMutex
	ticker     *time.Ticker
	cancelFunc context.CancelFunc
}

// New different refresh rate TimeTicker
func New(meth Method) *TimeTicker {
	now := time.Now()
	bs := StampTail / meth

	ctx, cancel := context.WithCancel(context.Background())
	t := &TimeTicker{
		start:      now,
		bits:       int32(bs),
		timestamp:  int32(Method(now.Unix()%int64(bs)) / meth),
		m:          &sync.RWMutex{},
		ticker:     time.NewTicker(time.Duration(meth) * time.Second),
		cancelFunc: cancel,
	}

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("TimeTicker canceled, start at: ", t.start)
			return
		case <-t.ticker.C:
			t.refresh()
		}
	}()

	return t
}

// refresh
func (t *TimeTicker) refresh() {
	t.m.Lock()
	t.timestamp++
	t.timestamp %= t.bits
	t.m.Unlock()
}

// Stop clear TimeTicker and cancel refresh goroutine
func (t *TimeTicker) Stop() {
	t.cancelFunc()
	return
}

// Get return current timestamp from TimeTicker
func (t *TimeTicker) Get() int32 {
	t.m.RLock()
	defer t.m.RUnlock()
	return t.timestamp
}
