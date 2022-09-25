package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type WriteService struct {
	Root string
}

func (s *WriteService) Serve(req []byte) ([]byte, error) {
	p, err := ParseWriteRequest(req)
	if err != nil {
		return nil, err
	}

	fpath, err := GenerateFilepath(s.Root, p.Path)
	if err != nil {
		return nil, err
	}

	err = WriteFile(fpath, p.Data, false)
	if err != nil {
		return nil, err
	}

	return []byte{}, nil
}

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

func WriteFile(name string, data []byte, trunc bool) (err error) {
	flag := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	if trunc {
		flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}

	errMkdir := os.MkdirAll(filepath.Dir(name), 0o777)
	if errMkdir != nil {
		return errMkdir
	}

	f, errOpen := os.OpenFile(name, flag, 0o666)
	if errOpen != nil {
		return errOpen
	}
	defer func() {
		errClose := f.Close()
		if err == nil {
			err = errClose
		}
	}()

	if len(data) > 0 {
		_, errWrite := f.Write(data)
		if errWrite != nil {
			return errWrite
		}
	}
	return nil
}
