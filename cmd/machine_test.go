package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	tmp := t.TempDir()
	root := filepath.Join(tmp, "root")
	testT := time.Date(2006, time.January, 2, 15, 4, 5, 999999999, time.FixedZone("UTC-7", -7*60*60))

	cases := []struct {
		inM       *Machine
		inS       string
		inT       time.Time
		inFpath   string
		inB       []byte
		wantIsErr bool
		wantFpath string
		wantB     []byte
	}{
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: "",
			},
			inS:       "",
			inT:       testT,
			inFpath:   filepath.Join(root, "dummy"),
			inB:       []byte{},
			wantIsErr: false,
			wantFpath: filepath.Join(root, "default", "default-20060102150405.ext"),
			wantB:     []byte{},
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: "",
			},
			inS:       "",
			inT:       testT,
			inFpath:   filepath.Join(root, "default", "default-20060102150405.ext"),
			inB:       []byte("in"),
			wantIsErr: false,
			wantFpath: filepath.Join(root, "default", "default-20060102150405.ext"),
			wantB:     []byte{},
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: "",
			},
			inS:       "out",
			inT:       testT,
			inFpath:   filepath.Join(root, "dummy"),
			inB:       []byte{},
			wantIsErr: false,
			wantFpath: filepath.Join(root, "default", "default-20060102150405.ext"),
			wantB:     []byte("out"),
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: "",
			},
			inS:       "out",
			inT:       testT,
			inFpath:   filepath.Join(root, "default", "default-20060102150405.ext"),
			inB:       []byte("in"),
			wantIsErr: false,
			wantFpath: filepath.Join(root, "default", "default-20060102150405.ext"),
			wantB:     []byte("out"),
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inS:       "",
			inT:       testT,
			inFpath:   filepath.Join(root, "dummy"),
			inB:       []byte{},
			wantIsErr: false,
			wantFpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			wantB:     []byte{},
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inS:       "",
			inT:       testT,
			inFpath:   filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			inB:       []byte("in"),
			wantIsErr: false,
			wantFpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			wantB:     []byte("in"),
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inS:       "out",
			inT:       testT,
			inFpath:   filepath.Join(root, "dummy"),
			inB:       []byte{},
			wantIsErr: false,
			wantFpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			wantB:     []byte("out"),
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inS:       "out",
			inT:       testT,
			inFpath:   filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			inB:       []byte("in"),
			wantIsErr: false,
			wantFpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			wantB:     []byte("inout"),
		},
	}

	for i, c := range cases {
		err := os.Mkdir(root, 0o777)
		if err != nil {
			t.Fatal(err)
		}

		err = os.MkdirAll(filepath.Dir(c.inFpath), 0o777)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(c.inFpath, c.inB, 0o666)
		if err != nil {
			t.Fatal(err)
		}

		gotErr := c.inM.ExecWrite(c.inS, c.inT)
		gotIsErr := gotErr != nil

		if gotIsErr != c.wantIsErr {
			t.Errorf("case %d: err: expected %t, got %t", i, c.wantIsErr, gotIsErr)
		}
		if c.inM.fpath != c.wantFpath {
			t.Errorf(`case %d: fpath: expected "%s", got "%s"`, i, c.wantFpath, c.inM.fpath)
		}

		gotB, err := os.ReadFile(c.wantFpath)
		if err != nil {
			t.Errorf("case %d: file: %v", i, err)
		}
		if err == nil && !bytes.Equal(gotB, c.wantB) {
			t.Errorf(`case %d: file: expected "%s", got "%s"`, i, string(c.wantB), string(gotB))
		}

		err = os.RemoveAll(root)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestMachineWrite(t *testing.T) {
	tmp := t.TempDir()
	root := filepath.Join(tmp, "root")
	testT := time.Date(2006, time.January, 2, 15, 4, 5, 999999999, time.FixedZone("UTC-7", -7*60*60))

	cases := []struct {
		inM       *Machine
		inTitle   string
		inS       string
		inT       time.Time
		inFpath   string
		inB       []byte
		wantIsErr bool
		wantFpath string
		wantB     []byte
	}{
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: "",
			},
			inTitle:   "",
			inS:       "",
			inT:       testT,
			inFpath:   filepath.Join(root, "dummy"),
			inB:       []byte{},
			wantIsErr: false,
			wantFpath: filepath.Join(root, "default", "default-20060102150405.ext"),
			wantB:     []byte{},
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: "",
			},
			inTitle:   "",
			inS:       "",
			inT:       testT,
			inFpath:   filepath.Join(root, "default", "default-20060102150405.ext"),
			inB:       []byte("in"),
			wantIsErr: false,
			wantFpath: filepath.Join(root, "default", "default-20060102150405.ext"),
			wantB:     []byte{},
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: "",
			},
			inTitle:   "",
			inS:       "out",
			inT:       testT,
			inFpath:   filepath.Join(root, "dummy"),
			inB:       []byte{},
			wantIsErr: false,
			wantFpath: filepath.Join(root, "default", "default-20060102150405.ext"),
			wantB:     []byte("out"),
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: "",
			},
			inTitle:   "",
			inS:       "out",
			inT:       testT,
			inFpath:   filepath.Join(root, "default", "default-20060102150405.ext"),
			inB:       []byte("in"),
			wantIsErr: false,
			wantFpath: filepath.Join(root, "default", "default-20060102150405.ext"),
			wantB:     []byte("out"),
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: "",
			},
			inTitle:   "title",
			inS:       "",
			inT:       testT,
			inFpath:   filepath.Join(root, "dummy"),
			inB:       []byte{},
			wantIsErr: false,
			wantFpath: filepath.Join(root, "title", "title-20060102150405.ext"),
			wantB:     []byte{},
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: "",
			},
			inTitle:   "title",
			inS:       "",
			inT:       testT,
			inFpath:   filepath.Join(root, "title", "title-20060102150405.ext"),
			inB:       []byte("in"),
			wantIsErr: false,
			wantFpath: filepath.Join(root, "title", "title-20060102150405.ext"),
			wantB:     []byte{},
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: "",
			},
			inTitle:   "title",
			inS:       "out",
			inT:       testT,
			inFpath:   filepath.Join(root, "dummy"),
			inB:       []byte{},
			wantIsErr: false,
			wantFpath: filepath.Join(root, "title", "title-20060102150405.ext"),
			wantB:     []byte("out"),
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: "",
			},
			inTitle:   "title",
			inS:       "out",
			inT:       testT,
			inFpath:   filepath.Join(root, "title", "title-20060102150405.ext"),
			inB:       []byte("in"),
			wantIsErr: false,
			wantFpath: filepath.Join(root, "title", "title-20060102150405.ext"),
			wantB:     []byte("out"),
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inTitle:   "",
			inS:       "",
			inT:       testT,
			inFpath:   filepath.Join(root, "dummy"),
			inB:       []byte{},
			wantIsErr: false,
			wantFpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			wantB:     []byte{},
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inTitle:   "",
			inS:       "",
			inT:       testT,
			inFpath:   filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			inB:       []byte("in"),
			wantIsErr: false,
			wantFpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			wantB:     []byte("in"),
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inTitle:   "",
			inS:       "out",
			inT:       testT,
			inFpath:   filepath.Join(root, "dummy"),
			inB:       []byte{},
			wantIsErr: false,
			wantFpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			wantB:     []byte("out"),
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inTitle:   "",
			inS:       "out",
			inT:       testT,
			inFpath:   filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			inB:       []byte("in"),
			wantIsErr: false,
			wantFpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			wantB:     []byte("inout"),
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inTitle:   "title",
			inS:       "",
			inT:       testT,
			inFpath:   filepath.Join(root, "dummy"),
			inB:       []byte{},
			wantIsErr: false,
			wantFpath: filepath.Join(root, "title", "title-20060102150405.ext"),
			wantB:     []byte{},
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inTitle:   "title",
			inS:       "",
			inT:       testT,
			inFpath:   filepath.Join(root, "title", "title-20060102150405.ext"),
			inB:       []byte("in"),
			wantIsErr: false,
			wantFpath: filepath.Join(root, "title", "title-20060102150405.ext"),
			wantB:     []byte{},
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inTitle:   "title",
			inS:       "out",
			inT:       testT,
			inFpath:   filepath.Join(root, "dummy"),
			inB:       []byte{},
			wantIsErr: false,
			wantFpath: filepath.Join(root, "title", "title-20060102150405.ext"),
			wantB:     []byte("out"),
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inTitle:   "title",
			inS:       "out",
			inT:       testT,
			inFpath:   filepath.Join(root, "title", "title-20060102150405.ext"),
			inB:       []byte("in"),
			wantIsErr: false,
			wantFpath: filepath.Join(root, "title", "title-20060102150405.ext"),
			wantB:     []byte("out"),
		},
		{
			inM: &Machine{
				cfg: MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inTitle:   "/",
			inS:       "in",
			inT:       testT,
			inFpath:   filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			inB:       []byte{},
			wantIsErr: true,
			wantFpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			wantB:     []byte{},
		},
	}

	for i, c := range cases {
		err := os.Mkdir(root, 0o777)
		if err != nil {
			t.Fatal(err)
		}

		err = os.MkdirAll(filepath.Dir(c.inFpath), 0o777)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(c.inFpath, c.inB, 0o666)
		if err != nil {
			t.Fatal(err)
		}

		gotErr := c.inM.write(c.inTitle, c.inS, c.inT)
		gotIsErr := gotErr != nil

		if gotIsErr != c.wantIsErr {
			t.Errorf("case %d: err: expected %t, got %t", i, c.wantIsErr, gotIsErr)
		}
		if c.inM.fpath != c.wantFpath {
			t.Errorf(`case %d: fpath: expected "%s", got "%s"`, i, c.wantFpath, c.inM.fpath)
		}

		gotB, err := os.ReadFile(c.wantFpath)
		if err != nil {
			t.Errorf("case %d: file: %v", i, err)
		}
		if err == nil && !bytes.Equal(gotB, c.wantB) {
			t.Errorf(`case %d: file: expected "%s", got "%s"`, i, string(c.wantB), string(gotB))
		}

		err = os.RemoveAll(root)
		if err != nil {
			t.Fatal(err)
		}
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
