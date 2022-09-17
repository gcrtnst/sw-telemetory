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
		return m.ExecNew(cmd[1:])
	case 'w':
		return m.ExecWrite(cmd[1:])
	default:
		return fmt.Errorf("unknown command '%s'", cmd[0:1])
	}
}

func (m *Machine) ExecNew(title string) error {
	if title == "" {
		title = m.cfg.Title
	}
	fpath, err := GenerateFilepath(m.cfg.Root, title, m.cfg.Ext, time.Now())
	if err != nil {
		return err
	}
	m.setFilepath(fpath)
	return nil
}

func (m *Machine) ExecWrite(s string) (err error) {
	if m.fpath == "" {
		fpath, err := GenerateFilepath(m.cfg.Root, m.cfg.Title, m.cfg.Ext, time.Now())
		if err != nil {
			return err
		}
		m.setFilepath(fpath)
	}

	err = os.MkdirAll(filepath.Dir(m.fpath), m.cfg.DirMode)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(m.fpath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, m.cfg.FileMode)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			f.Close()
			return
		}
		err = f.Close()
	}()
	_, err = f.Write([]byte(s + "\n"))
	return err
}

func (m *Machine) setFilepath(fpath string) {
	m.fpath = fpath
	m.cfg.Log.Printf("writing to \"%s\"", fpath)
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
