package main

import (
	"io"
)

type BitReader struct {
	r io.Reader
	b byte
	n uint
}

func NewBitReader(r io.Reader) *BitReader {
	return &BitReader{
		r: r,
		b: 0,
		n: 8,
	}
}

func (u *BitReader) ReadBit() (bool, error) {
	if u.n >= 8 {
		buf := make([]byte, 1)
		_, err := io.ReadFull(u.r, buf)
		if err != nil {
			return false, err
		}
		u.b = buf[0]
		u.n = 0
	}
	b := (u.b & (1 << u.n)) != 0
	u.n++
	return b, nil
}

func (u *BitReader) ReadByte() (byte, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(u.r, buf)
	if err == io.EOF && u.n < 8 {
		b := u.b >> u.n
		u.n = 8
		return b, nil
	}
	if err != nil {
		return 0, err
	}
	b := (u.b >> u.n) | (buf[0] << (8 - u.n))
	u.b = buf[0]
	return b, nil
}

func (u *BitReader) Read(p []byte) (int, error) {
	n, err := u.r.Read(p)
	for i := 0; i < n; i++ {
		b := p[i]
		p[i] = (u.b >> u.n) | (b << (8 - u.n))
		u.b = b
	}
	if err == io.EOF && u.n < 8 {
		if n >= len(p) {
			err = nil
		} else {
			p[n] = u.b >> u.n
			u.n = 8
			n++
		}
	}
	return n, err
}
