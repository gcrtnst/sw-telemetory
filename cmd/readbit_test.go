package main

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestBitReaderReadBit(t *testing.T) {
	tests := []struct {
		inRd     []byte
		inRB     byte
		inRN     uint
		outRet   bool
		outErr   error
		outRdLen int
		outRB    byte
		outRN    uint
	}{
		{
			inRd:     []byte{},
			inRB:     0,
			inRN:     8,
			outRet:   false,
			outErr:   io.EOF,
			outRdLen: 0,
			outRB:    0,
			outRN:    8,
		},
		{
			inRd:     []byte{0b11111110},
			inRB:     0,
			inRN:     8,
			outRet:   false,
			outErr:   nil,
			outRdLen: 0,
			outRB:    0xFE,
			outRN:    1,
		},
		{
			inRd:     []byte{0b00000001},
			inRB:     0,
			inRN:     8,
			outRet:   true,
			outErr:   nil,
			outRdLen: 0,
			outRB:    0x01,
			outRN:    1,
		},
		{
			inRd:     []byte{0b11111111},
			inRB:     0x7F,
			inRN:     7,
			outRet:   false,
			outErr:   nil,
			outRdLen: 1,
			outRB:    0x7F,
			outRN:    8,
		},
		{
			inRd:     []byte{0b00000000},
			inRB:     0x80,
			inRN:     7,
			outRet:   true,
			outErr:   nil,
			outRdLen: 1,
			outRB:    0x80,
			outRN:    8,
		},
	}
	for _, tt := range tests {
		rd := bytes.NewReader(tt.inRd)
		r := &BitReader{r: rd, b: tt.inRB, n: tt.inRN}
		ret, err := r.ReadBit()
		if ret != tt.outRet {
			t.Errorf("case %v: ret: expected %t, got %t", tt, tt.outRet, ret)
		}
		if err != tt.outErr {
			t.Errorf("case %v: err: expected %#v, got %#v", tt, tt.outErr, err)
		}
		if rd.Len() != tt.outRdLen {
			t.Errorf("case %v: rdlen: expected %d, got %d", tt, tt.outRdLen, rd.Len())
		}
		if r.b != tt.outRB {
			t.Errorf("case %v: r.b: expected %d, got %d", tt, tt.outRB, r.b)
		}
		if r.n != tt.outRN {
			t.Errorf("case %v: r.n: expected %d, got %d", tt, tt.outRN, r.n)
		}
	}
}

func TestBitReaderReadByte(t *testing.T) {
	merr := errors.New("TEST")
	tests := []struct {
		inRd   *MockReader
		inRB   byte
		inRN   uint
		outRet byte
		outErr error
		outRdN int
		outRB  byte
		outRN  uint
	}{
		{
			inRd:   &MockReader{b: []byte{}, err: io.EOF},
			inRB:   0b11011110,
			inRN:   8,
			outRet: 0,
			outErr: io.EOF,
			outRdN: 0,
			outRB:  0b11011110,
			outRN:  8,
		},
		{
			inRd:   &MockReader{b: []byte{}, err: io.EOF},
			inRB:   0b11011110,
			inRN:   2,
			outRet: 0b00110111,
			outErr: nil,
			outRdN: 0,
			outRB:  0b11011110,
			outRN:  8,
		},
		{
			inRd:   &MockReader{b: []byte{}, err: merr},
			inRB:   0b11011110,
			inRN:   2,
			outRet: 0,
			outErr: merr,
			outRdN: 0,
			outRB:  0b11011110,
			outRN:  2,
		},
		{
			inRd:   &MockReader{b: []byte{0b10101101, 0b10111110, 0b11101111}, err: nil},
			inRB:   0b11011110,
			inRN:   8,
			outRet: 0b10101101,
			outErr: nil,
			outRdN: 1,
			outRB:  0b10101101,
			outRN:  8,
		},
		{
			inRd:   &MockReader{b: []byte{0b10101101, 0b10111110, 0b11101111}, err: nil},
			inRB:   0b11011110,
			inRN:   2,
			outRet: 0b01110111,
			outErr: nil,
			outRdN: 1,
			outRB:  0b10101101,
			outRN:  2,
		},
	}
	for _, tt := range tests {
		r := &BitReader{r: tt.inRd, b: tt.inRB, n: tt.inRN}
		ret, err := r.ReadByte()
		if ret != tt.outRet {
			t.Errorf("case %v: ret: expected %b, got %b", tt, tt.outRet, ret)
		}
		if err != tt.outErr {
			t.Errorf("case %v: err: expected %#v, got %#v", tt, tt.outErr, err)
		}
		if tt.inRd.n != tt.outRdN {
			t.Errorf("case %v: rdn: expected %d, got %d", tt, tt.outRdN, tt.inRd.n)
		}
		if r.b != tt.outRB {
			t.Errorf("case %v: r.b: expected %d, got %d", tt, tt.outRB, r.b)
		}
		if r.n != tt.outRN {
			t.Errorf("case %v: r.n: expected %d, got %d", tt, tt.outRN, r.n)
		}
	}
}

func TestBitReaderRead(t *testing.T) {
	merr := errors.New("TEST")
	tests := []struct {
		inPLen int
		inRd   *MockReader
		inRB   byte
		inRN   uint
		outP   []byte
		outN   int
		outErr error
		outRdN int
		outRB  byte
		outRN  uint
	}{
		{
			inPLen: 0,
			inRd:   &MockReader{b: []byte{}, err: nil},
			inRB:   0b11011110,
			inRN:   2,
			outP:   []byte{},
			outN:   0,
			outErr: nil,
			outRdN: 0,
			outRB:  0b11011110,
			outRN:  2,
		},
		{
			inPLen: 0,
			inRd:   &MockReader{b: []byte{}, err: merr},
			inRB:   0b11011110,
			inRN:   2,
			outP:   []byte{},
			outN:   0,
			outErr: merr,
			outRdN: 0,
			outRB:  0b11011110,
			outRN:  2,
		},
		{
			inPLen: 0,
			inRd:   &MockReader{b: []byte{}, err: io.EOF},
			inRB:   0b11011110,
			inRN:   2,
			outP:   []byte{},
			outN:   0,
			outErr: nil,
			outRdN: 0,
			outRB:  0b11011110,
			outRN:  2,
		},
		{
			inPLen: 0,
			inRd:   &MockReader{b: []byte{}, err: io.EOF},
			inRB:   0b11011110,
			inRN:   8,
			outP:   []byte{},
			outN:   0,
			outErr: io.EOF,
			outRdN: 0,
			outRB:  0b11011110,
			outRN:  8,
		},
		{
			inPLen: 3,
			inRd:   &MockReader{b: []byte{0b10101101, 0b10111110}, err: nil},
			inRB:   0b11011110,
			inRN:   2,
			outP:   []byte{0b01110111, 0b10101011, 0b00000000},
			outN:   2,
			outErr: nil,
			outRdN: 2,
			outRB:  0b10111110,
			outRN:  2,
		},
		{
			inPLen: 3,
			inRd:   &MockReader{b: []byte{0b10101101, 0b10111110}, err: nil},
			inRB:   0b11011110,
			inRN:   8,
			outP:   []byte{0b10101101, 0b10111110, 0b00000000},
			outN:   2,
			outErr: nil,
			outRdN: 2,
			outRB:  0b10111110,
			outRN:  8,
		},
		{
			inPLen: 3,
			inRd:   &MockReader{b: []byte{0b10101101, 0b10111110}, err: merr},
			inRB:   0b11011110,
			inRN:   2,
			outP:   []byte{0b01110111, 0b10101011, 0b00000000},
			outN:   2,
			outErr: merr,
			outRdN: 2,
			outRB:  0b10111110,
			outRN:  2,
		},
		{
			inPLen: 3,
			inRd:   &MockReader{b: []byte{0b10101101, 0b10111110}, err: io.EOF},
			inRB:   0b11011110,
			inRN:   2,
			outP:   []byte{0b01110111, 0b10101011, 0b00101111},
			outN:   3,
			outErr: io.EOF,
			outRdN: 2,
			outRB:  0b10111110,
			outRN:  8,
		},
		{
			inPLen: 3,
			inRd:   &MockReader{b: []byte{0b10101101, 0b10111110}, err: io.EOF},
			inRB:   0b11011110,
			inRN:   8,
			outP:   []byte{0b10101101, 0b10111110, 0b00000000},
			outN:   2,
			outErr: io.EOF,
			outRdN: 2,
			outRB:  0b10111110,
			outRN:  8,
		},
		{
			inPLen: 2,
			inRd:   &MockReader{b: []byte{0b10101101, 0b10111110, 0b11101111}, err: nil},
			inRB:   0b11011110,
			inRN:   2,
			outP:   []byte{0b01110111, 0b10101011},
			outN:   2,
			outErr: nil,
			outRdN: 2,
			outRB:  0b10111110,
			outRN:  2,
		},
		{
			inPLen: 2,
			inRd:   &MockReader{b: []byte{0b10101101, 0b10111110, 0b11101111}, err: nil},
			inRB:   0b11011110,
			inRN:   8,
			outP:   []byte{0b10101101, 0b10111110},
			outN:   2,
			outErr: nil,
			outRdN: 2,
			outRB:  0b10111110,
			outRN:  8,
		},
		{
			inPLen: 2,
			inRd:   &MockReader{b: []byte{0b10101101, 0b10111110, 0b11101111}, err: merr},
			inRB:   0b11011110,
			inRN:   2,
			outP:   []byte{0b01110111, 0b10101011},
			outN:   2,
			outErr: merr,
			outRdN: 2,
			outRB:  0b10111110,
			outRN:  2,
		},
		{
			inPLen: 2,
			inRd:   &MockReader{b: []byte{0b10101101, 0b10111110, 0b11101111}, err: io.EOF},
			inRB:   0b11011110,
			inRN:   2,
			outP:   []byte{0b01110111, 0b10101011},
			outN:   2,
			outErr: nil,
			outRdN: 2,
			outRB:  0b10111110,
			outRN:  2,
		},
		{
			inPLen: 2,
			inRd:   &MockReader{b: []byte{0b10101101, 0b10111110, 0b11101111}, err: io.EOF},
			inRB:   0b11011110,
			inRN:   8,
			outP:   []byte{0b10101101, 0b10111110},
			outN:   2,
			outErr: io.EOF,
			outRdN: 2,
			outRB:  0b10111110,
			outRN:  8,
		},
	}
	for _, tt := range tests {
		p := make([]byte, tt.inPLen)
		r := &BitReader{r: tt.inRd, b: tt.inRB, n: tt.inRN}
		n, err := r.Read(p)
		if string(p) != string(tt.outP) {
			t.Errorf("case %v: p: expected %v, got %v", tt, tt.outP, p)
		}
		if n != tt.outN {
			t.Errorf("case %v: n: expected %d, got %d", tt, tt.outN, n)
		}
		if err != tt.outErr {
			t.Errorf("case %v: err: expected %#v, got %#v", tt, tt.outErr, err)
		}
		if tt.inRd.n != tt.outRdN {
			t.Errorf("case %v: rdn: expected %d, got %d", tt, tt.outRdN, tt.inRd.n)
		}
		if r.b != tt.outRB {
			t.Errorf("case %v: r.b: expected %b, got %b", tt, tt.outRB, r.b)
		}
		if r.n != tt.outRN {
			t.Errorf("case %v: r.n: expected %d, got %d", tt, tt.outRN, r.n)
		}
	}
}

type MockReader struct {
	b   []byte
	n   int
	err error
}

func (r *MockReader) Read(p []byte) (int, error) {
	r.n = copy(p, r.b)
	return r.n, r.err
}
