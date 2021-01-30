package main

import (
	"errors"
	"testing"
)

func TestCloseMemberInit(t *testing.T) {
	for i := 0; i < 1000; i++ {
		mcnt := 0
		c := &mockCloser{func() error {
			mcnt++
			return nil
		}}
		m := &CloseMember{
			c:    c,
			greq: &Wall{},
			mreq: &Wall{},
			done: make(chan struct{}),
			err:  nil,
		}
		m.greq.Break()
		m.init()
		if cnt := mcnt; cnt != 1 {
			t.Errorf("expected cnt==1, got cnt==%d", cnt)
			break
		}
	}

	mcnt := 0
	c := &mockCloser{func() error {
		mcnt++
		return nil
	}}
	m := &CloseMember{
		c:    c,
		greq: &Wall{},
		mreq: &Wall{},
		done: make(chan struct{}),
		err:  nil,
	}
	m.init()
	m.greq.Break()
	<-m.done
	if cnt := mcnt; cnt != 1 {
		t.Errorf("expected cnt==1, got cnt==%d", cnt)
	}
}

func TestCloseGroupNested(t *testing.T) {
	err := errors.New("TEST")
	for i := 0; i < 1000; i++ {
		mcnt := 0
		g := &CloseGroup{}
		c10 := &mockCloser{func() error {
			mcnt++
			return err
		}}
		c20 := &mockCloser{func() error {
			g.CloseAll()
			return nil
		}}
		m11 := g.Add(c10)
		m12 := g.Add(m11)
		m21 := g.Add(c20)
		m22 := g.Add(m21)

		_ = m22.Close()
		<-m12.done
		if cnt := mcnt; cnt != 1 {
			t.Errorf("expected cnt==1, got cnt==%d", cnt)
		}
		if m12.err != err {
			t.Errorf("got wrong error: %#v", m12.err)
		}
	}
}

type mockCloser struct {
	f func() error
}

func (c *mockCloser) Close() error {
	return c.f()
}
