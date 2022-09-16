package main

import (
	"bytes"
	"errors"
	"io"
	"net"
	"strings"
	"testing"
)

func TestReaderRead(t *testing.T) {
	cases := []struct {
		inChunk []string
		wantB   []byte
		wantErr error
	}{
		{
			inChunk: []string{},
			wantB:   []byte{},
			wantErr: nil,
		},
		{
			inChunk: []string{""},
			wantB:   []byte{},
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			inChunk: []string{"GET  HTTP/1.1\r\n"},
			wantB:   []byte{},
			wantErr: nil,
		},
		{
			inChunk: []string{"GET BBB HTTP/1.1\r\n"},
			wantB:   []byte("BBB"),
			wantErr: nil,
		},
		{
			inChunk: []string{"GET " + strings.Repeat("B", 8192) + " HTTP/1.1\r\n"},
			wantB:   []byte(strings.Repeat("B", 8192)),
			wantErr: nil,
		},
		{
			inChunk: []string{"GET BBB HTTP/1.1\r\n", "GET  HTTP/1.1\r\n"},
			wantB:   []byte("BBB"),
			wantErr: nil,
		},
		{
			inChunk: []string{"GET BBB HTTP/1.1\r\n", "GET CCC HTTP/1.1\r\n"},
			wantB:   []byte("BBBCCC"),
			wantErr: nil,
		},
		{
			inChunk: []string{"GET BBB HTTP/1.1\r\n", "GET " + strings.Repeat("C", 8192) + " HTTP/1.1\r\n"},
			wantB:   []byte("BBB" + strings.Repeat("C", 8192)),
			wantErr: nil,
		},
		{
			inChunk: []string{"GET BBB HTTP/1.1\r\n", "GET  HTTP/1.1\r\n", "GET CCC HTTP/1.1\r\n"},
			wantB:   []byte("BBBCCC"),
			wantErr: nil,
		},
	}

	for i, c := range cases {
		inLis := &mockMultiConnListener{chunk: c.inChunk, idx: 0}
		inRd := NewReader(inLis)

		gotB, gotErr := io.ReadAll(inRd)
		if !bytes.Equal(gotB, c.wantB) {
			t.Errorf("case %d: b: expected %#v, got %#v", i, c.wantB, gotB)
		}
		if gotErr != c.wantErr {
			t.Errorf(`case %d: err: expected "%s", got "%s"`, i, c.wantErr, gotErr)
		}
	}
}

func TestReaderReadUnit(t *testing.T) {
	cases := []struct {
		inPSize int
		inRxTxt string
		inBuf   []byte
		wantP   []byte
		wantN   int
		wantErr error
		wantBuf []byte
	}{
		{
			inPSize: 0,
			inRxTxt: "",
			inBuf:   []byte{0x42, 0x43},
			wantP:   []byte{},
			wantN:   0,
			wantErr: nil,
			wantBuf: []byte{0x42, 0x43},
		},
		{
			inPSize: 1,
			inRxTxt: "",
			inBuf:   []byte{0x42, 0x43},
			wantP:   []byte{0x42},
			wantN:   1,
			wantErr: nil,
			wantBuf: []byte{0x43},
		},
		{
			inPSize: 2,
			inRxTxt: "",
			inBuf:   []byte{0x42, 0x43},
			wantP:   []byte{0x42, 0x43},
			wantN:   2,
			wantErr: nil,
			wantBuf: []byte{},
		},
		{
			inPSize: 3,
			inRxTxt: "",
			inBuf:   []byte{0x42, 0x43},
			wantP:   []byte{0x42, 0x43, 0x00},
			wantN:   2,
			wantErr: nil,
			wantBuf: []byte{},
		},
		{
			inPSize: 1,
			inRxTxt: "",
			inBuf:   []byte{},
			wantP:   []byte{0x00},
			wantN:   0,
			wantErr: io.ErrUnexpectedEOF,
			wantBuf: []byte{},
		},
		{
			inPSize: 1,
			inRxTxt: "GET  HTTP/1.1\r\n",
			inBuf:   []byte{},
			wantP:   []byte{0x00},
			wantN:   0,
			wantErr: nil,
			wantBuf: []byte{},
		},
		{
			inPSize: 1,
			inRxTxt: "GET BC HTTP/1.1\r\n",
			inBuf:   []byte{},
			wantP:   []byte{0x42},
			wantN:   1,
			wantErr: nil,
			wantBuf: []byte{0x43},
		},
		{
			inPSize: 2,
			inRxTxt: "GET BC HTTP/1.1\r\n",
			inBuf:   []byte{},
			wantP:   []byte{0x42, 0x43},
			wantN:   2,
			wantErr: nil,
			wantBuf: []byte{},
		},
		{
			inPSize: 3,
			inRxTxt: "GET BC HTTP/1.1\r\n",
			inBuf:   []byte{},
			wantP:   []byte{0x42, 0x43, 0x00},
			wantN:   2,
			wantErr: nil,
			wantBuf: []byte{},
		},
	}

	for i, c := range cases {
		inRx := NewReceiver(&mockListener{
			acceptConn: &mockConn{
				readInner: strings.NewReader(c.inRxTxt),
				closeErr:  nil,
			},
			acceptErr: nil,
			closeErr:  nil,
		})
		inRd := &Reader{
			rx:  inRx,
			buf: c.inBuf,
		}

		gotP := make([]byte, c.inPSize)
		gotN, gotErr := inRd.Read(gotP)

		if !bytes.Equal(gotP, c.wantP) {
			t.Errorf("case %d: p: expected %#v, got %#v", i, c.wantP, gotP)
		}
		if gotN != c.wantN {
			t.Errorf("case %d: n: expected %d, got %d", i, c.wantN, gotN)
		}
		if gotErr != c.wantErr {
			t.Errorf(`case %d: err: expected "%s", got "%s"`, i, c.wantErr, gotErr)
		}
		if !bytes.Equal(inRd.buf, c.wantBuf) {
			t.Errorf("case %d: buf: expected %#v, got %#v", i, c.wantBuf, inRd.buf)
		}
	}
}

func TestReaderClose(t *testing.T) {
	inLis := &mockListener{closeErr: errors.New("")}
	inRd := NewReader(inLis)
	gotErr := inRd.Close()
	if gotErr != inLis.closeErr {
		t.Errorf(`err: expected "%s", got "%s"`, inLis.closeErr, gotErr)
	}
	if !inLis.closeDone {
		t.Errorf("lis not closed")
	}
}

type mockMultiConnListener struct {
	chunk []string
	idx   int
}

func (m *mockMultiConnListener) Close() error   { return nil }
func (m *mockMultiConnListener) Addr() net.Addr { panic("not implemented") }

func (m *mockMultiConnListener) Accept() (net.Conn, error) {
	if m.idx >= len(m.chunk) {
		return nil, io.EOF
	}
	conn := &mockConn{readInner: strings.NewReader(m.chunk[m.idx]), closeErr: nil}
	m.idx++
	return conn, nil
}
