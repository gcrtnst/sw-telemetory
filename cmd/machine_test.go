package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMachineConfigValidate(t *testing.T) {
	cases := []struct {
		inCfg     *MachineConfig
		wantIsErr bool
	}{
		{
			inCfg:     NewMachineConfig(),
			wantIsErr: false,
		},
		{
			inCfg: &MachineConfig{
				Root:  "root",
				Title: "title",
				Ext:   ".ext",
				Log:   log.Default(),
			},
			wantIsErr: false,
		},
		{
			inCfg: &MachineConfig{
				Root:  "",
				Title: "title",
				Ext:   ".ext",
				Log:   log.Default(),
			},
			wantIsErr: true,
		},
		{
			inCfg: &MachineConfig{
				Root:  "root",
				Title: "title.",
				Ext:   ".ext",
				Log:   log.Default(),
			},
			wantIsErr: true,
		},
		{
			inCfg: &MachineConfig{
				Root:  "root",
				Title: "title",
				Ext:   "ext",
				Log:   log.Default(),
			},
			wantIsErr: true,
		},
	}

	for i, c := range cases {
		gotErr := c.inCfg.Validate()
		gotIsErr := gotErr != nil

		if gotIsErr != c.wantIsErr {
			t.Errorf("case %d: err: expected %t, got %t", i, c.wantIsErr, gotIsErr)
		}
	}
}

func TestMachineExec(t *testing.T) {
	tmp := t.TempDir()
	root := filepath.Join(tmp, "root")
	testT := time.Date(2006, time.January, 2, 15, 4, 5, 999999999, time.FixedZone("UTC-7", -7*60*60))

	cases := []struct {
		inM         *Machine
		inCmd       string
		inT         time.Time
		wantIsErr   bool
		wantTmpName string
		wantTmpData []byte
	}{
		{
			inM: &Machine{
				cfg: &MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inCmd:       "n",
			inT:         testT,
			wantIsErr:   false,
			wantTmpName: filepath.Join(root, "default", "default-20060102150405.ext"),
			wantTmpData: []byte{},
		},
		{
			inM: &Machine{
				cfg: &MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inCmd:       "ntitle",
			inT:         testT,
			wantIsErr:   false,
			wantTmpName: filepath.Join(root, "title", "title-20060102150405.ext"),
			wantTmpData: []byte{},
		},
		{
			inM: &Machine{
				cfg: &MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inCmd:       "w",
			inT:         testT,
			wantIsErr:   false,
			wantTmpName: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			wantTmpData: []byte{},
		},
		{
			inM: &Machine{
				cfg: &MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inCmd:       "wout",
			inT:         testT,
			wantIsErr:   false,
			wantTmpName: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			wantTmpData: []byte("out"),
		},
		{
			inM: &Machine{
				cfg: &MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inCmd:     "",
			inT:       testT,
			wantIsErr: true,
		},
		{
			inM: &Machine{
				cfg: &MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inCmd:     "z",
			inT:       testT,
			wantIsErr: true,
		},
	}

	for i, c := range cases {
		err := os.Mkdir(root, 0o777)
		if err != nil {
			t.Fatal(err)
		}

		gotErr := c.inM.Exec(c.inCmd, c.inT)
		gotIsErr := gotErr != nil

		if gotIsErr != c.wantIsErr {
			t.Errorf("case %d: err: expected %t, got %t", i, c.wantIsErr, gotErr)
		}
		if !gotIsErr {
			var gotTmpData []byte
			gotTmpData, err = os.ReadFile(c.wantTmpName)
			if err != nil {
				t.Errorf("case %d: data: %v", i, err)
			}
			if !bytes.Equal(gotTmpData, c.wantTmpData) {
				t.Errorf(`case %d: data: expected "%s", got "%s"`, i, string(c.wantTmpData), string(gotTmpData))
			}
		}

		err = os.RemoveAll(root)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestMachineExecNew(t *testing.T) {
	tmp := t.TempDir()
	root := filepath.Join(tmp, "root")
	testT := time.Date(2006, time.January, 2, 15, 4, 5, 999999999, time.FixedZone("UTC-7", -7*60*60))

	cases := []struct {
		inM       *Machine
		inTitle   string
		inT       time.Time
		inFpath   string
		inB       []byte
		wantIsErr bool
		wantFpath string
		wantB     []byte
	}{
		{
			inM: &Machine{
				cfg: &MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: "",
			},
			inTitle:   "",
			inT:       testT,
			inFpath:   filepath.Join(root, "dummy"),
			inB:       []byte{},
			wantIsErr: false,
			wantFpath: filepath.Join(root, "default", "default-20060102150405.ext"),
			wantB:     []byte{},
		},
		{
			inM: &Machine{
				cfg: &MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: "",
			},
			inTitle:   "",
			inT:       testT,
			inFpath:   filepath.Join(root, "default", "default-20060102150405.ext"),
			inB:       []byte("in"),
			wantIsErr: false,
			wantFpath: filepath.Join(root, "default", "default-20060102150405.ext"),
			wantB:     []byte{},
		},
		{
			inM: &Machine{
				cfg: &MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: "",
			},
			inTitle:   "title",
			inT:       testT,
			inFpath:   filepath.Join(root, "dummy"),
			inB:       []byte{},
			wantIsErr: false,
			wantFpath: filepath.Join(root, "title", "title-20060102150405.ext"),
			wantB:     []byte{},
		},
		{
			inM: &Machine{
				cfg: &MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: "",
			},
			inTitle:   "title",
			inT:       testT,
			inFpath:   filepath.Join(root, "title", "title-20060102150405.ext"),
			inB:       []byte("in"),
			wantIsErr: false,
			wantFpath: filepath.Join(root, "title", "title-20060102150405.ext"),
			wantB:     []byte{},
		},
		{
			inM: &Machine{
				cfg: &MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inTitle:   "",
			inT:       testT,
			inFpath:   filepath.Join(root, "dummy"),
			inB:       []byte{},
			wantIsErr: false,
			wantFpath: filepath.Join(root, "default", "default-20060102150405.ext"),
			wantB:     []byte{},
		},
		{
			inM: &Machine{
				cfg: &MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inTitle:   "",
			inT:       testT,
			inFpath:   filepath.Join(root, "default", "default-20060102150405.ext"),
			inB:       []byte("in"),
			wantIsErr: false,
			wantFpath: filepath.Join(root, "default", "default-20060102150405.ext"),
			wantB:     []byte{},
		},
		{
			inM: &Machine{
				cfg: &MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inTitle:   "title",
			inT:       testT,
			inFpath:   filepath.Join(root, "dummy"),
			inB:       []byte{},
			wantIsErr: false,
			wantFpath: filepath.Join(root, "title", "title-20060102150405.ext"),
			wantB:     []byte{},
		},
		{
			inM: &Machine{
				cfg: &MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inTitle:   "title",
			inT:       testT,
			inFpath:   filepath.Join(root, "title", "title-20060102150405.ext"),
			inB:       []byte("in"),
			wantIsErr: false,
			wantFpath: filepath.Join(root, "title", "title-20060102150405.ext"),
			wantB:     []byte{},
		},
		{
			inM: &Machine{
				cfg: &MachineConfig{
					Root:  root,
					Title: "default",
					Ext:   ".ext",
					Log:   log.New(io.Discard, "", 0),
				},
				fpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			},
			inTitle:   "/",
			inT:       testT,
			inFpath:   filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			inB:       []byte("in"),
			wantIsErr: true,
			wantFpath: filepath.Join(root, "fpath", "fpath-20060102150405.ext"),
			wantB:     []byte("in"),
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

		gotErr := c.inM.ExecNew(c.inTitle, c.inT)
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
				cfg: &MachineConfig{
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
				cfg: &MachineConfig{
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
				cfg: &MachineConfig{
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
				cfg: &MachineConfig{
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
				cfg: &MachineConfig{
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
				cfg: &MachineConfig{
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
				cfg: &MachineConfig{
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
				cfg: &MachineConfig{
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
			inRoot:    "",
			inTitle:   "title",
			inExt:     ".ext",
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

func TestValidateRootTitleExt(t *testing.T) {
	sep := string(os.PathSeparator)

	cases := []struct {
		inRoot    string
		inTitle   string
		inExt     string
		wantIsErr bool
	}{
		{
			inRoot:    "root",
			inTitle:   "title",
			inExt:     ".ext",
			wantIsErr: false,
		},
		{
			inRoot:    "root",
			inTitle:   "title",
			inExt:     "",
			wantIsErr: false,
		},
		{
			inRoot:    "",
			inTitle:   "title",
			inExt:     ".ext",
			wantIsErr: true,
		},
		{
			inRoot:    "root",
			inTitle:   "",
			inExt:     ".ext",
			wantIsErr: true,
		},
		{
			inRoot:    "root",
			inTitle:   sep + "title",
			inExt:     ".ext",
			wantIsErr: true,
		},
		{
			inRoot:    "root",
			inTitle:   "title.ext",
			inExt:     ".ext",
			wantIsErr: true,
		},
		{
			inRoot:    "root",
			inTitle:   "title",
			inExt:     sep + "ext",
			wantIsErr: true,
		},
		{
			inRoot:    "root",
			inTitle:   "title",
			inExt:     "ext" + sep,
			wantIsErr: true,
		},
		{
			inRoot:    "root",
			inTitle:   "title",
			inExt:     "ext",
			wantIsErr: true,
		},
	}

	for i, c := range cases {
		gotErr := ValidateRootTitleExt(c.inRoot, c.inTitle, c.inExt)
		gotIsErr := gotErr != nil

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

func TestWriteFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "tmp")

	cases := []struct {
		inName      string
		inData      []byte
		inTrunc     bool
		inTmpName   string
		inTmpData   []byte
		wantIsErr   bool
		wantTmpName string
		wantTmpData []byte
	}{
		{
			inName:      filepath.Join(tmp, "name"),
			inData:      []byte("data"),
			inTrunc:     false,
			inTmpName:   filepath.Join(tmp, "dummy"),
			inTmpData:   []byte("tmp"),
			wantIsErr:   false,
			wantTmpName: filepath.Join(tmp, "name"),
			wantTmpData: []byte("data"),
		},
		{
			inName:      filepath.Join(tmp, "name"),
			inData:      []byte("data"),
			inTrunc:     false,
			inTmpName:   filepath.Join(tmp, "name"),
			inTmpData:   []byte("tmp"),
			wantIsErr:   false,
			wantTmpName: filepath.Join(tmp, "name"),
			wantTmpData: []byte("tmpdata"),
		},
		{
			inName:      filepath.Join(tmp, "name"),
			inData:      []byte("data"),
			inTrunc:     true,
			inTmpName:   filepath.Join(tmp, "dummy"),
			inTmpData:   []byte("tmp"),
			wantIsErr:   false,
			wantTmpName: filepath.Join(tmp, "name"),
			wantTmpData: []byte("data"),
		},
		{
			inName:      filepath.Join(tmp, "name"),
			inData:      []byte("data"),
			inTrunc:     true,
			inTmpName:   filepath.Join(tmp, "name"),
			inTmpData:   []byte("tmp"),
			wantIsErr:   false,
			wantTmpName: filepath.Join(tmp, "name"),
			wantTmpData: []byte("data"),
		},
		{
			inName:      filepath.Join(tmp, "name"),
			inData:      []byte(""),
			inTrunc:     false,
			inTmpName:   filepath.Join(tmp, "dummy"),
			inTmpData:   []byte("tmp"),
			wantIsErr:   false,
			wantTmpName: filepath.Join(tmp, "name"),
			wantTmpData: []byte(""),
		},
		{
			inName:      filepath.Join(tmp, "name"),
			inData:      []byte(""),
			inTrunc:     false,
			inTmpName:   filepath.Join(tmp, "name"),
			inTmpData:   []byte("tmp"),
			wantIsErr:   false,
			wantTmpName: filepath.Join(tmp, "name"),
			wantTmpData: []byte("tmp"),
		},
		{
			inName:      filepath.Join(tmp, "name"),
			inData:      []byte(""),
			inTrunc:     true,
			inTmpName:   filepath.Join(tmp, "dummy"),
			inTmpData:   []byte("tmp"),
			wantIsErr:   false,
			wantTmpName: filepath.Join(tmp, "name"),
			wantTmpData: []byte(""),
		},
		{
			inName:      filepath.Join(tmp, "name"),
			inData:      []byte(""),
			inTrunc:     true,
			inTmpName:   filepath.Join(tmp, "name"),
			inTmpData:   []byte("tmp"),
			wantIsErr:   false,
			wantTmpName: filepath.Join(tmp, "name"),
			wantTmpData: []byte(""),
		},
		{
			inName:      filepath.Join(tmp, "tmp", "tmp", "name"),
			inData:      []byte("data"),
			inTrunc:     false,
			inTmpName:   filepath.Join(tmp, "dummy"),
			inTmpData:   []byte("tmp"),
			wantIsErr:   false,
			wantTmpName: filepath.Join(tmp, "tmp", "tmp", "name"),
			wantTmpData: []byte("data"),
		},
	}

	for i, c := range cases {
		err := os.Mkdir(tmp, 0o777)
		if err != nil {
			t.Fatal(err)
		}

		err = os.MkdirAll(filepath.Dir(c.inTmpName), 0o777)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(c.inTmpName, c.inTmpData, 0o666)
		if err != nil {
			t.Fatal(err)
		}

		gotErr := WriteFile(c.inName, c.inData, c.inTrunc)
		gotIsErr := gotErr != nil

		if gotIsErr != c.wantIsErr {
			t.Errorf("case %d: err: expected %t, got %t", i, c.wantIsErr, gotIsErr)
		}

		gotTmpData, err := os.ReadFile(c.wantTmpName)
		if err != nil {
			t.Errorf("case %d: data: %v", i, err)
		}
		if !bytes.Equal(gotTmpData, c.wantTmpData) {
			t.Errorf(`case %d: data: expected "%s", got "%s"`, i, string(c.wantTmpData), string(gotTmpData))
		}

		err = os.RemoveAll(tmp)
		if err != nil {
			t.Fatal(err)
		}
	}
}
