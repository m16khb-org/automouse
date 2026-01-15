package watcher

import (
	"context"
	"time"

	"github.com/go-vgo/robotgo"
)

type Position struct {
	X int
	Y int
}

type Watcher struct {
	interval time.Duration
}

func New(interval time.Duration) *Watcher {
	if interval <= 0 {
		interval = 100 * time.Millisecond
	}
	return &Watcher{interval: interval}
}

func (w *Watcher) Start(ctx context.Context) <-chan Position {
	positions := make(chan Position)

	go func() {
		defer close(positions)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				x, y := robotgo.Location()
				select {
				case positions <- Position{X: x, Y: y}:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return positions
}
