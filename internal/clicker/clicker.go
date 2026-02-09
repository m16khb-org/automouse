package clicker

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/go-vgo/robotgo"
)

type clickFn func(x, y int)
type getPosFn func() (int, int)

func defaultClick(x, y int) {
	robotgo.Move(x, y)
	robotgo.Click("left")
}

func defaultGetPos() (int, int) {
	return robotgo.Location()
}

type Clicker struct {
	X             int
	Y             int
	Interval      time.Duration
	UseCurrentPos bool

	running   atomic.Bool
	clickFunc clickFn
	getPosFunc getPosFn
}

func New(x, y int, interval time.Duration, useCurrentPos bool) *Clicker {
	if interval <= 0 {
		interval = 1000 * time.Millisecond
	}
	return &Clicker{
		X:             x,
		Y:             y,
		Interval:      interval,
		UseCurrentPos: useCurrentPos,
		clickFunc:     defaultClick,
		getPosFunc:    defaultGetPos,
	}
}

func (c *Clicker) IsRunning() bool {
	return c.running.Load()
}

func (c *Clicker) Run(ctx context.Context, toggleCh <-chan struct{}) {
	ticker := time.NewTicker(c.Interval)
	defer ticker.Stop()

	if c.UseCurrentPos {
		fmt.Printf("Auto-clicker ready at CURRENT MOUSE POSITION with %v interval\n", c.Interval)
	} else {
		fmt.Printf("Auto-clicker ready at (%d, %d) with %v interval\n", c.X, c.Y, c.Interval)
	}
	fmt.Println("Press Enter to start/stop clicking. Press ESC or Ctrl+C to exit.")

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
				if c.UseCurrentPos {
					x, y := c.getPosFunc()
					c.clickFunc(x, y)
				} else {
					c.clickFunc(c.X, c.Y)
				}
			}
		}
	}
}
