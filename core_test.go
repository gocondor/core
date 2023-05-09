// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gocondor/core/logger"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func TestNew(t *testing.T) {
	app := New()
	if fmt.Sprintf("%T", app) != "*core.App" {
		t.Errorf("failed testing new core")
	}
}

func TestSetEnv(t *testing.T) {
	env, err := godotenv.Read("./testingdata/.env")
	if err != nil {
		t.Errorf("failed reading .env file")
	}
	app := New()
	app.SetEnv(env)

	if os.Getenv("KEY_ONE") != "VAL_ONE" || os.Getenv("KEY_TWO") != "VAL_TWO" {
		t.Errorf("failed to set env vars")
	}
}

func TestMakeHTTPHandlerFunc(t *testing.T) {
	app := New()
	tmpFile := filepath.Join(t.TempDir(), uuid.NewString())
	app.SetLogsDriver(&logger.LogFileDriver{
		FilePath: filepath.Join(t.TempDir(), uuid.NewString()),
	})
	hs := []Handler{
		func(c *Context) {
			f, _ := os.Create(tmpFile)
			f.WriteString("DFT2V56H")
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

func TestMethodNotAllowedHandler(t *testing.T) {
	a := New()
	a.SetLogsDriver(&logger.LogNullDriver{})
	a.Bootstrap()
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
	app := New()
	UseMiddleware(func(c *Context) { c.LogInfo("Testing!") })
	if len(app.middlewares.GetMiddlewares()) != 1 {
		t.Errorf("failed testing use middleware")
	}
}

func TestChainReset(t *testing.T) {
	c := &chain{}
	c.nodes = []Handler{
		func(c *Context) { c.LogInfo("Testing1!") }, func(c *Context) { c.LogInfo("Testing2!") },
	}

	c.reset()
	if len(c.nodes) != 0 {
		t.Errorf("failed testing reset chain")
	}
}

func TestNext(t *testing.T) {
	app := New()
	app.t = 0
	tfPath := filepath.Join(t.TempDir(), uuid.NewString())
	hs := []Handler{
		func(c *Context) { c.Next() },
		func(c *Context) {
			f, _ := os.Create(tfPath)
			f.WriteString("DFT2V56H")
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
		func(c *Context) { c.LogInfo("testing!") },
		func(c *Context) {
			f, _ := os.Create(tf)
			f.WriteString("DFT2V56H")
		},
	}
	c.getByIndex(1)(makeCTX(t))
	d, _ := os.ReadFile(tf)
	if string(d) != "DFT2V56H" {
		t.Errorf("failed testing chain get by index")
	}
}

func TestPrepareChain(t *testing.T) {
	app := New()
	UseMiddleware(func(c *Context) { c.LogInfo("Testing!") })
	hs := []Handler{
		func(c *Context) { c.LogInfo("testing1!") },
		func(c *Context) { c.LogInfo("testing2!") },
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
		func(c *Context) {
			tf, _ := os.Create(f1Path)
			defer tf.Close()
			tf.WriteString("DFT2V56H")
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
			textBody:           "",
			jsonBody:           []byte(""),
			HttpResponseWriter: httptest.NewRecorder(),
		},
		logger: logger.NewLogger(&logger.LogFileDriver{
			FilePath: lgsPath,
		}),
	}
}

func TestRevHAndlers(t *testing.T) {
	app := New()
	t1 := func(c *Context) { c.LogInfo("Testing1!") }
	t2 := func(c *Context) { c.LogInfo("Testing2!") }

	handlers := []Handler{t1, t2}
	reved := app.revHandlers(handlers)
	if reflect.ValueOf(handlers[0]).Pointer() != reflect.ValueOf(reved[1]).Pointer() {
		t.Errorf("failed testing reverse handlers")
	}

	if reflect.ValueOf(handlers[1]).Pointer() != reflect.ValueOf(reved[0]).Pointer() {
		t.Errorf("failed testing reverse handlers")
	}
}
