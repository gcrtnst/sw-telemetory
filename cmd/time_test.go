package main

import (
	"bytes"
	"testing"
	"time"
)

func TestTimeServiceServe(t *testing.T) {
	cases := []struct {
		inNow     time.Time
		inReq     []byte
		wantResp  []byte
		wantIsErr bool
	}{
		{
			inNow:     time.Date(2006, time.January, 2, 15, 4, 5, 999999999, time.FixedZone("UTC-7", -7*60*60)),
			inReq:     []byte{},
			wantResp:  []byte("20060102150405"),
			wantIsErr: false,
		},
		{
			inNow:     time.Date(2006, time.January, 2, 15, 4, 5, 999999999, time.FixedZone("UTC-7", -7*60*60)),
			inReq:     []byte("T"),
			wantResp:  nil,
			wantIsErr: true,
		},
	}

	for i, c := range cases {
		inSrv := &TimeService{
			Now: func() time.Time { return c.inNow },
		}

		gotResp, gotErr := inSrv.Serve(c.inReq)
		gotIsErr := gotErr != nil

		if ((gotResp == nil) != (c.wantResp == nil)) || !bytes.Equal(gotResp, c.wantResp) {
			t.Errorf("case %d: resp: expected %#v, got %#v", i, c.wantResp, gotResp)
		}
		if gotIsErr != c.wantIsErr {
			t.Errorf("case %d: err: expected %t, got %t", i, c.wantIsErr, gotIsErr)
		}
	}
}
