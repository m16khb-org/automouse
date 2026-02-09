package clicker

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

type mockClickFunc struct {
	callCount atomic.Int32
}

func (m *mockClickFunc) click(x, y int) {
	m.callCount.Add(1)
}

func TestClicker_Run_ClicksAtInterval(t *testing.T) {
	mock := &mockClickFunc{}
	c := &Clicker{
		X:             100,
		Y:             200,
		Interval:      50 * time.Millisecond,
		UseCurrentPos: false,
		clickFunc:     mock.click,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()

	toggleCh := make(chan struct{}, 1)
	toggleCh <- struct{}{} // Start immediately

	c.Run(ctx, toggleCh)

	clicks := int(mock.callCount.Load())
	if clicks < 2 || clicks > 5 {
		t.Errorf("expected 2-5 clicks, got %d", clicks)
	}
}

func TestClicker_Run_ClicksAtCurrentPos(t *testing.T) {
	mock := &mockClickFunc{}
	posX, posY := 500, 600
	c := &Clicker{
		X:             0,
		Y:             0,
		Interval:      50 * time.Millisecond,
		UseCurrentPos: true,
		clickFunc:     mock.click,
		getPosFunc:    func() (int, int) { return posX, posY },
	}

	ctx, cancel := context.WithTimeout(context.Background(), 130*time.Millisecond)
	defer cancel()

	toggleCh := make(chan struct{}, 1)
	toggleCh <- struct{}{} // Start immediately

	c.Run(ctx, toggleCh)

	clicks := int(mock.callCount.Load())
	if clicks < 1 || clicks > 4 {
		t.Errorf("expected 1-4 clicks for current pos mode, got %d", clicks)
	}
}

func TestClicker_IsRunning(t *testing.T) {
	c := New(0, 0, 100*time.Millisecond, false)

	if c.IsRunning() {
		t.Error("should not be running initially")
	}

	c.running.Store(true)
	if !c.IsRunning() {
		t.Error("should be running after setting")
	}
}

func TestNew_DefaultInterval(t *testing.T) {
	c := New(100, 200, 0, false)
	if c.Interval != 1000*time.Millisecond {
		t.Errorf("expected default 1000ms, got %v", c.Interval)
	}
}

func TestNew_UseCurrentPos(t *testing.T) {
	c := New(0, 0, 100*time.Millisecond, true)
	if !c.UseCurrentPos {
		t.Error("expected UseCurrentPos to be true")
	}

	c2 := New(100, 200, 100*time.Millisecond, false)
	if c2.UseCurrentPos {
		t.Error("expected UseCurrentPos to be false")
	}
}
