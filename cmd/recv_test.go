package main

import (
	"bytes"
	"errors"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestReceiverRecv(t *testing.T) {
	cases := []struct {
		inLis    *mockListener
		wantBody []byte
		wantErr  error
	}{
		{
			inLis: &mockListener{
				acceptConn: &mockConn{
					readInner: strings.NewReader("GETS"),
					closeErr:  nil,
				},
				acceptErr: nil,
			},
			wantBody: nil,
			wantErr:  ErrInvalidChunkPrefix,
		},
		{
			inLis: &mockListener{
				acceptConn: &mockConn{
					readInner: strings.NewReader("GET abc HTTP/1.1\r\n"),
					closeErr:  nil,
				},
				acceptErr: nil,
			},
			wantBody: []byte("abc"),
			wantErr:  nil,
		},
	}

	for i, c := range cases {
		inRecv := NewReceiver(c.inLis)
		gotBody, gotErr := inRecv.Recv()
		if !((gotBody == nil && c.wantBody == nil) || (gotBody != nil && c.wantBody != nil && bytes.Equal(gotBody, c.wantBody))) {
			t.Errorf("case %d: body: expected %#v, got %#v", i, c.wantBody, gotBody)
		}
		if gotErr != c.wantErr {
			t.Errorf(`case %d: err: expected "%s", got "%s"`, i, c.wantErr, gotErr)
		}
	}
}

func TestReceiverRecvChunk(t *testing.T) {
	testErr := errors.New("")

	cases := []struct {
		inLis     *mockListener
		wantChunk []byte
		wantErr   error
	}{
		{
			inLis: &mockListener{
				acceptConn: nil,
				acceptErr:  testErr,
			},
			wantChunk: []byte{},
			wantErr:   testErr,
		},
		{
			inLis: &mockListener{
				acceptConn: &mockConn{
					readInner: strings.NewReader("GETS"),
					closeErr:  nil,
				},
				acceptErr: nil,
			},
			wantChunk: []byte("GETS"),
			wantErr:   ErrInvalidChunkPrefix,
		},
		{
			inLis: &mockListener{
				acceptConn: &mockConn{
					readInner: strings.NewReader("GET  HTTP/1.1\r\n"),
					closeErr:  testErr,
				},
				acceptErr: nil,
			},
			wantChunk: []byte("GET  HTTP/1.1\r\n"),
			wantErr:   testErr,
		},
		{
			inLis: &mockListener{
				acceptConn: &mockConn{
					readInner: strings.NewReader("GET  HTTP/1.1\r\n"),
					closeErr:  nil,
				},
				acceptErr: nil,
			},
			wantChunk: []byte("GET  HTTP/1.1\r\n"),
			wantErr:   nil,
		},
		{
			inLis: &mockListener{
				acceptConn: &mockByteReaderConn{
					readInner: strings.NewReader("GET  HTTP/1.1\r\n"),
				},
				acceptErr: nil,
			},
			wantChunk: []byte("GET  HTTP/1.1\r\n"),
			wantErr:   nil,
		},
	}

	for i, c := range cases {
		inRecv := NewReceiver(c.inLis)
		gotChunk, gotErr := inRecv.RecvChunk()

		if !bytes.Equal(gotChunk, c.wantChunk) {
			t.Errorf("case %d: chunk: expected %#v, got %#v", i, c.wantChunk, gotChunk)
		}
		if gotErr != c.wantErr {
			t.Errorf(`case %d: err: expected "%s", got "%s"`, i, c.wantErr, gotErr)
		}
		if conn, ok := c.inLis.acceptConn.(*mockConn); ok && !conn.closeDone {
			t.Errorf("case %d: conn not closed", i)
		}
	}
}

func TestReceiverClose(t *testing.T) {
	inRecv := NewReceiver(&mockListener{closeErr: errors.New("")})
	gotErr := inRecv.Close()
	if gotErr != inRecv.lis.(*mockListener).closeErr {
		t.Errorf(`err: expected "%s", got "%s"`, inRecv.lis.(*mockListener).closeErr, gotErr)
	}
	if !inRecv.lis.(*mockListener).closeDone {
		t.Error("conn not closed")
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

type mockListener struct {
	acceptConn net.Conn
	acceptErr  error
	closeErr   error
	closeDone  bool
}

func (m *mockListener) Accept() (net.Conn, error) { return m.acceptConn, m.acceptErr }
func (m *mockListener) Addr() net.Addr            { panic("not implemented") }

func (m *mockListener) Close() error {
	m.closeDone = true
	return m.closeErr
}

type mockConn struct {
	readInner *strings.Reader
	closeErr  error
	closeDone bool
}

func (m *mockConn) Read(b []byte) (int, error)         { return m.readInner.Read(b) }
func (m *mockConn) Write(b []byte) (int, error)        { panic("not implemented") }
func (m *mockConn) LocalAddr() net.Addr                { panic("not implemented") }
func (m *mockConn) RemoteAddr() net.Addr               { panic("not implemented") }
func (m *mockConn) SetDeadline(t time.Time) error      { panic("not implemented") }
func (m *mockConn) SetReadDeadline(t time.Time) error  { panic("not implemented") }
func (m *mockConn) SetWriteDeadline(t time.Time) error { panic("not implemented") }

func (m *mockConn) Close() error {
	m.closeDone = true
	return m.closeErr
}

type mockByteReaderConn struct {
	readInner *strings.Reader
}

func (m *mockByteReaderConn) ReadByte() (byte, error)            { return m.readInner.ReadByte() }
func (m *mockByteReaderConn) Read(b []byte) (int, error)         { panic("not implemented") }
func (m *mockByteReaderConn) Write(b []byte) (int, error)        { panic("not implemented") }
func (m *mockByteReaderConn) Close() error                       { return nil }
func (m *mockByteReaderConn) LocalAddr() net.Addr                { panic("not implemented") }
func (m *mockByteReaderConn) RemoteAddr() net.Addr               { panic("not implemented") }
func (m *mockByteReaderConn) SetDeadline(t time.Time) error      { panic("not implemented") }
func (m *mockByteReaderConn) SetReadDeadline(t time.Time) error  { panic("not implemented") }
func (m *mockByteReaderConn) SetWriteDeadline(t time.Time) error { panic("not implemented") }

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

type mockListenerOld struct {
	conns  []net.Conn
	err    error
	closed bool
}

func (l *mockListenerOld) Accept() (net.Conn, error) {
	if len(l.conns) <= 0 {
		return nil, l.err
	}
	conn := l.conns[0]
	l.conns = l.conns[1:]
	return conn, nil
}

func (l *mockListenerOld) Close() error {
	l.closed = true
	return nil
}

func (l *mockListenerOld) Addr() net.Addr {
	return &mockAddrOld{
		network: "TEST",
		address: "TEST",
	}
}

type mockConnOld struct {
	rd       *strings.Reader
	closed   bool
	deadline time.Time
}

func (c *mockConnOld) Read(b []byte) (int, error) {
	return c.rd.Read(b)
}

func (c *mockConnOld) Write(b []byte) (int, error) {
	panic("not implemented")
}

func (c *mockConnOld) Close() error {
	c.closed = true
	return nil
}

func (c *mockConnOld) LocalAddr() net.Addr {
	panic("not implemented")
}

func (c *mockConnOld) RemoteAddr() net.Addr {
	panic("not implemented")
}

func (c *mockConnOld) SetDeadline(t time.Time) error {
	c.deadline = t
	return nil
}

func (c *mockConnOld) SetReadDeadline(t time.Time) error {
	panic("not implemented")
}

func (c *mockConnOld) SetWriteDeadline(t time.Time) error {
	panic("not implemented")
}

type mockAddrOld struct {
	network string
	address string
}

func (a *mockAddrOld) Network() string { return a.network }
func (a *mockAddrOld) String() string  { return a.address }
