package main

import (
	"errors"
	"testing"
)

func TestCloseGroupCloseAllAfterAdd(t *testing.T) {
	mock1 := mockCloser{ch: make(chan struct{})}
	mock2 := mockCloser{ch: make(chan struct{})}

	cg := &CloseGroup{}
	cg.Add(mock1)
	cg.Add(mock2)
	cg.CloseAll()
	<-mock1.ch
	<-mock2.ch
}

func TestCloseGroupCloseAllBeforeAdd(t *testing.T) {
	mock1 := mockCloser{ch: make(chan struct{})}
	mock2 := mockCloser{ch: make(chan struct{})}

	cg := &CloseGroup{}
	cg.CloseAll()
	cg.Add(mock1)
	cg.Add(mock2)
	<-mock1.ch
	<-mock2.ch
}

func TestCloseGroupCloseAllMultiple(t *testing.T) {
	mock1 := mockCloser{ch: make(chan struct{})}
	mock2 := mockCloser{ch: make(chan struct{})}

	cg := &CloseGroup{}
	cg.Add(mock1)
	cg.Add(mock2)
	cg.CloseAll()
	<-mock1.ch
	<-mock2.ch

	cg.CloseAll()
}

func TestCloseMemberCloseBeforeCloseAll(t *testing.T) {
	mock := mockCloser{ch: make(chan struct{}), err: errors.New("")}

	cg := &CloseGroup{}
	err := cg.Add(mock).Close()
	select {
	case <-mock.ch:
	default:
		t.Error()
	}
	if err != mock.err {
		t.Error()
	}

	cg.CloseAll()
}

func TestCloseMemberCloseAfterCloseAll(t *testing.T) {
	mock := mockCloser{ch: make(chan struct{}), err: errors.New("")}

	cg := &CloseGroup{}
	cm := cg.Add(mock)
	cg.CloseAll()
	<-mock.ch
	err := cm.Close()
	if err != mock.err {
		t.Error()
	}
}

func TestCloseMemberCloseMultiple(t *testing.T) {
	mock := mockCloser{ch: make(chan struct{}), err: errors.New("")}

	cg := &CloseGroup{}
	cm := cg.Add(mock)
	err1 := cm.Close()
	select {
	case <-mock.ch:
	default:
		t.Error()
	}
	if err1 != mock.err {
		t.Error()
	}

	err2 := cm.Close()
	if err2 != mock.err {
		t.Error()
	}
}

func TestCloseMemberCloseCatchAssignError(t *testing.T) {
	mock := mockCloser{ch: make(chan struct{}), err: errors.New("")}

	var err error
	cg := &CloseGroup{}
	cg.Add(mock).CloseCatch(&err)
	select {
	case <-mock.ch:
	default:
		t.Error()
	}
	if err != mock.err {
		t.Error()
	}
}

func TestCloseMemberCloseCatchIgnoreError(t *testing.T) {
	err_back := errors.New("back")
	err_mock := errors.New("mock")
	mock := mockCloser{ch: make(chan struct{}), err: err_mock}

	err := err_back
	cg := &CloseGroup{}
	cg.Add(mock).CloseCatch(&err)
	select {
	case <-mock.ch:
	default:
		t.Error()
	}
	if err != err_back {
		t.Error()
	}
}

type mockCloser struct {
	ch  chan struct{}
	err error
}

func (m mockCloser) Close() error {
	close(m.ch)
	return m.err
}
