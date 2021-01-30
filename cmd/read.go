package main

import (
	"net"
	"time"
)

type Reader struct {
	recv *Receiver
	buf  []byte
}

func NewReader(lis net.Listener) *Reader {
	return NewReaderTimeout(lis, 0)
}

func NewReaderTimeout(lis net.Listener, timeout time.Duration) *Reader {
	return &Reader{
		recv: NewReceiverTimeout(lis, timeout),
		buf:  []byte{},
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
	buf, err := r.recv.Recv()
	if err != nil {
		return 0, err
	}
	r.buf = buf
	n = r.readBuffer(p)
	return n, nil
}

func (r *Reader) Close() error {
	return r.recv.Close()
}

func (r *Reader) readBuffer(p []byte) int {
	n := copy(p, r.buf)
	r.buf = r.buf[n:]
	return n
}
