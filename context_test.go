package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gocondor/core/logger"
	"github.com/google/uuid"
)

func TestDebugAny(t *testing.T) {
	r := httptest.NewRequest(GET, LOCALHOST, nil)
	w := httptest.NewRecorder()
	c := &Context{
		Request: &Request{
			HttpRequest:    r,
			httpPathParams: nil,
		},
		Response: &Response{
			headers:            []header{},
			textBody:           "",
			jsonBody:           []byte(""),
			HttpResponseWriter: w,
		},
		logger:    logger.NewLogger(&logger.LogNullDriver{}),
		Validator: nil,
		JWT:       nil,
	}
	h := func(c *Context) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			c.DebugAny("test-debug-msg")
		}
	}(c)
	h(w, r)
	b, err := io.ReadAll(w.Body)
	if err != nil {
		t.Errorf("failed testing debug any")
	}
	if !strings.Contains(string(b), "test-debug-msg") {
		t.Errorf("failed testing debug any")
	}
}

func TestLogInfo(t *testing.T) {
	tmpF := filepath.Join(t.TempDir(), uuid.NewString())
	r := httptest.NewRequest(GET, LOCALHOST, nil)
	w := httptest.NewRecorder()
	msg := "test-log-info"
	c := makeCTXLogTestCTX(t, w, r, tmpF)
	h := func(c *Context) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			c.LogInfo(msg)
		}
	}(c)
	h(w, r)
	fc, err := os.ReadFile(tmpF)
	if err != nil {
		t.Errorf("failed testing log info")
	}
	if !(strings.Contains(string(fc), msg) || strings.Contains(string(fc), "info:")) {
		t.Errorf("failed testing log info")
	}
}

func TestLogWarning(t *testing.T) {
	tmpF := filepath.Join(t.TempDir(), uuid.NewString())
	r := httptest.NewRequest(GET, LOCALHOST, nil)
	w := httptest.NewRecorder()
	msg := "test-log-warning"
	c := makeCTXLogTestCTX(t, w, r, tmpF)
	h := func(c *Context) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			c.LogWarning(msg)
		}
	}(c)
	h(w, r)
	fc, err := os.ReadFile(tmpF)
	if err != nil {
		t.Errorf("failed testing log warning")
	}
	if !(strings.Contains(string(fc), msg) || strings.Contains(string(fc), "warning:")) {
		t.Errorf("failed testing log warning")
	}
}

func TestLogDebug(t *testing.T) {
	tmpF := filepath.Join(t.TempDir(), uuid.NewString())
	r := httptest.NewRequest(GET, LOCALHOST, nil)
	w := httptest.NewRecorder()
	msg := "test-log-debug"
	c := makeCTXLogTestCTX(t, w, r, tmpF)
	h := func(c *Context) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			c.LogDebug(msg)
		}
	}(c)
	h(w, r)
	fc, err := os.ReadFile(tmpF)
	if err != nil {
		t.Errorf("failed testing log debug")
	}
	if !(strings.Contains(string(fc), msg) || strings.Contains(string(fc), "debug:")) {
		t.Errorf("failed testing log debug")
	}
}

func TestLogError(t *testing.T) {
	tmpF := filepath.Join(t.TempDir(), uuid.NewString())
	r := httptest.NewRequest(GET, LOCALHOST, nil)
	w := httptest.NewRecorder()
	msg := "test-log-error"
	c := makeCTXLogTestCTX(t, w, r, tmpF)
	h := func(c *Context) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			c.LogError(msg)
		}
	}(c)
	h(w, r)
	fc, err := os.ReadFile(tmpF)
	if err != nil {
		t.Errorf("failed testing log error")
	}
	if !(strings.Contains(string(fc), msg) || strings.Contains(string(fc), "error:")) {
		t.Errorf("failed testing log error")
	}
}

func makeCTXLogTestCTX(t *testing.T, w http.ResponseWriter, r *http.Request, tmpFilePath string) *Context {
	t.Helper()
	return &Context{
		Request: &Request{
			HttpRequest:    r,
			httpPathParams: nil,
		},
		Response: &Response{
			headers:            []header{},
			textBody:           "",
			jsonBody:           []byte(""),
			HttpResponseWriter: w,
		},
		logger:    logger.NewLogger(&logger.LogFileDriver{FilePath: tmpFilePath}),
		Validator: nil,
		JWT:       nil,
	}
}
