package main

import (
	"fmt"
	"log"
)

type RouterService struct {
	M      map[byte]Service
	Logger *log.Logger
}

func (r *RouterService) Serve(req []byte) ([]byte, error) {
	if len(req) < 1 {
		return r.serveError("empty request")
	}

	sid := req[0]
	s, ok := r.M[sid]
	if sid == 0xC5 || !ok {
		return r.serveError(fmt.Sprintf("invalid sid 0x%X", sid))
	}

	sr, err := s.Serve(req[1:])
	if err != nil {
		return r.serveError(err.Error())
	}

	resp := make([]byte, len(sr)+1)
	resp[0] = sid
	copy(resp[1:], sr)
	return resp, nil
}

func (r *RouterService) serveError(s string) ([]byte, error) {
	if r.Logger != nil {
		r.Logger.Printf("error: %s", s)
	}
	return []byte("\xC5" + s), nil
}
