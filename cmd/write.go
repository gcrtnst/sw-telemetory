package main

import (
	"bytes"
	"errors"
)

type WriteRequest struct {
	Path string
	Data []byte
}

func ParseWriteRequest(req []byte) (*WriteRequest, error) {
	i := bytes.IndexByte(req, 0)
	if i < 0 {
		return nil, errors.New("invalid write request")
	}

	src_data := req[i+1:]
	dst_data := make([]byte, len(src_data))
	copy(dst_data, src_data)
	return &WriteRequest{
		Path: string(req[:i]),
		Data: dst_data,
	}, nil
}
