package main

import "sync"

type Wall struct {
	mu sync.Mutex
	ch chan struct{}
}

func (w *Wall) Wait() <-chan struct{} {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.ch == nil {
		w.ch = make(chan struct{})
	}
	return w.ch
}

func (w *Wall) Break() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.ch == nil {
		w.ch = make(chan struct{})
	}
	select {
	case <-w.ch:
	default:
		close(w.ch)
	}
}
