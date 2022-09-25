package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestRouterServiceServeNormal(t *testing.T) {
	cases := []struct {
		inMockSID   byte
		inMockResp  []byte
		inReq       []byte
		wantResp    []byte
		wantMockReq []byte
	}{
		{
			inMockSID:   0x80,
			inMockResp:  []byte("RESP"),
			inReq:       []byte("\x80REQ"),
			wantResp:    []byte("\x80RESP"),
			wantMockReq: []byte("REQ"),
		},
		{
			inMockSID:   0x80,
			inMockResp:  nil,
			inReq:       []byte("\x80REQ"),
			wantResp:    []byte("\x80"),
			wantMockReq: []byte("REQ"),
		},
	}

	for i, c := range cases {
		ms := &mockService{
			req:  nil,
			resp: c.inMockResp,
			err:  nil,
		}
		inR := &RouterService{
			M: map[byte]Service{c.inMockSID: ms},
		}

		gotResp, gotErr := inR.Serve(c.inReq)
		if !bytes.Equal(gotResp, c.wantResp) {
			t.Errorf("case %d: resp: expected %#v, got %#v", i, c.wantResp, gotResp)
		}
		if gotErr != nil {
			t.Errorf(`case %d: err: expected nil, got "%v"`, i, gotErr)
		}
		if !bytes.Equal(ms.req, c.wantMockReq) {
			t.Errorf("case %d: mock req: expected %#v, got %#v", i, c.wantMockReq, ms.req)
		}
	}
}

func TestRouterServiceServeError(t *testing.T) {
	cases := []struct {
		inMockSID   byte
		inMockResp  []byte
		inMockErr   error
		inReq       []byte
		wantMockReq []byte
	}{
		{
			inMockSID:   0x80,
			inMockResp:  nil,
			inMockErr:   nil,
			inReq:       nil,
			wantMockReq: nil,
		},
		{
			inMockSID:   0x80,
			inMockResp:  nil,
			inMockErr:   nil,
			inReq:       []byte{},
			wantMockReq: nil,
		},
		{
			inMockSID:   0xC5,
			inMockResp:  nil,
			inMockErr:   nil,
			inReq:       []byte{0xC5},
			wantMockReq: nil,
		},
		{
			inMockSID:   0x80,
			inMockResp:  nil,
			inMockErr:   nil,
			inReq:       []byte{0x81},
			wantMockReq: nil,
		},
		{
			inMockSID:   0x80,
			inMockResp:  nil,
			inMockErr:   errors.New(""),
			inReq:       []byte("\x80REQ"),
			wantMockReq: []byte("REQ"),
		},
	}

	for i, c := range cases {
		ms := &mockService{
			req:  nil,
			resp: c.inMockResp,
			err:  c.inMockErr,
		}
		inR := &RouterService{
			M: map[byte]Service{c.inMockSID: ms},
		}

		gotResp, gotErr := inR.Serve(c.inReq)
		if len(gotResp) < 1 || gotResp[0] != 0xC5 {
			t.Errorf("case %d: resp: wrong resp %#v", i, gotResp)
		}
		if gotErr != nil {
			t.Errorf(`case %d: err: expected nil, got "%v"`, i, gotErr)
		}
		if ((ms.req == nil) != (c.wantMockReq == nil)) || !bytes.Equal(ms.req, c.wantMockReq) {
			t.Errorf("case %d: mock req: expected %#v, got %#v", i, c.wantMockReq, ms.req)
		}
	}
}

type mockService struct {
	req  []byte
	resp []byte
	err  error
}

func (s *mockService) Serve(req []byte) ([]byte, error) {
	s.req = req
	return s.resp, s.err
}
