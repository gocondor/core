package core

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gocondor/core/logger"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
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
			body:               nil,
			HttpResponseWriter: w,
		},
		logger:       logger.NewLogger(&logger.LogNullDriver{}),
		GetValidator: nil,
		GetJWT:       nil,
	}
	h := func(c *Context) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			var msg interface{}
			msg = "test-debug-pointer"
			c.DebugAny(&msg)
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
	if !strings.Contains(string(b), "test-debug-pointer") {
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

func TestGetPathParams(t *testing.T) {
	r := httptest.NewRequest(GET, LOCALHOST, nil)
	w := httptest.NewRecorder()
	pathParams := httprouter.Params{
		{
			Key:   "param1",
			Value: "param1val",
		},
		{
			Key:   "param2",
			Value: "param2val",
		},
	}
	a := New()
	h := a.makeHTTPRouterHandlerFunc([]Handler{
		func(c *Context) *Response {
			rsp := fmt.Sprintf("param1: %v | param2: %v", c.GetPathParam("param1"), c.GetPathParam("param2"))
			return c.Response.Text(rsp)
		},
	})
	h(w, r, pathParams)
	b, err := io.ReadAll(w.Body)
	if err != nil {
		t.Log("failed testing get path params")
	}
	bStr := string(b)
	if !(strings.Contains(bStr, "param1val") || strings.Contains(bStr, "param2val")) {
		t.Errorf("failed testing get path params")
	}
}

func TestGetRequestParams(t *testing.T) {
	app := New()
	hr := httprouter.New()
	gcr := NewRouter()
	gcr.Post("/pt", func(c *Context) *Response {
		fmt.Fprintln(c.Response.HttpResponseWriter, c.GetRequestParam("param"))
		return nil
	})
	gcr.Get("/gt", func(c *Context) *Response {
		fmt.Fprintln(c.Response.HttpResponseWriter, c.GetRequestParam("param"))
		return nil
	})
	hr = app.RegisterRoutes(gcr.GetRoutes(), hr)
	s := httptest.NewServer(hr)
	defer s.Close()
	rsp, err := http.PostForm(s.URL+"/pt", url.Values{"param": {"paramValPost"}})
	if err != nil {
		t.Logf("failed test get request params")
	}

	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Logf("failed test get request params")
	}
	if strings.TrimSpace(string(b)) != "paramValPost" {
		t.Errorf("failed test get request params")
	}
	rsp.Body.Close()
	rsp, err = http.Get(s.URL + "/gt?param=paramValGet")
	b, err = io.ReadAll(rsp.Body)
	if err != nil {
		t.Errorf("failed test get request params")
	}
	if strings.TrimSpace(string(b)) != "paramValGet" {
		t.Errorf("failed test get request param")
	}
	rsp.Body.Close()
}

func TestRequestParamsExists(t *testing.T) {
	app := New()
	hr := httprouter.New()
	gcr := NewRouter()
	gcr.Post("/pt", func(c *Context) *Response {
		fmt.Fprintln(c.Response.HttpResponseWriter, c.RequestParamExists("param"))
		return nil
	})
	gcr.Get("/gt", func(c *Context) *Response {
		fmt.Fprintln(c.Response.HttpResponseWriter, c.RequestParamExists("param"))
		return nil
	})
	hr = app.RegisterRoutes(gcr.GetRoutes(), hr)
	s := httptest.NewServer(hr)
	defer s.Close()
	rsp, err := http.PostForm(s.URL+"/pt", url.Values{"param": {"paramValPost"}})
	if err != nil {
		t.Logf("failed test get request params")
	}

	b, err := io.ReadAll(rsp.Body)

	if err != nil {
		t.Logf("failed test get request params")
	}
	if strings.TrimSpace(string(b)) != "true" {
		t.Errorf("failed test get request params")
	}
	rsp.Body.Close()

	rsp, err = http.Get(s.URL + "/gt?param=paramValGet")
	b, err = io.ReadAll(rsp.Body)
	if err != nil {
		t.Errorf("failed test get request params")
	}
	if strings.TrimSpace(string(b)) != "true" {
		t.Errorf("failed test get request param")
	}
	rsp.Body.Close()
}

func TestGetHeader(t *testing.T) {
	app := New()
	hr := httprouter.New()
	gcr := NewRouter()
	gcr.Post("/pt", func(c *Context) *Response {
		fmt.Fprintln(c.Response.HttpResponseWriter, c.GetHeader("headerkey"))
		return nil
	})
	gcr.Get("/gt", func(c *Context) *Response {
		fmt.Fprintln(c.Response.HttpResponseWriter, c.GetHeader("headerkey"))
		return nil
	})
	hr = app.RegisterRoutes(gcr.GetRoutes(), hr)
	s := httptest.NewServer(hr)
	defer s.Close()
	clt := &http.Client{}
	req, err := http.NewRequest("POST", s.URL+"/pt", nil)
	if err != nil {
		t.Errorf("failed test get header")
	}
	req.Header.Add("headerkey", "headerPostVal")
	rsp, err := clt.Do(req)
	if err != nil {
		t.Logf("failed test get request params")
	}
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Logf("failed test get request params")
	}
	if strings.TrimSpace(string(b)) != "headerPostVal" {
		t.Errorf("failed test get request params")
	}
	req, err = http.NewRequest("GET", s.URL+"/gt", nil)
	if err != nil {
		t.Errorf("failed test get header")
	}
	req.Header.Add("headerkey", "headerGetVal")
	rsp, err = clt.Do(req)
	if err != nil {
		t.Logf("failed test get request params")
	}
	b, err = io.ReadAll(rsp.Body)
	if err != nil {
		t.Logf("failed test get request params")
	}
	if strings.TrimSpace(string(b)) != "headerGetVal" {
		t.Errorf("failed test get request params")
	}
}

func TestGetUploadedFile(t *testing.T) {
	app := New()
	hr := httprouter.New()
	gcr := NewRouter()
	gcr.Post("/pt", func(c *Context) *Response {
		uploadedFile := c.GetUploadedFile("myfile")
		rs := fmt.Sprintf("file name: %v | size: %v", uploadedFile.Name, uploadedFile.Size)
		fmt.Fprintln(c.Response.HttpResponseWriter, rs)
		return nil
	})

	hr = app.RegisterRoutes(gcr.GetRoutes(), hr)
	s := httptest.NewServer(hr)
	defer s.Close()
	wd, _ := os.Getwd()
	tfp := filepath.Join(wd, "testingdata/testdata.json")
	file, err := os.Open(tfp)
	if err != nil {
		t.Error("failed test get upload file")
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("myfile", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	clt := &http.Client{}
	req, err := http.NewRequest("POST", s.URL+"/pt", body)
	if err != nil {
		t.Errorf("failed test get uploaded file")
	}
	req.Header.Add(CONTENT_TYPE, writer.FormDataContentType())
	rsp, err := clt.Do(req)
	if err != nil {
		t.Logf("failed test get get uploaded file")
	}
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Logf("failed test get get uploaded file")
	}
	rsp.Body.Close()
	asrtFile, err := os.Stat(tfp)
	if err != nil {
		t.Logf("failed test get get uploaded file")
	}
	if !strings.Contains(string(b), fmt.Sprintf("size: %v", asrtFile.Size())) {
		t.Errorf("failed test get uploaded file")
	}
	if !strings.Contains(string(b), "testdata.json") {
		t.Errorf("failed test get uploaded file")
	}
}

func TestMoveFile(t *testing.T) {
	tmpDir := t.TempDir()
	c := makeCTX(t)
	c.MoveFile("./testingdata/dummy.md", tmpDir, "dummy.md")
	fi, err := os.Stat(filepath.Join(tmpDir, "dummy.md"))
	if err != nil {
		t.Errorf("failed test move file")
	}
	if fi.Name() != "dummy.md" {
		t.Errorf("failed test move file")
	}
	t.Cleanup(func() {
		c.MoveFile(filepath.Join(tmpDir, "dummy.md"), "./testingdata", "dummy.md")
		os.Remove(filepath.Join(tmpDir, "dummy.md"))
	})
}

func TestCastToString(t *testing.T) {
	c := makeCTX(t)
	s := c.CastToString(25)
	if fmt.Sprintf("%T", s) != "string" {
		t.Errorf("failed test cast to string")
	}
	s = c.CastToString(25.54)
	if fmt.Sprintf("%T", s) != "string" {
		t.Errorf("failed test cast to string")
	}
	var v interface{} = "434"
	s = c.CastToString(v)
	if fmt.Sprintf("%T", s) != "string" {
		t.Errorf("failed test cast to string")
	}
	var vs interface{} = "this is a string"
	s = c.CastToString(vs)
	if fmt.Sprintf("%T", s) != "string" {
		t.Errorf("failed test cast to string")
	}
}

func TestCastToInt(t *testing.T) {
	c := makeCTX(t)
	i := c.CastToInt(4)
	if !(i == 4 && fmt.Sprintf("%T", i) == "int") {
		t.Errorf("failed test cast to int")
	}
	ii := c.CastToInt(4.434)
	if !(ii == 4 && fmt.Sprintf("%T", i) == "int") {
		t.Errorf("failed test cast to int")
	}
	iii := c.CastToInt("4")
	if !(iii == 4 && fmt.Sprintf("%T", i) == "int") {
		t.Errorf("failed test cast to int")
	}
	iiii := c.CastToInt("4.434")
	if !(iiii == 4 && fmt.Sprintf("%T", i) == "int") {
		t.Errorf("failed test cast to int")
	}
	var iInterface interface{}
	iInterface = 4
	i = c.CastToInt(iInterface)
	if !(i == 4 && fmt.Sprintf("%T", i) == "int") {
		t.Errorf("failed test cast to int")
	}
	iInterface = 4.545
	ii = c.CastToInt(iInterface)
	if !(ii == 4 && fmt.Sprintf("%T", i) == "int") {
		t.Errorf("failed test cast to int")
	}
	iInterface = "4"
	iii = c.CastToInt(iInterface)
	if !(iii == 4 && fmt.Sprintf("%T", i) == "int") {
		t.Errorf("failed test cast to int")
	}
	iInterface = "4.434"
	iiii = c.CastToInt(iInterface)
	if !(iiii == 4 && fmt.Sprintf("%T", i) == "int") {
		t.Errorf("failed test cast to int")
	}
}

func TestCastToFloat(t *testing.T) {
	c := makeCTX(t)
	f := c.CastToFloat(4)
	if !(f == 4 && fmt.Sprintf("%T", f) == "float64") {
		t.Errorf("failed test cast to float")
	}
	var varf32 float32 = 4.434
	ff32 := c.CastToFloat(varf32)
	if !(ff32 == 4.434 && fmt.Sprintf("%T", ff32) == "float64") {
		t.Errorf("failed test cast to float")
	}
	var varf64 float64 = 4.434
	ff64 := c.CastToFloat(varf64)
	if !(ff64 == 4.434 && fmt.Sprintf("%T", ff64) == "float64") {
		t.Errorf("failed test cast to float")
	}

	fff := c.CastToFloat("4")
	if !(fff == 4 && fmt.Sprintf("%T", fff) == "float64") {
		t.Errorf("failed test cast to float")
	}
	ffff := c.CastToFloat("4.434")
	if !(ffff == 4.434 && fmt.Sprintf("%T", ffff) == "float64") {
		t.Errorf("failed test cast to float")
	}
	var iInterface interface{}
	iInterface = 4
	f = c.CastToFloat(iInterface)
	if !(f == 4 && fmt.Sprintf("%T", f) == "float64") {
		t.Errorf("failed test cast to float")
	}
	iInterface = 4.434
	iff := c.CastToFloat(iInterface)
	if !(iff == 4.434 && fmt.Sprintf("%T", iff) == "float64") {
		t.Errorf("failed test cast to float")
	}
	iInterface = "4"
	fff = c.CastToFloat(iInterface)
	if !(fff == 4 && fmt.Sprintf("%T", fff) == "float64") {
		t.Errorf("failed test cast to float")
	}
	iInterface = "4.434"
	ffff = c.CastToFloat(iInterface)
	if !(ffff == 4.434 && fmt.Sprintf("%T", ffff) == "float64") {
		t.Errorf("failed test cast to float")
	}
}

func TestGetBaseDirPath(t *testing.T) {
	c := makeCTX(t)
	p := c.GetBaseDirPath()
	pwd, err := os.Getwd()
	if err != nil {
		t.Errorf("failed test get base dir path")
	}
	if p != pwd {
		t.Errorf("failed test get base dir path")
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
			body:               nil,
			HttpResponseWriter: w,
		},
		logger:       logger.NewLogger(&logger.LogFileDriver{FilePath: tmpFilePath}),
		GetValidator: nil,
		GetJWT:       nil,
	}
}
