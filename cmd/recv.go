package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"time"
)

type Receiver struct {
	lis     net.Listener
	timeout time.Duration
	cg      CloseGroup
}

func NewReceiver(lis net.Listener) *Receiver {
	return NewReceiverTimeout(lis, 0)
}

func NewReceiverTimeout(lis net.Listener, timeout time.Duration) *Receiver {
	return &Receiver{
		lis:     lis,
		timeout: timeout,
		cg:      CloseGroup{},
	}
}

func (r *Receiver) Recv() ([]byte, error) {
	chunk, err := r.recvChunk()
	if err != nil {
		return nil, err
	}
	return ExtractBody(chunk)
}

func (r *Receiver) Close() error {
	err := r.lis.Close()
	r.cg.CloseAll()
	return err
}

func (r *Receiver) recvChunk() ([]byte, error) {
	conn, err := r.lis.Accept()
	if err != nil {
		return nil, err
	}
	defer r.cg.Add(conn).Close()

	if r.timeout > 0 {
		_ = conn.SetDeadline(time.Now().Add(r.timeout))
	} else {
		_ = conn.SetDeadline(time.Time{})
	}
	if conn, ok := conn.(*net.TCPConn); ok {
		_ = conn.SetLinger(0)
	}

	br, ok := conn.(io.ByteReader)
	if !ok {
		br = bufio.NewReader(conn)
	}
	return ReadChunk(br)
}

var (
	chunkPrefix = []byte("GET ")
	chunkSuffix = []byte(" HTTP/1.1\r\n")
)

func ReadChunk(br io.ByteReader) ([]byte, error) {
	return ReadUntil(br, chunkSuffix)
}

func ReadUntil(br io.ByteReader, delim []byte) ([]byte, error) {
	if len(delim) <= 0 {
		return make([]byte, 0), nil
	}

	buf := make([]byte, 0, 4096)
	for {
		b, err := br.ReadByte()
		if err != nil {
			return buf, err
		}
		buf = append(buf, b)
		if bytes.HasSuffix(buf, delim) {
			return buf, nil
		}
	}
}

func ExtractBody(chunk []byte) ([]byte, error) {
	if len(chunk) < len(chunkPrefix)+len(chunkSuffix) || !bytes.HasPrefix(chunk, chunkPrefix) || !bytes.HasSuffix(chunk, chunkSuffix) {
		return make([]byte, 0), errors.New("invalid chunk")
	}
	return chunk[len(chunkPrefix) : len(chunk)-len(chunkSuffix)], nil
}
