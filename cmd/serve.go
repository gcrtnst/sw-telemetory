package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"
)

type ServerConfig struct {
	Root     string
	Title    string
	Ext      string
	DirMode  os.FileMode
	FileMode os.FileMode
	Log      *log.Logger
}

func NewServerConfig() ServerConfig {
	return ServerConfig{
		Root:     DefaultRoot,
		Title:    DefaultTitle,
		Ext:      DefaultExt,
		DirMode:  DefaultDirMode,
		FileMode: DefaultFileMode,
		Log:      log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (cfg ServerConfig) MachineConfig() MachineConfig {
	return MachineConfig{
		Root:     cfg.Root,
		Title:    cfg.Title,
		Ext:      cfg.Ext,
		DirMode:  cfg.DirMode,
		FileMode: cfg.FileMode,
		Log:      cfg.Log,
	}
}

type Server struct {
	cfg ServerConfig
	cg  CloseGroup
}

func NewServer(cfg ServerConfig) *Server {
	return &Server{
		cfg: cfg,
		cg:  CloseGroup{},
	}
}

func (s *Server) Close() error {
	s.cg.CloseAll()
	return nil
}

func (s *Server) ListenAndServe(port int) error {
	addr := &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: port,
	}
	lis, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(lis)
}

func (s *Server) Serve(lis net.Listener) (err error) {
	s.cfg.Log.Printf("listening on %s://%s", lis.Addr().Network(), lis.Addr().String())

	m := NewMachine(s.cfg.MachineConfig())
	rd := NewReader(lis)
	defer s.cg.Add(rd).CloseCatch(&err)
	sc := bufio.NewScanner(rd)
	sc.Split(bufio.ScanLines)

	for sc.Scan() {
		err = m.Exec(sc.Text(), time.Now())
		if err != nil {
			return
		}
	}
	return sc.Err()
}
