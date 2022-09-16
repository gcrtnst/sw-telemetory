package main

import (
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestReaderRead(t *testing.T) {
	tests := []struct {
		inPSize   int
		inRdBuf   string
		inRecvBuf string
		outP      string
		outN      int
		outErr    error
		outRdBuf  string
	}{
		{
			inPSize:   0,
			inRdBuf:   "",
			inRecvBuf: "",
			outP:      "",
			outN:      0,
			outErr:    nil,
			outRdBuf:  "",
		},
		{
			inPSize:   1,
			inRdBuf:   "TEST",
			inRecvBuf: "",
			outP:      "T",
			outN:      1,
			outErr:    nil,
			outRdBuf:  "EST",
		},
		{
			inPSize:   5,
			inRdBuf:   "TEST",
			inRecvBuf: "",
			outP:      "TEST\x00",
			outN:      4,
			outErr:    nil,
			outRdBuf:  "",
		},
		{
			inPSize:   1,
			inRdBuf:   "",
			inRecvBuf: "",
			outP:      "\x00",
			outN:      0,
			outErr:    io.ErrUnexpectedEOF,
			outRdBuf:  "",
		},
		{
			inPSize:   1,
			inRdBuf:   "",
			inRecvBuf: string(chunkPrefix) + string(chunkSuffix),
			outP:      "\x00",
			outN:      0,
			outErr:    nil,
			outRdBuf:  "",
		},
		{
			inPSize:   1,
			inRdBuf:   "",
			inRecvBuf: string(chunkPrefix) + "TEST" + string(chunkSuffix),
			outP:      "T",
			outN:      1,
			outErr:    nil,
			outRdBuf:  "EST",
		},
		{
			inPSize:   5,
			inRdBuf:   "",
			inRecvBuf: string(chunkPrefix) + "TEST" + string(chunkSuffix),
			outP:      "TEST\x00",
			outN:      4,
			outErr:    nil,
			outRdBuf:  "",
		},
	}
	for i, tt := range tests {
		conn := &mockConnOld{
			rd:       strings.NewReader(tt.inRecvBuf),
			closed:   false,
			deadline: time.Time{},
		}
		lis := &mockListenerOld{
			conns: []net.Conn{conn},
		}
		recv := &Receiver{
			lis: lis,
			cg:  CloseGroup{},
		}
		rd := &Reader{
			recv: recv,
			buf:  []byte(tt.inRdBuf),
		}
		p := make([]byte, tt.inPSize)
		n, err := rd.Read(p)
		if string(p) != tt.outP {
			t.Errorf("case %d: expected \"%s\", got \"%s\"", i, string(p), tt.outP)
		}
		if n != tt.outN {
			t.Errorf("case %d: expected %d, got %d", i, n, tt.outN)
		}
		if err != tt.outErr {
			t.Errorf("case %d: expected %#v, got %#v", i, err, tt.outErr)
		}
		if string(rd.buf) != tt.outRdBuf {
			t.Errorf("case %d: expected \"%s\", got \"%s\"", i, string(rd.buf), tt.outRdBuf)
		}
	}
}
