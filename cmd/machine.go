package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MachineConfig struct {
	Root  string
	Title string
	Ext   string
	Log   *log.Logger
}

func NewMachineConfig() *MachineConfig {
	return &MachineConfig{
		Root:  DefaultRoot,
		Title: DefaultTitle,
		Ext:   DefaultExt,
		Log:   log.Default(),
	}
}

func (cfg *MachineConfig) Validate() error {
	if err := ValidateRootTitleExt(cfg.Root, cfg.Title, cfg.Ext); err != nil {
		return err
	}
	if cfg.Log == nil {
		return errors.New("nil logger")
	}
	return nil
}

type Machine struct {
	cfg   *MachineConfig
	fpath string
}

func NewMachine(cfg *MachineConfig) (*Machine, error) {
	c := *cfg
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return &Machine{
		cfg:   &c,
		fpath: "",
	}, nil
}

func (m *Machine) Exec(cmd string, t time.Time) error {
	if len(cmd) <= 0 {
		return errors.New("empty command")
	}
	switch cmd[0] {
	case 'n':
		return m.ExecNew(cmd[1:], t)
	case 'w':
		return m.ExecWrite(cmd[1:], t)
	default:
		return fmt.Errorf("unknown command '%s'", cmd[0:1])
	}
}

func (m *Machine) ExecNew(title string, t time.Time) error {
	err := m.internalNew(title, t)
	if err != nil {
		return err
	}

	return m.internalWrite("", true)
}

func (m *Machine) ExecWrite(s string, t time.Time) error {
	trunc := false
	if m.fpath == "" {
		trunc = true

		err := m.internalNew("", t)
		if err != nil {
			return err
		}
	}
	return m.internalWrite(s, trunc)
}

func (m *Machine) internalNew(title string, t time.Time) error {
	if title == "" {
		title = m.cfg.Title
	}
	fpath, err := GenerateFilepathOld(m.cfg.Root, title, m.cfg.Ext, t)
	if err != nil {
		return err
	}
	m.cfg.Log.Printf(`writing to "%s"`, fpath)
	m.fpath = fpath
	return nil
}

func (m *Machine) internalWrite(s string, trunc bool) error {
	if m.fpath == "" {
		panic("m.fpath is empty")
	}
	return WriteFile(m.fpath, []byte(s), trunc)
}

func GenerateFilepathOld(root, title, ext string, t time.Time) (string, error) {
	err := ValidateRootTitleExt(root, title, ext)
	if err != nil {
		return "", err
	}

	fpath := filepath.Join(root, title, title+"-"+t.Format("20060102150405")+ext)
	return fpath, nil
}

func ValidateRootTitleExt(root, title, ext string) error {
	if root == "" {
		return errors.New("empty root")
	}
	if title == "" {
		return errors.New("empty title")
	}
	if i := IndexPathSeparator(title); i >= 0 {
		return fmt.Errorf("title has disallowed character '%s'", title[i:i+1])
	}
	if strings.ContainsRune(title, '.') {
		return errors.New("title has disallowed character '.'")
	}
	if i := IndexPathSeparator(ext); i >= 0 {
		return fmt.Errorf("file extension has disallowed character '%s'", ext[i:i+1])
	}
	if ext != "" && ext[0] != '.' {
		return errors.New("file extension does not start with '.'")
	}
	return nil
}

func IndexPathSeparator(s string) int {
	for i := range s {
		if os.IsPathSeparator(s[i]) {
			return i
		}
	}
	return -1
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
