package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMachineExec(t *testing.T) {
	root, err := ioutil.TempDir("", "*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err := os.RemoveAll(root)
		if err != nil {
			t.Error(err)
		}
	})

	cfg := NewMachineConfig()
	cfg.Root = root
	cfg.Log = log.New(ioutil.Discard, "", log.LstdFlags)
	m := NewMachine(cfg)

	tests := []struct {
		cmd string
		err bool
	}{
		{"", true},
		{"n", false},
		{"w", false},
		{"/", true},
	}
	for _, tt := range tests {
		err := m.Exec(tt.cmd)
		if (err != nil) != tt.err {
			t.Errorf("input %#v: wrong error", tt.cmd)
		}
	}
}

func TestMachineExecNew(t *testing.T) {
	cfg := NewMachineConfig()
	cfg.Log = log.New(ioutil.Discard, "", log.LstdFlags)
	m := NewMachine(cfg)

	tests := []struct {
		title string
		err   bool
	}{
		{"", false},
		{"TEST", false},
		{"TEST" + string(os.PathSeparator), true},
	}
	for _, tt := range tests {
		err := m.ExecNew(tt.title)
		if (err != nil) != tt.err {
			t.Errorf("input %#v: wrong error", tt.title)
		}
	}
}

func TestMachineExecWrite(t *testing.T) {
	root, err := ioutil.TempDir("", "*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err := os.RemoveAll(root)
		if err != nil {
			t.Error(err)
		}
	})

	cfg := NewMachineConfig()
	cfg.Root = root
	cfg.Log = log.New(ioutil.Discard, "", log.LstdFlags)
	m := NewMachine(cfg)

	s := "TEST"
	m.ExecWrite(s)
	f, err := os.Open(m.fpath)
	if err != nil {
		t.Error(err)
		return
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
		return
	}
	if string(b) != s+"\n" {
		t.Error("wrong data was written")
		return
	}
}

func TestGenerateFilepath(t *testing.T) {
	testSep := string(os.PathSeparator)
	testT := time.Date(2006, time.January, 2, 15, 4, 5, 999999999, time.FixedZone("UTC-7", -7*60*60))

	cases := []struct {
		inRoot    string
		inTitle   string
		inExt     string
		inT       time.Time
		wantFpath string
		wantIsErr bool
	}{
		{
			inRoot:    "root",
			inTitle:   "title",
			inExt:     ".ext",
			inT:       testT,
			wantFpath: strings.ReplaceAll("root/title/title-20060102150405.ext", "/", testSep),
			wantIsErr: false,
		},
		{
			inRoot:    "root",
			inTitle:   "title",
			inExt:     "",
			inT:       testT,
			wantFpath: strings.ReplaceAll("root/title/title-20060102150405", "/", testSep),
			wantIsErr: false,
		},
		{
			inRoot:    "",
			inTitle:   "title",
			inExt:     ".ext",
			inT:       testT,
			wantFpath: "",
			wantIsErr: true,
		},
		{
			inRoot:    "root",
			inTitle:   "",
			inExt:     ".ext",
			inT:       testT,
			wantFpath: "",
			wantIsErr: true,
		},
		{
			inRoot:    "root",
			inTitle:   "/title",
			inExt:     ".ext",
			inT:       testT,
			wantFpath: "",
			wantIsErr: true,
		},
		{
			inRoot:    "root",
			inTitle:   "title.ext",
			inExt:     ".ext",
			inT:       testT,
			wantFpath: "",
			wantIsErr: true,
		},
		{
			inRoot:    "root",
			inTitle:   "title",
			inExt:     "/ext",
			inT:       testT,
			wantFpath: "",
			wantIsErr: true,
		},
		{
			inRoot:    "root",
			inTitle:   "title",
			inExt:     "ext/",
			inT:       testT,
			wantFpath: "",
			wantIsErr: true,
		},
		{
			inRoot:    "root",
			inTitle:   "title",
			inExt:     "ext",
			inT:       testT,
			wantFpath: "",
			wantIsErr: true,
		},
	}

	for i, c := range cases {
		gotFpath, gotErr := GenerateFilepath(c.inRoot, c.inTitle, c.inExt, c.inT)
		gotIsErr := gotErr != nil

		if gotFpath != c.wantFpath {
			t.Errorf(`case %d: fpath: expected "%s", got "%s"`, i, c.wantFpath, gotFpath)
		}
		if gotIsErr != c.wantIsErr {
			t.Errorf("case %d: err: expected %t, got %t", i, c.wantIsErr, gotIsErr)
		}
	}
}

func TestIndexPathSeparator(t *testing.T) {
	sep := string(os.PathSeparator)

	cases := []struct {
		inS   string
		wantI int
	}{
		{"", -1},
		{"test", -1},
		{sep + "test", 0},
		{"te" + sep + "st", 2},
		{"test" + sep, 4},
	}

	for i, c := range cases {
		gotI := IndexPathSeparator(c.inS)
		if gotI != c.wantI {
			t.Errorf("case %d: expected %d, got %d", i, c.wantI, gotI)
		}
	}
}
