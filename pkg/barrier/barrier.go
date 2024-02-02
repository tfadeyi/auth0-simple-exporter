package barrier

import (
	"context"
	"github.com/juju/errors"
	"sync"
	"time"
)

type (
	Barrier struct {
		Capacity     int
		Value        int
		RefillPeriod time.Duration
		mutex        sync.RWMutex
	}
)

func New(cap int, rate time.Duration) *Barrier {
	return &Barrier{
		Capacity:     cap,
		Value:        cap,
		RefillPeriod: rate,
	}
}

func (b *Barrier) RecordBadEvent() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.Value <= 0 {
		return
	}
	b.Value = b.Value - 1
}

func (b *Barrier) IsBadState() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.Value == 0
}

func (b *Barrier) refill() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.Value == b.Capacity {
		return
	}
	b.Value = b.Value + 1
}

func (b *Barrier) Start(ctx context.Context) error {
	if b.Value < 0 {
		return errors.New("The instance cannot be initialised with a negative number")
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.NewTicker(b.RefillPeriod).C:
			b.refill()
		}
	}
}
