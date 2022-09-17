package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
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
	sep := string(os.PathSeparator)
	tests := []struct {
		title, ext string
		err        bool
	}{
		{"title", "ext", false},
		{"title" + sep, "ext", true},
		{"title.", "ext", true},
		{"title", "ext" + sep, true},
		{"title", "ext.", false},
	}
	for _, tt := range tests {
		_, err := GenerateFilepath("root", tt.title, tt.ext)
		if (err != nil) != tt.err {
			t.Errorf("input %#v, %#v: wrong error", tt.title, tt.ext)
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
