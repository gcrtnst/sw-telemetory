package main

import (
	"errors"
	"io"
	"net"
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

func TestReadChunk(t *testing.T) {
	testErr := errors.New("")

	cases := []struct {
		inS     string
		inErr   error
		wantBuf string
		wantErr error
		wantLen int
	}{
		{"", io.EOF, "", io.ErrUnexpectedEOF, 0},
		{"G", io.EOF, "G", io.ErrUnexpectedEOF, 0},
		{"P", io.EOF, "P", ErrInvalidChunkPrefix, 0},
		{"GET ", io.EOF, "GET ", io.ErrUnexpectedEOF, 0},
		{"GETS", io.EOF, "GETS", ErrInvalidChunkPrefix, 0},
		{"GETS HTTP/1.1\r\n", io.EOF, "GETS", ErrInvalidChunkPrefix, 11},
		{"", testErr, "", testErr, 0},
		{"GET ", testErr, "GET ", testErr, 0},
		{"GET HTTP/1.1\r\n", io.EOF, "GET HTTP/1.1\r\n", io.ErrUnexpectedEOF, 0},
		{"GET  HTTP/1.2\r\n", io.EOF, "GET  HTTP/1.2\r\n", io.ErrUnexpectedEOF, 0},
		{"GET  HTTP/1.1\r\n", io.EOF, "GET  HTTP/1.1\r\n", nil, 0},
		{"GET   HTTP/1.1\r\n", io.EOF, "GET   HTTP/1.1\r\n", nil, 0},
		{"GET  HTTP/1.1\r\n ", io.EOF, "GET  HTTP/1.1\r\n", nil, 1},
	}

	for i, c := range cases {
		in := &mockByteReader{br: strings.NewReader(c.inS), err: c.inErr}
		gotBuf, gotErr := ReadChunk(in)
		if string(gotBuf) != c.wantBuf {
			t.Errorf("case %d: buf: expected %#v, got %#v", i, c.wantBuf, string(gotBuf))
		}
		if gotErr != c.wantErr {
			t.Errorf(`case %d: err: expected "%s", got "%s"`, i, c.wantErr, gotErr)
		}
		if in.br.Len() != c.wantLen {
			t.Errorf("case %d: len: expected %d, got %d", i, c.wantLen, in.br.Len())
		}
	}
}

type mockByteReader struct {
	br  *strings.Reader
	err error
}

func (m *mockByteReader) ReadByte() (byte, error) {
	b, err := m.br.ReadByte()
	if err == io.EOF {
		return b, m.err
	}
	return b, err
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
