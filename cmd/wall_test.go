package main

import (
	"testing"
)

func TestWallWaitNoBreak(t *testing.T) {
	w := &Wall{}
	select {
	case <-w.Wait():
		t.Fail()
	default:
	}
}

func TestWallWaitBeforeBreak(t *testing.T) {
	w := &Wall{}
	wait := w.Wait()
	w.Break()
	select {
	case <-wait:
	default:
		t.Fail()
	}
}

func TestWallWaitAfterBreak(t *testing.T) {
	w := &Wall{}
	w.Break()
	select {
	case <-w.Wait():
	default:
		t.Fail()
	}
}

func TestWallWaitNil(t *testing.T) {
	w := &Wall{ch: nil}
	wait := w.Wait()

	if wait == nil {
		t.Error()
	}
	select {
	case <-wait:
		t.Error()
	default:
	}
}

func TestWallWaitOpen(t *testing.T) {
	ch := make(chan struct{})
	w := &Wall{ch: ch}
	wait := w.Wait()

	if w.ch != ch {
		t.Error()
	}
	if wait != w.ch {
		t.Error()
	}
	if wait == nil {
		t.Error()
	}
	select {
	case <-wait:
		t.Error()
	default:
	}
}

func TestWallWaitClose(t *testing.T) {
	ch := make(chan struct{})
	close(ch)
	w := &Wall{ch: ch}
	wait := w.Wait()

	if w.ch != ch {
		t.Error()
	}
	if wait != w.ch {
		t.Error()
	}
	if wait == nil {
		t.Error()
	}
	select {
	case <-wait:
	default:
		t.Error()
	}
}

func TestWallBreakNil(t *testing.T) {
	w := &Wall{ch: nil}
	w.Break()

	if w.ch == nil {
		t.Error()
	}
	select {
	case <-w.ch:
	default:
		t.Error()
	}
}

func TestWallBreakOpen(t *testing.T) {
	ch := make(chan struct{})
	w := &Wall{ch: ch}
	w.Break()

	if w.ch != ch {
		t.Error()
	}
	select {
	case <-w.ch:
	default:
		t.Error()
	}
}

func TestWallBreakClose(t *testing.T) {
	ch := make(chan struct{})
	close(ch)
	w := &Wall{ch: ch}
	w.Break()

	if w.ch != ch {
		t.Error()
	}
	select {
	case <-w.ch:
	default:
		t.Error()
	}
}
