package main

import (
	"bytes"
	"runtime"
	"testing"
)

func TestParseWriteRequestNormal(t *testing.T) {
	cases := []struct {
		inReq     []byte
		wantPath  string
		wantData  []byte
		wantIsErr bool
	}{
		{
			inReq:    []byte("\x00"),
			wantPath: "",
			wantData: []byte{},
		},
		{
			inReq:    []byte("\x00data"),
			wantPath: "",
			wantData: []byte("data"),
		},
		{
			inReq:    []byte("path\x00"),
			wantPath: "path",
			wantData: []byte{},
		},
		{
			inReq:    []byte("path\x00data"),
			wantPath: "path",
			wantData: []byte("data"),
		},
	}

	for i, c := range cases {
		inReq := make([]byte, len(c.inReq))
		copy(inReq, c.inReq)
		gotReq, gotErr := ParseWriteRequest(inReq)
		for i := range inReq {
			inReq[i]++
		}

		if gotErr != nil {
			t.Errorf("case %d: error: %v", i, gotErr)
		}
		if gotReq == nil {
			t.Errorf("case %d: req is nil", i)
		}
		if gotReq != nil && gotReq.Path != c.wantPath {
			t.Errorf("case %d: req.Path: expected %#v, got %#v", i, c.wantPath, gotReq.Path)
		}
		if gotReq != nil && (((gotReq.Data == nil) != (c.wantData == nil)) || (!bytes.Equal(gotReq.Data, c.wantData))) {
			t.Errorf("case %d: req.Data: expected %#v, got %#v", i, c.wantData, gotReq.Data)
		}
	}
}

func TestValidPath(t *testing.T) {
	testValidPath(t, []validPathTestCase{
		{
			inPath: ".",
			wantOK: true,
		},
		{
			inPath: "../",
			wantOK: false,
		},
		{
			inPath: "/",
			wantOK: false,
		},
	})
}

func TestValidPathNonWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.SkipNow()
	}

	testValidPath(t, []validPathTestCase{
		{
			inPath: `..\`,
			wantOK: true,
		},
		{
			inPath: `C:`,
			wantOK: true,
		},
	})
}

func TestValidPathWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.SkipNow()
	}

	testValidPath(t, []validPathTestCase{
		{
			inPath: `..\`,
			wantOK: false,
		},
		{
			inPath: `C:`,
			wantOK: false,
		},
	})
}

type validPathTestCase struct {
	inPath string
	wantOK bool
}

func testValidPath(t *testing.T, cases []validPathTestCase) {
	for i, c := range cases {
		gotOK := ValidPath(c.inPath)
		if gotOK != c.wantOK {
			t.Errorf("case %d: expected %t, got %t", i, c.wantOK, gotOK)
		}
	}
}
