package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestServiceHandlerServeHTTP(t *testing.T) {
	tests := []struct {
		inReq          *http.Request
		inSvcResp      []byte
		inSvcErr       error
		wantRespCode   int
		wantRespHeader map[string]string
		wantRespBody   []byte
		wantSvcQuery   map[string]string
	}{
		{
			inReq:        httptest.NewRequest("GET", "/path/to/service?k1=v1&k2=v2", nil),
			inSvcResp:    []byte("MOCK"),
			inSvcErr:     nil,
			wantRespCode: http.StatusOK,
			wantRespHeader: map[string]string{
				"Content-Type":           "text/plain; charset=utf-8",
				"X-Content-Type-Options": "nosniff",
				"Content-Length":         "9",
			},
			wantRespBody: []byte("SVCOKMOCK"),
			wantSvcQuery: map[string]string{
				"k1": "v1",
				"k2": "v2",
			},
		},
		{
			inReq:          httptest.NewRequest("POST", "/path/to/service?k1=v1&k2=v2", bytes.NewReader([]byte{})),
			inSvcResp:      []byte("MOCK"),
			inSvcErr:       nil,
			wantRespCode:   http.StatusMethodNotAllowed,
			wantRespHeader: map[string]string{},
			wantRespBody:   nil,
			wantSvcQuery:   map[string]string{},
		},
		{
			inReq:          httptest.NewRequest("GET", "/path/to/service?k1=v1&k2=v%", nil),
			inSvcResp:      []byte("MOCK"),
			inSvcErr:       nil,
			wantRespCode:   http.StatusBadRequest,
			wantRespHeader: map[string]string{},
			wantRespBody:   nil,
			wantSvcQuery:   map[string]string{},
		},
		{
			inReq:          httptest.NewRequest("GET", "/path/to/service?k1=v1&k2=v2", nil),
			inSvcResp:      nil,
			inSvcErr:       errors.New(""),
			wantRespCode:   http.StatusInternalServerError,
			wantRespHeader: map[string]string{},
			wantRespBody:   nil,
			wantSvcQuery: map[string]string{
				"k1": "v1",
				"k2": "v2",
			},
		},
	}

	for i, tt := range tests {
		s := &mockService{
			inResp: tt.inSvcResp,
			inErr:  tt.inSvcErr,
			gotV:   make(url.Values),
		}
		h := &ServiceHandler{
			S:        s,
			ErrorLog: log.New(io.Discard, "", 0),
		}
		w := httptest.NewRecorder()

		h.ServeHTTP(w, tt.inReq)
		gotResp := w.Result()
		gotRespBody, _ := io.ReadAll(gotResp.Body)

		if gotResp.StatusCode != tt.wantRespCode {
			t.Errorf("case %d: resp code: expected %d, got %d", i, tt.wantRespCode, gotResp.StatusCode)
		}
		for k, wantRespHeaderValue := range tt.wantRespHeader {
			gotRespHeaderValue := gotResp.Header.Get(k)
			if gotRespHeaderValue != wantRespHeaderValue {
				t.Errorf(`case %d: resp header: key "%s": expected "%s", got "%s"`, i, k, wantRespHeaderValue, gotRespHeaderValue)
			}
		}
		if tt.wantRespBody != nil && !bytes.Equal(gotRespBody, tt.wantRespBody) {
			t.Errorf("case %d: resp body: expected %#v, got %#v", i, tt.wantRespBody, gotRespBody)
		}
		for k, wantSvcQueryValue := range tt.wantSvcQuery {
			gotSvcQueryValue := s.gotV.Get(k)
			if gotSvcQueryValue != wantSvcQueryValue {
				t.Errorf(`case %d: svc query: key "%s": expected "%s", got "%s"`, i, k, wantSvcQueryValue, gotSvcQueryValue)
			}
		}
	}
}

func TestTimeServiceServeAPI(t *testing.T) {
	inNow := time.Date(2006, time.January, 2, 15, 4, 5, 999999999, time.FixedZone("UTC-7", -7*60*60))
	wantResp := []byte("20060102150405")

	inS := &TimeService{
		Now: func() time.Time { return inNow },
	}
	gotResp, gotErr := inS.ServeAPI(nil)

	if !(gotResp != nil && bytes.Equal(gotResp, wantResp)) {
		t.Errorf("resp: expected %#v, got %#v", wantResp, gotResp)
	}
	if gotErr != nil {
		t.Errorf("err: %v", gotErr)
	}
}

func TestWriteServiceServeAPINormal(t *testing.T) {
	tmp := t.TempDir()
	tmp = filepath.Join(tmp, "tmp")

	tests := []struct {
		inS             *WriteService
		inV             url.Values
		wantTmpFileName string
		wantTmpFileData []byte
	}{
		{
			inS: &WriteService{Root: tmp},
			inV: url.Values{
				"path": []string{"path/to/file"},
			},
			wantTmpFileName: filepath.Join(tmp, "path", "to", "file"),
			wantTmpFileData: []byte{},
		},
		{
			inS: &WriteService{Root: tmp},
			inV: url.Values{
				"path": []string{"path/to/file"},
				"data": []string{""},
			},
			wantTmpFileName: filepath.Join(tmp, "path", "to", "file"),
			wantTmpFileData: []byte{},
		},
		{
			inS: &WriteService{Root: tmp},
			inV: url.Values{
				"path": []string{"path/to/file"},
				"data": []string{"data"},
			},
			wantTmpFileName: filepath.Join(tmp, "path", "to", "file"),
			wantTmpFileData: []byte("data"),
		},
	}

	for i, tt := range tests {
		err := os.Mkdir(tmp, 0o777)
		if err != nil {
			t.Fatal(err)
		}

		gotResp, gotErr := tt.inS.ServeAPI(tt.inV)
		if gotResp == nil || len(gotResp) > 0 {
			t.Errorf("case %d: resp: expected %#v, got %#v", i, []byte{}, gotResp)
		}
		if gotErr != nil {
			t.Errorf("case %d: err: %v", i, gotErr)
		}

		gotTmpFileData, err := os.ReadFile(tt.wantTmpFileName)
		if err != nil {
			t.Errorf("case %d: file: %v", i, err)
		}
		if err == nil && !bytes.Equal(gotTmpFileData, tt.wantTmpFileData) {
			t.Errorf("case %d: file: expected %#v, got %#v", i, tt.wantTmpFileData, gotTmpFileData)
		}

		err = os.RemoveAll(tmp)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestWriteServiceServeAPIError(t *testing.T) {
	tmp := t.TempDir()

	tests := []struct {
		inS *WriteService
		inV url.Values
	}{
		{
			inS: &WriteService{Root: tmp},
			inV: url.Values{
				"path": []string{""},
				"data": []string{"data"},
			},
		},
		{
			inS: &WriteService{Root: tmp},
			inV: url.Values{
				"data": []string{"data"},
			},
		},
	}

	for i, tt := range tests {
		gotResp, gotErr := tt.inS.ServeAPI(tt.inV)

		if gotResp != nil {
			t.Errorf("case %d: resp: expected nil, got %#v", i, gotResp)
		}
		if gotErr == nil {
			t.Errorf("case %d: err: expected non-nil error, got nil", i)
		}
	}
}

func TestGenerateFilepath(t *testing.T) {
	tests := []struct {
		inRoot    string
		inPath    string
		wantFpath string
		wantIsErr bool
	}{
		{
			inRoot:    "root",
			inPath:    "path/to/file",
			wantFpath: strings.ReplaceAll("root/path/to/file", "/", string(os.PathSeparator)),
			wantIsErr: false,
		},
		{
			inRoot:    "root",
			inPath:    "",
			wantFpath: "",
			wantIsErr: true,
		},
	}

	for i, tt := range tests {
		gotFpath, gotErr := GenerateFilepath(tt.inRoot, tt.inPath)
		gotIsErr := gotErr != nil

		if gotFpath != tt.wantFpath {
			t.Errorf(`case %d: fpath: expected "%s", got "%s"`, i, tt.wantFpath, gotFpath)
		}
		if gotIsErr != tt.wantIsErr {
			t.Errorf(`case %d: err: expected %t, got %t`, i, tt.wantIsErr, gotIsErr)
		}
	}
}

func TestValidPath(t *testing.T) {
	testValidPath(t, []validPathTestCase{
		{
			inPath: ".",
			wantOK: true,
		},
		{
			inPath: "../",
			wantOK: false,
		},
		{
			inPath: "/",
			wantOK: false,
		},
	})
}

func TestValidPathNonWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.SkipNow()
	}

	testValidPath(t, []validPathTestCase{
		{
			inPath: `..\`,
			wantOK: true,
		},
		{
			inPath: `C:`,
			wantOK: true,
		},
	})
}

func TestValidPathWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.SkipNow()
	}

	testValidPath(t, []validPathTestCase{
		{
			inPath: `..\`,
			wantOK: false,
		},
		{
			inPath: `C:`,
			wantOK: false,
		},
	})
}

type validPathTestCase struct {
	inPath string
	wantOK bool
}

func testValidPath(t *testing.T, tests []validPathTestCase) {
	for i, tt := range tests {
		gotOK := ValidPath(tt.inPath)
		if gotOK != tt.wantOK {
			t.Errorf("case %d: expected %t, got %t", i, tt.wantOK, gotOK)
		}
	}
}

func TestWriteFile(t *testing.T) {
	tmp := t.TempDir()
	tmp = filepath.Join(tmp, "tmp")

	tests := []struct {
		inTmpDirName    string
		inTmpFileName   string
		inTmpFileData   []byte
		inName          string
		inData          []byte
		wantTmpFileName string
		wantTmpFileData []byte
	}{
		{
			inTmpDirName:    tmp,
			inTmpFileName:   filepath.Join(tmp, "dummy"),
			inTmpFileData:   []byte{},
			inName:          filepath.Join(tmp, "file"),
			inData:          []byte{},
			wantTmpFileName: filepath.Join(tmp, "file"),
			wantTmpFileData: []byte{},
		},
		{
			inTmpDirName:    tmp,
			inTmpFileName:   filepath.Join(tmp, "dummy"),
			inTmpFileData:   []byte{},
			inName:          filepath.Join(tmp, "file"),
			inData:          []byte("data"),
			wantTmpFileName: filepath.Join(tmp, "file"),
			wantTmpFileData: []byte("data"),
		},
		{
			inTmpDirName:    tmp,
			inTmpFileName:   filepath.Join(tmp, "dummy"),
			inTmpFileData:   []byte{},
			inName:          filepath.Join(tmp, "path", "to", "file"),
			inData:          []byte("data"),
			wantTmpFileName: filepath.Join(tmp, "path", "to", "file"),
			wantTmpFileData: []byte("data"),
		},
		{
			inTmpDirName:    tmp,
			inTmpFileName:   filepath.Join(tmp, "file"),
			inTmpFileData:   []byte("...."),
			inName:          filepath.Join(tmp, "file"),
			inData:          []byte("data"),
			wantTmpFileName: filepath.Join(tmp, "file"),
			wantTmpFileData: []byte("....data"),
		},
	}

	for i, tt := range tests {
		err := os.Mkdir(tmp, 0o777)
		if err != nil {
			t.Fatal(err)
		}

		err = os.MkdirAll(tt.inTmpDirName, 0o777)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(tt.inTmpFileName, tt.inTmpFileData, 0o666)
		if err != nil {
			t.Fatal(err)
		}

		gotErr := WriteFile(tt.inName, tt.inData)
		if gotErr != nil {
			t.Errorf("case %d: err: %v", i, gotErr)
		}

		gotTmpFileData, err := os.ReadFile(tt.wantTmpFileName)
		if err != nil {
			t.Errorf("case %d: file: %v", i, err)
		}
		if err == nil && !bytes.Equal(gotTmpFileData, tt.wantTmpFileData) {
			t.Errorf("case %d: file: expected %#v, got %#v", i, tt.wantTmpFileData, gotTmpFileData)
		}

		err = os.RemoveAll(tmp)
		if err != nil {
			t.Fatal(err)
		}
	}
}

type mockService struct {
	inResp []byte
	inErr  error
	gotV   url.Values
}

func (s *mockService) ServeAPI(v url.Values) ([]byte, error) {
	s.gotV = v
	return s.inResp, s.inErr
}
