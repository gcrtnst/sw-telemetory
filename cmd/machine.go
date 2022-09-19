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
	Root     string
	Title    string
	Ext      string
	DirMode  os.FileMode
	FileMode os.FileMode
	Log      *log.Logger
}

func NewMachineConfig() MachineConfig {
	return MachineConfig{
		Root:     DefaultRoot,
		Title:    DefaultTitle,
		Ext:      DefaultExt,
		DirMode:  DefaultDirMode,
		FileMode: DefaultFileMode,
		Log:      log.New(os.Stderr, "", log.LstdFlags),
	}
}

type Machine struct {
	cfg   MachineConfig
	fpath string
}

func NewMachine(cfg MachineConfig) *Machine {
	return &Machine{
		cfg:   cfg,
		fpath: "",
	}
}

func (m *Machine) Exec(cmd string) error {
	if len(cmd) <= 0 {
		return errors.New("empty command")
	}
	switch cmd[0] {
	case 'n':
		return m.ExecNew(cmd[1:], time.Now())
	case 'w':
		return m.ExecWrite(cmd[1:], time.Now())
	default:
		return fmt.Errorf("unknown command '%s'", cmd[0:1])
	}
}

func (m *Machine) ExecNew(title string, t time.Time) error {
	err := m.internalNew(title, t)
	if err != nil {
		return err
	}

	return m.internalWrite("", os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
}

func (m *Machine) ExecWrite(s string, t time.Time) error {
	flag := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	if m.fpath == "" {
		flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC

		err := m.internalNew("", t)
		if err != nil {
			return err
		}
	}
	return m.internalWrite(s, flag)
}

func (m *Machine) internalNew(title string, t time.Time) error {
	if title == "" {
		title = m.cfg.Title
	}
	fpath, err := GenerateFilepath(m.cfg.Root, title, m.cfg.Ext, t)
	if err != nil {
		return err
	}
	m.cfg.Log.Printf(`writing to "%s"`, fpath)
	m.fpath = fpath
	return nil
}

func (m *Machine) internalWrite(s string, flag int) error {
	if m.fpath == "" {
		panic("m.fpath is empty")
	}
	return WriteFile(m.fpath, []byte(s), flag)
}

func GenerateFilepath(root, title, ext string, t time.Time) (string, error) {
	if root == "" {
		return "", errors.New("empty root")
	}
	if title == "" {
		return "", errors.New("empty title")
	}
	if i := IndexPathSeparator(title); i >= 0 {
		return "", fmt.Errorf("title has disallowed character '%s'", title[i:i+1])
	}
	if strings.ContainsRune(title, '.') {
		return "", errors.New("title has disallowed character '.'")
	}
	if i := IndexPathSeparator(ext); i >= 0 {
		return "", fmt.Errorf("file extension has disallowed character '%s'", ext[i:i+1])
	}
	if ext != "" && ext[0] != '.' {
		return "", errors.New("file extension does not start with '.'")
	}

	fpath := filepath.Join(root, title, title+"-"+t.Format("20060102150405")+ext)
	return fpath, nil
}

func IndexPathSeparator(s string) int {
	for i := range s {
		if os.IsPathSeparator(s[i]) {
			return i
		}
	}
	return -1
}

func WriteFile(name string, data []byte, flag int) (err error) {
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
