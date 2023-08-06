// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gocondor/core/env"
	"github.com/gocondor/core/logger"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func TestNew(t *testing.T) {
	app := createNewApp(t)
	if fmt.Sprintf("%T", app) != "*core.App" {
		t.Errorf("failed testing new core")
	}
}

func TestSetEnv(t *testing.T) {
	envVars, err := godotenv.Read("./testingdata/.env")
	if err != nil {
		t.Errorf("failed reading .env file")
	}

	env.SetEnvVars(envVars)

	if os.Getenv("KEY_ONE") != "VAL_ONE" || os.Getenv("KEY_TWO") != "VAL_TWO" {
		t.Errorf("failed to set env vars")
	}
}

func TestMakeHTTPHandlerFunc(t *testing.T) {
	app := createNewApp(t)
	tmpFile := filepath.Join(t.TempDir(), uuid.NewString())
	app.SetLogsDriver(&logger.LogFileDriver{
		FilePath: filepath.Join(t.TempDir(), uuid.NewString()),
	})
	hs := []Handler{
		func(c *Context) *Response {
			f, _ := os.Create(tmpFile)
			f.WriteString("DFT2V56H")
			c.Response.SetHeader("header-key", "header-val")
			return c.Response.Text("DFT2V56H")
		},
	}
	h := app.makeHTTPRouterHandlerFunc(hs)
	r := httptest.NewRequest(GET, LOCALHOST, nil)
	w := httptest.NewRecorder()
	h(w, r, []httprouter.Param{{Key: "tkey", Value: "tvalue"}})
	rsp := w.Result()
	if rsp.StatusCode != http.StatusOK {
		t.Errorf("failed testing make http handler func")
	}
	s, _ := os.ReadFile(tmpFile)
	if string(s) != "DFT2V56H" {
		t.Errorf("failed testing make http handler func")
	}
}

func TestMakeHTTPHandlerFuncVerifyJson(t *testing.T) {
	app := createNewApp(t)
	tmpFile := filepath.Join(t.TempDir(), uuid.NewString())
	app.SetLogsDriver(&logger.LogFileDriver{
		FilePath: filepath.Join(t.TempDir(), uuid.NewString()),
	})
	hs := []Handler{
		func(c *Context) *Response {
			f, _ := os.Create(tmpFile)
			f.WriteString("DFT2V56H")
			c.Response.SetHeader("header-key", "header-val")
			return c.Response.Json("{\"testKey\": \"testVal\"}")
		},
	}
	h := app.makeHTTPRouterHandlerFunc(hs)
	r := httptest.NewRequest(GET, LOCALHOST, nil)
	w := httptest.NewRecorder()
	h(w, r, []httprouter.Param{{Key: "tkey", Value: "tvalue"}})
	rsp := w.Result()
	if rsp.StatusCode != http.StatusOK {
		t.Errorf("failed testing make http handler func with json verify")
	}
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Errorf("failed testing make http handler func with json verify")
	}
	var j map[string]interface{}
	err = json.Unmarshal(b, &j)
	if err != nil {
		t.Errorf("failed testing make http handler func with json verify: %v", err)
	}
	if j["testKey"] != "testVal" {
		t.Errorf("failed testing make http handler func with json verify")
	}
}

func TestMethodNotAllowedHandler(t *testing.T) {
	app := createNewApp(t)
	app.SetLogsDriver(&logger.LogNullDriver{})
	app.Bootstrap()
	m := &methodNotAllowed{}
	r := httptest.NewRequest(GET, LOCALHOST, nil)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)
	rsp := w.Result()
	if rsp.StatusCode != 405 {
		t.Errorf("failed testing method not allowed")
	}
}

func TestNotFoundHandler(t *testing.T) {
	n := &notFoundHandler{}
	r := httptest.NewRequest(GET, LOCALHOST, nil)
	w := httptest.NewRecorder()
	n.ServeHTTP(w, r)
	rsp := w.Result()
	if rsp.StatusCode != 404 {
		t.Errorf("failed testing not found handler")
	}

}

func TestUseMiddleware(t *testing.T) {
	app := createNewApp(t)
	UseMiddleware(func(c *Context) *Response { c.LogInfo("Testing!"); return nil })
	if len(app.middlewares.GetMiddlewares()) != 1 {
		t.Errorf("failed testing use middleware")
	}
}

func TestChainReset(t *testing.T) {
	c := &chain{}
	c.nodes = []Handler{
		func(c *Context) *Response { c.LogInfo("Testing1!"); return nil },
		func(c *Context) *Response { c.LogInfo("Testing2!"); return nil },
	}

	c.reset()
	if len(c.nodes) != 0 {
		t.Errorf("failed testing reset chain")
	}
}

func TestNext(t *testing.T) {
	app := createNewApp(t)
	app.t = 0
	tfPath := filepath.Join(t.TempDir(), uuid.NewString())
	hs := []Handler{
		func(c *Context) *Response { c.Next(); return nil },
		func(c *Context) *Response {
			f, _ := os.Create(tfPath)
			f.WriteString("DFT2V56H")
			return nil
		},
	}
	app.prepareChain(hs)
	app.chain.execute(makeCTX(t))
	cnt, _ := os.ReadFile(tfPath)
	if string(cnt) != "DFT2V56H" {
		t.Errorf("failed testing next")
	}
}

func TestChainGetByIndex(t *testing.T) {
	c := &chain{}
	tf := filepath.Join(t.TempDir(), uuid.NewString())
	c.nodes = []Handler{
		func(c *Context) *Response { c.LogInfo("testing!"); return nil },
		func(c *Context) *Response {
			f, _ := os.Create(tf)
			f.WriteString("DFT2V56H")
			return nil
		},
	}
	c.getByIndex(1)(makeCTX(t))
	d, _ := os.ReadFile(tf)
	if string(d) != "DFT2V56H" {
		t.Errorf("failed testing chain get by index")
	}
}

func TestPrepareChain(t *testing.T) {
	app := createNewApp(t)
	UseMiddleware(func(c *Context) *Response { c.LogInfo("Testing!"); return nil })
	hs := []Handler{
		func(c *Context) *Response { c.LogInfo("testing1!"); return nil },
		func(c *Context) *Response { c.LogInfo("testing2!"); return nil },
	}
	app.prepareChain(hs)
	if len(app.chain.nodes) != 3 {
		t.Errorf("failed preparing chain")
	}
}

func TestChainExecute(t *testing.T) {
	tmpDir := t.TempDir()
	f1Path := filepath.Join(tmpDir, uuid.NewString())
	c := &chain{}
	c.nodes = []Handler{
		func(c *Context) *Response {
			tf, _ := os.Create(f1Path)
			defer tf.Close()
			tf.WriteString("DFT2V56H")
			return nil
		},
	}
	ctx := makeCTX(t)
	c.execute(ctx)
	cnt, _ := os.ReadFile(f1Path)
	if string(cnt) != "DFT2V56H" {
		t.Errorf("failed testing execute chain")
	}
}

func makeCTX(t *testing.T) *Context {
	t.Helper()
	lgsPath := filepath.Join(t.TempDir(), uuid.NewString())
	return &Context{
		Request: &Request{
			HttpRequest:    httptest.NewRequest(GET, LOCALHOST, nil),
			httpPathParams: nil,
		},
		Response: &Response{
			headers:            []header{},
			body:               nil,
			HttpResponseWriter: httptest.NewRecorder(),
		},
		logger: logger.NewLogger(&logger.LogFileDriver{
			FilePath: lgsPath,
		}),
		GetValidator: nil,
		GetJWT:       nil,
	}
}

func TestRevHAndlers(t *testing.T) {
	app := createNewApp(t)
	t1 := func(c *Context) *Response { c.LogInfo("Testing1!"); return nil }
	t2 := func(c *Context) *Response { c.LogInfo("Testing2!"); return nil }

	handlers := []Handler{t1, t2}
	reved := app.revHandlers(handlers)
	if reflect.ValueOf(handlers[0]).Pointer() != reflect.ValueOf(reved[1]).Pointer() {
		t.Errorf("failed testing reverse handlers")
	}

	if reflect.ValueOf(handlers[1]).Pointer() != reflect.ValueOf(reved[0]).Pointer() {
		t.Errorf("failed testing reverse handlers")
	}
}

func TestRegisterGetRoute(t *testing.T) {
	app := New()
	hr := httprouter.New()
	gcr := NewRouter()
	gcr.Get("/", func(c *Context) *Response {
		fmt.Fprintln(c.Response.HttpResponseWriter, c.GetRequestParam("param"))
		return nil
	})
	gcr.Post("/", func(c *Context) *Response {
		fmt.Fprintln(c.Response.HttpResponseWriter, c.GetRequestParam("param"))
		return nil
	})
	gcr.Delete("/", func(c *Context) *Response {
		fmt.Fprintln(c.Response.HttpResponseWriter, c.GetRequestParam("param"))
		return nil
	})
	gcr.Patch("/", func(c *Context) *Response {
		fmt.Fprintln(c.Response.HttpResponseWriter, c.GetRequestParam("param"))
		return nil
	})
	gcr.Put("/", func(c *Context) *Response {
		fmt.Fprintln(c.Response.HttpResponseWriter, c.GetRequestParam("param"))
		return nil
	})
	gcr.Options("/", func(c *Context) *Response {
		fmt.Fprintln(c.Response.HttpResponseWriter, c.GetRequestParam("param"))
		return nil
	})
	gcr.Head("/", func(c *Context) *Response {
		fmt.Fprintln(c.Response.HttpResponseWriter, c.GetRequestParam("param"))
		return nil
	})
	hr = app.RegisterRoutes(gcr.GetRoutes(), hr)
	s := httptest.NewServer(hr)
	defer s.Close()
	clt := &http.Client{}
	req, err := http.NewRequest("GET", s.URL+"?param=valget", nil)
	if err != nil {
		t.Errorf("failed test register routes")
	}
	rsp, err := clt.Do(req)
	if err != nil {
		t.Errorf("failed test register routes")
	}
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Errorf("failed test register routes")
	}
	if strings.TrimSpace(string(b)) != "valget" {
		t.Errorf("failed test register routes")
	}
	rsp.Body.Close()
}

func TestRegisterPostRoute(t *testing.T) {
	app := New()
	hr := httprouter.New()
	gcr := NewRouter()
	gcr.Post("/", func(c *Context) *Response {
		fmt.Fprintln(c.Response.HttpResponseWriter, c.GetRequestParam("param"))
		return nil
	})
	hr = app.RegisterRoutes(gcr.GetRoutes(), hr)
	s := httptest.NewServer(hr)
	defer s.Close()
	clt := &http.Client{}
	req, err := http.NewRequest("POST", s.URL+"?param=valpost", nil)
	if err != nil {
		t.Errorf("failed test register post route")
	}
	rsp, err := clt.Do(req)
	if err != nil {
		t.Errorf("failed test register post route")
	}
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Errorf("failed test register post route")
	}
	if strings.TrimSpace(string(b)) != "valpost" {
		t.Errorf("failed test register post route")
	}
	rsp.Body.Close()
}

func TestRegisterDeleteRoute(t *testing.T) {
	app := New()
	hr := httprouter.New()
	gcr := NewRouter()
	gcr.Delete("/", func(c *Context) *Response {
		fmt.Fprintln(c.Response.HttpResponseWriter, c.GetRequestParam("param"))
		return nil
	})
	hr = app.RegisterRoutes(gcr.GetRoutes(), hr)
	s := httptest.NewServer(hr)
	defer s.Close()
	clt := &http.Client{}
	req, err := http.NewRequest("DELETE", s.URL+"?param=valdelete", nil)
	if err != nil {
		t.Errorf("failed test register delete route")
	}
	rsp, err := clt.Do(req)
	if err != nil {
		t.Errorf("failed test register delete route")
	}
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Errorf("failed test register delete route")
	}
	if strings.TrimSpace(string(b)) != "valdelete" {
		t.Errorf("failed test register delete route")
	}
	rsp.Body.Close()
}

func TestRegisterPatchRoute(t *testing.T) {
	app := New()
	hr := httprouter.New()
	gcr := NewRouter()
	gcr.Patch("/", func(c *Context) *Response {
		fmt.Fprintln(c.Response.HttpResponseWriter, c.GetRequestParam("param"))
		return nil
	})
	hr = app.RegisterRoutes(gcr.GetRoutes(), hr)
	s := httptest.NewServer(hr)
	defer s.Close()
	clt := &http.Client{}
	req, err := http.NewRequest("PATCH", s.URL+"?param=valpatch", nil)
	if err != nil {
		t.Errorf("failed test register patch route")
	}
	rsp, err := clt.Do(req)
	if err != nil {
		t.Errorf("failed test register patch route")
	}
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Errorf("failed test register patch route")
	}
	if strings.TrimSpace(string(b)) != "valpatch" {
		t.Errorf("failed test register patch route")
	}
	rsp.Body.Close()
}

func TestRegisterPutRoute(t *testing.T) {
	app := New()
	hr := httprouter.New()
	gcr := NewRouter()
	gcr.Put("/", func(c *Context) *Response {
		fmt.Fprintln(c.Response.HttpResponseWriter, c.GetRequestParam("param"))
		return nil
	})
	hr = app.RegisterRoutes(gcr.GetRoutes(), hr)
	s := httptest.NewServer(hr)
	defer s.Close()
	clt := &http.Client{}
	req, err := http.NewRequest("PUT", s.URL+"?param=valput", nil)
	if err != nil {
		t.Errorf("failed test register put route")
	}
	rsp, err := clt.Do(req)
	if err != nil {
		t.Errorf("failed test register put route")
	}
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Errorf("failed test register put route")
	}
	if strings.TrimSpace(string(b)) != "valput" {
		t.Errorf("failed test register put route")
	}
	rsp.Body.Close()
}

func TestRegisterOptionsRoute(t *testing.T) {
	app := New()
	hr := httprouter.New()
	gcr := NewRouter()
	gcr.Options("/", func(c *Context) *Response {
		fmt.Fprintln(c.Response.HttpResponseWriter, c.GetRequestParam("param"))
		return nil
	})
	hr = app.RegisterRoutes(gcr.GetRoutes(), hr)
	s := httptest.NewServer(hr)
	defer s.Close()
	clt := &http.Client{}
	req, err := http.NewRequest("OPTIONS", s.URL+"?param=valoptions", nil)
	if err != nil {
		t.Errorf("failed test register options route")
	}
	rsp, err := clt.Do(req)
	if err != nil {
		t.Errorf("failed test register options route")
	}
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Errorf("failed test register options route")
	}
	if strings.TrimSpace(string(b)) != "valoptions" {
		t.Errorf("failed test register options route")
	}
	rsp.Body.Close()
}

func TestRegisterHeadRoute(t *testing.T) {
	app := New()
	hr := httprouter.New()
	gcr := NewRouter()
	tfp := filepath.Join(t.TempDir(), uuid.NewString())
	gcr.Head("/", func(c *Context) *Response {
		param := c.GetRequestParam("param")
		p, _ := param.(string)
		f, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR, 777)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer f.Close()
		f.WriteString("fromhead")
		return nil
	})
	hr = app.RegisterRoutes(gcr.GetRoutes(), hr)
	s := httptest.NewServer(hr)
	defer s.Close()
	clt := &http.Client{}
	req, err := http.NewRequest("HEAD", s.URL+"?param="+tfp, nil)
	if err != nil {
		t.Errorf("failed test register head route")
	}
	rsp, err := clt.Do(req)
	if err != nil {
		t.Errorf("failed test register head route")
	}
	f, err := os.Open(tfp)
	if err != nil {
		t.Errorf("failed test register head route: %v", err.Error())
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Errorf("failed test register head route: %v", err.Error())
	}
	f.Close()
	if strings.TrimSpace(string(b)) != "fromhead" {
		t.Errorf("failed test register head route")
	}
	rsp.Body.Close()
}

func TestPanicHandler(t *testing.T) {
	loggr = logger.NewLogger(&logger.LogNullDriver{})
	r := httptest.NewRequest(GET, LOCALHOST, nil)
	w := httptest.NewRecorder()
	panicHandler(w, r, "")
	rsp := w.Result()
	b, _ := io.ReadAll(rsp.Body)
	if !strings.Contains(string(b), "stack trace") {
		t.Errorf("failed test panic handler")
	}
}

func createNewApp(t *testing.T) *App {
	t.Helper()
	a := New()
	a.SetLogsDriver(&logger.LogNullDriver{})
	a.SetRequestConfig(testingRequestC)

	return a
}
