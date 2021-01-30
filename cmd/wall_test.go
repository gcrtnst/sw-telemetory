package main

import (
	"testing"
)

func TestWall(t *testing.T) {
	w := &Wall{}
	select {
	case <-w.Wait():
		t.Error("wall broken")
	default:
	}

	w = &Wall{}
	wait := w.Wait()
	w.Break()
	select {
	case <-wait:
	default:
		t.Error("wall not broken")
	}

	w = &Wall{}
	w.Break()
	select {
	case <-w.Wait():
	default:
		t.Error("wall not broken")
	}
}

func TestWallWait(t *testing.T) {
	closed := make(chan struct{})
	close(closed)

	tests := []chan struct{}{
		nil,
		make(chan struct{}),
		closed,
	}
	for i, tt := range tests {
		w := &Wall{ch: tt}
		ch := w.Wait()
		if ch != w.ch {
			t.Errorf("case %d: ch != w.ch", i)
		}
		if tt == nil {
			if ch == nil {
				t.Errorf("case %d: ch == nil", i)
			}
		} else {
			if ch != tt {
				t.Errorf("case %d: ch != tt", i)
			}
		}
	}
}

func TestWallBreak(t *testing.T) {
	closed := make(chan struct{})
	close(closed)

	tests := []chan struct{}{
		nil,
		make(chan struct{}),
		closed,
	}
	for i, tt := range tests {
		w := &Wall{ch: tt}
		w.Break()
		if w.ch == nil {
			t.Errorf("case %d: non-nil channel", i)
			continue
		}
		select {
		case <-w.ch:
		default:
			t.Errorf("case %d: channel not closed", i)
		}
	}
}
