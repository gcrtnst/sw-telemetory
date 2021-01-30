package main

import (
	"errors"
	"io"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestReceiverRecvChunk(t *testing.T) {
	merr := errors.New("TEST")
	lis := &mockListener{
		err: merr,
	}
	r := &Receiver{
		lis:     lis,
		timeout: 0,
		cg:      CloseGroup{},
	}
	_, err := r.recvChunk()
	if err != merr {
		t.Errorf("got wrong error: %v", err)
	}

	buf := string(chunkPrefix) + "abc" + string(chunkSuffix)
	conn := &mockConn{
		rd:       strings.NewReader(buf),
		closed:   false,
		deadline: time.Time{},
	}
	lis = &mockListener{
		conns: []net.Conn{conn},
	}
	r = &Receiver{
		lis:     lis,
		timeout: 0,
		cg:      CloseGroup{},
	}
	chunk, _ := r.recvChunk()
	if string(chunk) != buf {
		t.Errorf("expected \"%s\", got \"%s\"", buf, string(chunk))
	}
	if !conn.closed {
		t.Errorf("conn not closed")
	}
	if !conn.deadline.IsZero() {
		t.Errorf("expected no deadline, got %v", conn.deadline)
	}

	conn = &mockConn{
		rd:       strings.NewReader("GET abc HTTP/1.1\r\n"),
		closed:   false,
		deadline: time.Time{},
	}
	lis = &mockListener{
		conns: []net.Conn{conn},
	}
	r = &Receiver{
		lis:     lis,
		timeout: 1,
		cg:      CloseGroup{},
	}
	r.recvChunk()
	if conn.deadline.IsZero() {
		t.Errorf("expected deadline, got no deadline")
	}
}

func TestReadUntil(t *testing.T) {
	tests := []struct {
		in    string
		delim string
		out   string
		err   error
	}{
		{"", "", "", nil},
		{"", "TEST", "", io.EOF},
		{"TEST", "", "", nil},
		{"TEST", "TEST", "TEST", nil},
		{"TESTTEST", "TEST", "TEST", nil},
		{"TEST", "TESTTEST", "TEST", io.EOF},
		{"FOO", "BAR", "FOO", io.EOF},
		{"FOOBAR", "BAR", "FOOBAR", nil},
	}
	for _, tt := range tests {
		r := strings.NewReader(tt.in)
		delim := []byte(tt.delim)
		out, err := ReadUntil(r, delim)
		if !reflect.DeepEqual(string(out), tt.out) {
			t.Errorf("input %#v, %#v: got %#v, want %#v", tt.in, tt.delim, out, tt.out)
		}
		if !errors.Is(err, tt.err) {
			t.Errorf("input %#v, %#v: got error %#v, want %#v", tt.in, tt.delim, err, tt.err)
		}
		if len(tt.in)-r.Len() != len(tt.out) {
			t.Errorf("input %#v, %#v: wrong amount consumed", tt.in, tt.delim)
		}
		if !reflect.DeepEqual(delim, []byte(tt.delim)) {
			t.Errorf("input %#v, %#v: delim modified", tt.in, tt.delim)
		}
	}
}

func TestExtractBody(t *testing.T) {
	prefix := string(chunkPrefix)
	suffix := string(chunkSuffix)

	tests := []struct {
		reqline string
		cmd     string
		err     bool
	}{
		{"", "", true},
		{prefix, "", true},
		{suffix, "", true},
		{prefix + "/ HTTP/1.1", "", true},
		{prefix[:len(prefix)-1] + suffix, "", true},
		{prefix + suffix, "", false},
		{prefix + "/" + suffix, "/", false},
		{prefix + "abc" + suffix, "abc", false},
		{prefix + "newline\nallowed" + suffix, "newline\nallowed", false},
	}
	for _, tt := range tests {
		body, err := ExtractBody([]byte(tt.reqline))
		if string(body) != tt.cmd {
			t.Errorf("input %#v: got %#v, want %#v", tt.reqline, body, tt.cmd)
		}
		if (err != nil) != tt.err {
			t.Errorf("input %#v: wrong error", tt.reqline)
		}
	}
}

type mockListener struct {
	conns  []net.Conn
	err    error
	closed bool
}

func (l *mockListener) Accept() (net.Conn, error) {
	if len(l.conns) <= 0 {
		return nil, l.err
	}
	conn := l.conns[0]
	l.conns = l.conns[1:]
	return conn, nil
}

func (l *mockListener) Close() error {
	l.closed = true
	return nil
}

func (l *mockListener) Addr() net.Addr {
	return &mockAddr{
		network: "TEST",
		address: "TEST",
	}
}

type mockConn struct {
	rd       *strings.Reader
	closed   bool
	deadline time.Time
}

func (c *mockConn) Read(b []byte) (int, error) {
	return c.rd.Read(b)
}

func (c *mockConn) Write(b []byte) (int, error) {
	panic("not implemented")
}

func (c *mockConn) Close() error {
	c.closed = true
	return nil
}

func (c *mockConn) LocalAddr() net.Addr {
	panic("not implemented")
}

func (c *mockConn) RemoteAddr() net.Addr {
	panic("not implemented")
}

func (c *mockConn) SetDeadline(t time.Time) error {
	c.deadline = t
	return nil
}

func (c *mockConn) SetReadDeadline(t time.Time) error {
	panic("not implemented")
}

func (c *mockConn) SetWriteDeadline(t time.Time) error {
	panic("not implemented")
}

type mockAddr struct {
	network string
	address string
}

func (a *mockAddr) Network() string { return a.network }
func (a *mockAddr) String() string  { return a.address }
