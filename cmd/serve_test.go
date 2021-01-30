package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"testing"
)

func TestServerServe(t *testing.T) {
	cfg := NewServerConfig()
	cfg.Log = log.New(ioutil.Discard, "", log.LstdFlags)

	conn1 := &mockConn{
		rd: strings.NewReader(string(chunkPrefix) + "n\n" + string(chunkSuffix)),
	}
	conn2 := &mockConn{
		rd: strings.NewReader(string(chunkPrefix) + "n\n" + string(chunkSuffix)),
	}
	lis := &mockListener{
		conns: []net.Conn{conn1, conn2},
		err:   errors.New("TEST"),
	}
	srv := NewServer(cfg)
	err := srv.Serve(lis)
	if err != lis.err {
		t.Errorf("expected %#v, got %#v", lis.err, err)
	}
	if !lis.closed {
		t.Errorf("listener not closed")
	}
	if !conn1.closed {
		t.Errorf("conn1 not closed")
	}
	if !conn2.closed {
		t.Errorf("conn2 not closed")
	}

	conn1 = &mockConn{
		rd: strings.NewReader(string(chunkPrefix) + "\x00\n" + string(chunkSuffix)),
	}
	conn2 = &mockConn{
		rd: strings.NewReader(string(chunkPrefix) + "\x00\n" + string(chunkSuffix)),
	}
	lis = &mockListener{
		conns: []net.Conn{conn1, conn2},
		err:   errors.New("TEST"),
	}
	srv = NewServer(cfg)
	err = srv.Serve(lis)
	if err == lis.err {
		t.Errorf("unexpected error %#v", err)
	}
	if !lis.closed {
		t.Errorf("listener not closed")
	}
	if !conn1.closed {
		t.Errorf("conn1 not closed")
	}
	if conn2.closed {
		t.Errorf("conn2 closed")
	}
}
