package main

import (
	"errors"
	"time"
)

type TimeService struct {
	Now func() time.Time
}

func (s *TimeService) Serve(req []byte) ([]byte, error) {
	if req == nil {
		panic("req is nil")
	}
	if len(req) > 0 {
		return nil, errors.New("extra request data")
	}

	t := s.Now()
	resp := t.Format("20060102150405")
	return []byte(resp), nil
}
