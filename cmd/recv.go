package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
)

var ErrInvalidChunkPrefix = errors.New("invalid chunk prefix")

var (
	chunkPrefix = []byte("GET ")
	chunkSuffix = []byte(" HTTP/1.1\r\n")
)

type Receiver struct {
	lis net.Listener
	cg  CloseGroup
}

func NewReceiver(lis net.Listener) *Receiver {
	return &Receiver{
		lis: lis,
		cg:  CloseGroup{},
	}
}

func (r *Receiver) Recv() ([]byte, error) {
	chunk, err := r.RecvChunk()
	if err != nil {
		return nil, err
	}
	return chunk[len(chunkPrefix) : len(chunk)-len(chunkSuffix)], nil
}

func (r *Receiver) Close() error {
	err := r.lis.Close()
	r.cg.CloseAll()
	return err
}

func (r *Receiver) RecvChunk() (chunk []byte, err error) {
	conn, errAccept := r.lis.Accept()
	if errAccept != nil {
		return []byte{}, errAccept
	}
	defer r.cg.Add(conn).CloseCatch(&err)

	br, ok := conn.(io.ByteReader)
	if !ok {
		br = bufio.NewReader(conn)
	}
	return ReadChunk(br)
}

func ReadChunk(br io.ByteReader) ([]byte, error) {
	buf := make([]byte, 0, 8192)
	for {
		b, err := br.ReadByte()
		if err == io.EOF {
			return buf, io.ErrUnexpectedEOF
		}
		if err != nil {
			return buf, err
		}

		buf = append(buf, b)
		if len(buf) <= len(chunkPrefix) && b != chunkPrefix[len(buf)-1] {
			return buf, ErrInvalidChunkPrefix
		}
		if len(buf) >= len(chunkPrefix)+len(chunkSuffix) && bytes.HasSuffix(buf, chunkSuffix) {
			return buf, nil
		}
	}
}
