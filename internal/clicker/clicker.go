// /Users/habin/workspace/automouse/internal/clicker/clicker.go
package clicker

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/go-vgo/robotgo"
)

type clickFn func(x, y int)

func defaultClick(x, y int) {
	robotgo.Move(x, y)
	robotgo.Click("left")
}

type Clicker struct {
	X        int
	Y        int
	Interval time.Duration

	running   atomic.Bool
	clickFunc clickFn
}

func New(x, y int, interval time.Duration) *Clicker {
	if interval <= 0 {
		interval = 1000 * time.Millisecond
	}
	return &Clicker{
		X:         x,
		Y:         y,
		Interval:  interval,
		clickFunc: defaultClick,
	}
}

func (c *Clicker) IsRunning() bool {
	return c.running.Load()
}

func (c *Clicker) Run(ctx context.Context, toggleCh <-chan struct{}) {
	ticker := time.NewTicker(c.Interval)
	defer ticker.Stop()

	fmt.Printf("Auto-clicker ready at (%d, %d) with %v interval\n", c.X, c.Y, c.Interval)
	fmt.Println("Press Enter to start/stop clicking. Press Ctrl+C to exit.")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("\nShutting down...")
			return

		case <-toggleCh:
			wasRunning := c.running.Load()
			c.running.Store(!wasRunning)
			if c.running.Load() {
				fmt.Println(">>> Clicking STARTED")
			} else {
				fmt.Println(">>> Clicking STOPPED")
			}

		case <-ticker.C:
			if c.running.Load() {
				c.clickFunc(c.X, c.Y)
			}
		}
	}
}
