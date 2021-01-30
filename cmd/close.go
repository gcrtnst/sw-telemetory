package main

import (
	"io"
)

type CloseGroup struct {
	req Wall
}

func (g *CloseGroup) Add(c io.Closer) *CloseMember {
	m := &CloseMember{
		c:    c,
		greq: &g.req,
		mreq: &Wall{},
		done: make(chan struct{}),
		err:  nil,
	}
	m.init()
	return m
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

func (m *CloseMember) Close() error {
	m.mreq.Break()
	<-m.done
	return m.err
}

func (m *CloseMember) CloseCatch(err *error) {
	if *err != nil {
		m.Close()
		return
	}
	*err = m.Close()
}

func (m *CloseMember) init() {
	select {
	case <-m.greq.Wait():
	case <-m.mreq.Wait():
	default:
		go m.worker()
		return
	}
	m.closeImmediately()
}

func (m *CloseMember) worker() {
	select {
	case <-m.greq.Wait():
	case <-m.mreq.Wait():
	}
	m.closeImmediately()
}

func (m *CloseMember) closeImmediately() {
	m.err = m.c.Close()
	close(m.done)
}
