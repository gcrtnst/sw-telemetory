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

func (rx *Receiver) Recv() ([]byte, error) {
	chunk, err := rx.RecvChunk()
	if err != nil {
		return nil, err
	}
	return chunk[len(chunkPrefix) : len(chunk)-len(chunkSuffix)], nil
}

func (rx *Receiver) RecvChunk() (chunk []byte, err error) {
	conn, errAccept := rx.lis.Accept()
	if errAccept != nil {
		return []byte{}, errAccept
	}
	defer CloseCatch(rx.cg.Add(conn), &err)

	br, ok := conn.(io.ByteReader)
	if !ok {
		br = bufio.NewReader(conn)
	}
	return ReadChunk(br)
}

func (rx *Receiver) Close() error {
	err := rx.lis.Close()
	rx.cg.CloseAll()
	return err
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
