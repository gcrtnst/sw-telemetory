package main

import (
	"net"
)

type Reader struct {
	rx  *Receiver
	buf []byte
}

func NewReader(lis net.Listener) *Reader {
	return &Reader{
		rx:  NewReceiver(lis),
		buf: []byte{},
	}
}

func (r *Reader) Read(p []byte) (int, error) {
	if len(p) <= 0 {
		return 0, nil
	}
	n := r.readBuffer(p)
	if n > 0 {
		return n, nil
	}
	buf, err := r.rx.Recv()
	if err != nil {
		return 0, err
	}
	r.buf = buf
	n = r.readBuffer(p)
	return n, nil
}

func (r *Reader) Close() error {
	return r.rx.Close()
}

func (r *Reader) readBuffer(p []byte) int {
	n := copy(p, r.buf)
	r.buf = r.buf[n:]
	return n
}
