package main

import (
	"context"
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

func TestCloseOnCancelCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mock := mockCloser{ch: make(chan struct{}), err: errors.New("")}

	c := CloseOnCancel(ctx, mock)
	cancel()
	<-mock.ch
	err := c.Close()
	if err != mock.err {
		t.Error()
	}
}

func TestCloseOnCancelClose(t *testing.T) {
	ctx := context.Background()
	mock := mockCloser{ch: make(chan struct{}), err: errors.New("")}

	c := CloseOnCancel(ctx, mock)
	err := c.Close()
	select {
	case <-mock.ch:
	default:
		t.Error()
	}
	if err != mock.err {
		t.Error()
	}
}

func TestCloseCatchAssignError(t *testing.T) {
	mock := mockCloser{ch: make(chan struct{}), err: errors.New("")}

	var err error
	CloseCatch(mock, &err)
	<-mock.ch
	if err != mock.err {
		t.Error()
	}
}

func TestCloseCatchIgnoreError(t *testing.T) {
	err_back := errors.New("back")
	err_mock := errors.New("mock")
	mock := mockCloser{ch: make(chan struct{}), err: err_mock}

	err := err_back
	CloseCatch(mock, &err)
	<-mock.ch
	if err != err_back {
		t.Error()
	}
}

func TestCatchContextError(t *testing.T) {
	open_ctx := context.Background()
	done_ctx, cancel := context.WithCancel(context.Background())
	cancel()
	test_err := errors.New("")

	cases := []struct {
		inCtx   context.Context
		inErr   error
		wantErr error
	}{
		{
			inCtx:   open_ctx,
			inErr:   nil,
			wantErr: nil,
		},
		{
			inCtx:   open_ctx,
			inErr:   test_err,
			wantErr: test_err,
		},
		{
			inCtx:   done_ctx,
			inErr:   nil,
			wantErr: context.Canceled,
		},
		{
			inCtx:   done_ctx,
			inErr:   test_err,
			wantErr: context.Canceled,
		},
	}

	for i, c := range cases {
		gotErr := c.inErr
		CatchContextError(c.inCtx, &gotErr)

		if gotErr != c.wantErr {
			t.Errorf(`case %d: err: expected "%s", got "%s"`, i, c.wantErr, gotErr)
		}
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
