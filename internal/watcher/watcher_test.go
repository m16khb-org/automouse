package watcher

import (
	"context"
	"testing"
	"time"
)

func TestWatcher_Start_SendsPositions(t *testing.T) {
	w := New(100 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	positions := w.Start(ctx)

	count := 0
	for range positions {
		count++
	}

	if count < 2 {
		t.Errorf("expected at least 2 positions, got %d", count)
	}
}

func TestWatcher_Start_StopsOnContextCancel(t *testing.T) {
	w := New(50 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())

	positions := w.Start(ctx)
	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case _, ok := <-positions:
		if ok {
			for range positions {
			}
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("channel did not close after context cancel")
	}
}

func TestNew_DefaultInterval(t *testing.T) {
	w := New(0)
	if w.interval != 100*time.Millisecond {
		t.Errorf("expected default 100ms, got %v", w.interval)
	}
}
