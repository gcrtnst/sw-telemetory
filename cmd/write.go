package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
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

func ValidPath(path string) bool {
	return fs.ValidPath(path) && !(runtime.GOOS == "windows" && strings.ContainsAny(path, `:\`))
}

func GenerateFilepath(root, path string) (string, error) {
	if !ValidPath(path) {
		return "", fmt.Errorf(`invalid path: "%s"`, path)
	}

	path = filepath.FromSlash(path)
	return filepath.Join(root, path), nil
}
