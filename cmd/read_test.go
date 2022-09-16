package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestReaderRead(t *testing.T) {
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
			recv: inRx,
			buf:  c.inBuf,
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
