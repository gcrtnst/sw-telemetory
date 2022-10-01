package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	port := flag.Int("port", 0, "listen port")
	root := flag.String("root", ".", "root directory for write command")
	flag.Parse()

	code := command(*port, *root)
	os.Exit(code)
}

func command(port int, root string) int {
	logger := log.Default()

	mux := http.NewServeMux()
	mux.Handle("/write", &ServiceHandler{
		S:        &WriteService{Root: root},
		ErrorLog: logger,
	})
	mux.Handle("/time", &ServiceHandler{
		S:        &TimeService{Now: time.Now},
		ErrorLog: logger,
	})

	s := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       75 * time.Second,
		ErrorLog:          logger,
	}

	lis, err := net.ListenTCP("tcp4", &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: port,
	})
	if err != nil {
		logger.Print(err)
		return 1
	}
	logger.Printf("listening on %s", lis.Addr().String())

	err = s.Serve(lis)
	logger.Print(err)
	return 1
}

type ServiceHandler struct {
	S        Service
	ErrorLog *log.Logger
}

func (h *ServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "" {
		h.ErrorLog.Printf(`"%s" from %s: method not allowed: %s`, r.URL.String(), r.RemoteAddr, r.Method)
		http.Error(w, "error", http.StatusMethodNotAllowed)
		return
	}

	u, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		h.ErrorLog.Printf(`"%s" from %s: invalid url query`, r.URL.String(), r.RemoteAddr)
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	b, err := h.S.ServeAPI(u)
	if err != nil {
		h.ErrorLog.Printf(`"%s" from %s: %s`, r.URL.String(), r.RemoteAddr, err.Error())
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	prefix := "SVCOK"
	resp := make([]byte, len(prefix)+len(b))
	copy(resp, prefix)
	copy(resp[len(prefix):], b)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Length", strconv.Itoa(len(resp)))
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

type Service interface {
	ServeAPI(url.Values) ([]byte, error)
}

type TimeService struct {
	Now func() time.Time
}

func (s *TimeService) ServeAPI(v url.Values) ([]byte, error) {
	t := s.Now()
	resp := t.Format("20060102150405")
	return []byte(resp), nil
}

type WriteService struct {
	Root string
}

func (s *WriteService) ServeAPI(v url.Values) ([]byte, error) {
	path := v.Get("path")
	data := v.Get("data")

	name, err := GenerateFilepath(s.Root, path)
	if err != nil {
		return nil, err
	}

	err = WriteFile(name, []byte(data))
	if err != nil {
		return nil, err
	}

	return []byte{}, nil
}

func GenerateFilepath(root, path string) (string, error) {
	if !ValidPath(path) {
		return "", fmt.Errorf(`invalid path: "%s"`, path)
	}

	path = filepath.FromSlash(path)
	return filepath.Join(root, path), nil
}

func ValidPath(path string) bool {
	return fs.ValidPath(path) && !(runtime.GOOS == "windows" && strings.ContainsAny(path, `:\`))
}

func WriteFile(name string, data []byte) (err error) {
	errMkdir := os.MkdirAll(filepath.Dir(name), 0o777)
	if errMkdir != nil {
		return errMkdir
	}

	f, errOpen := os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o666)
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
