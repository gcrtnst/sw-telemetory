package main

import (
	"io"
)

type CloseGroup struct {
	req Wall
}

func (g *CloseGroup) Add(c io.Closer) *CloseMember {
	return newCloseMember(c, &g.req)
}

func (g *CloseGroup) CloseAll() {
	g.req.Break()
}

type CloseMember struct {
	c    io.Closer
	greq *Wall
	mreq *Wall
	done chan struct{}
	err  error
}

func newCloseMember(c io.Closer, greq *Wall) *CloseMember {
	m := &CloseMember{
		c:    c,
		greq: greq,
		mreq: &Wall{},
		done: make(chan struct{}),
		err:  nil,
	}
	m.init()
	return m
}

func (m *CloseMember) Close() error {
	m.mreq.Break()
	<-m.done
	return m.err
}

func (m *CloseMember) CloseCatch(err *error) {
	e := m.Close()
	if *err == nil {
		*err = e
	}
}

func (m *CloseMember) init() {
	select {
	case <-m.greq.Wait():
	case <-m.mreq.Wait():
	default:
		go m.worker()
		return
	}
	m.closeForce()
}

func (m *CloseMember) worker() {
	select {
	case <-m.greq.Wait():
	case <-m.mreq.Wait():
	}
	m.closeForce()
}

func (m *CloseMember) closeForce() {
	m.err = m.c.Close()
	close(m.done)
}
