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

// 	server := httptest.NewServer(g)
// 	defer server.Close()
// 	_, err := http.Get(server.URL)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func TestRegisterRoutes(t *testing.T) {
// 	routes := []routing.Route{
// 		{
// 			Method: "get",
// 			Path:   "/:name",
// 			Handlers: []gin.HandlerFunc{
// 				func(c *gin.Context) {
// 					val, _ := c.Params.Get("name")
// 					c.JSON(http.StatusOK, gin.H{
// 						"name": val,
// 					})
// 				},
// 			},
// 		},
// 	}
// 	g := gin.New()
// 	app := New()
// 	app.RegisterRoutes(routes, g)
// 	s := httptest.NewServer(g)
// 	defer s.Close()

// 	res, _ := http.Get(fmt.Sprintf("%s/jack", s.URL))
// 	body, _ := ioutil.ReadAll(res.Body)
// 	type ResultStruct struct {
// 		Name string `json:"Name"`
// 	}
// 	var result ResultStruct
// 	json.Unmarshal(body, &result)
// 	if result.Name != "jack" {
// 		t.Errorf("failed assert execution of registered route")
// 	}
// }

// func TestSetEnabledFeatures(t *testing.T) {
// 	app := New()
// 	var Features *Features = &Features{
// 		Database: false,
// 		Cache:    false,
// 		GRPC:     false,
// 	}
// 	app.SetEnabledFeatures(Features)

// 	if app.Features.Database != false || app.Features.Cache != false || app.Features.GRPC != false {
// 		t.Errorf("failed setting features")
// 	}
// }

// func TestBootstrap(t *testing.T) {
// 	var Features *Features = &Features{
// 		Database: true,
// 		Cache:    false,
// 		GRPC:     false,
// 	}
// 	os.Setenv("DB_DRIVER", "sqlite") // set database driver to sqlite
// 	app := New()
// 	env, _ := godotenv.Read("./testingdata/.env")
// 	app.SetEnv(env)
// 	app.SetEnabledFeatures(Features)
// 	app.Bootstrap()

// 	m := middlewares.Resolve()
// 	if m == nil || fmt.Sprintf("%T", m) != "*middlewares.MiddlewaresUtil" {
// 		t.Errorf("failed asserting the initiation of MiddlewaresUtil")
// 	}

// 	r := routing.Resolve()
// 	if r == nil || fmt.Sprintf("%T", r) != "*routing.Router" {
// 		t.Errorf("failed asserting the initiation of Router")
// 	}

// 	d := database.Resolve()
// 	if d == nil || fmt.Sprintf("%T", d) != "*gorm.DB" {
// 		t.Errorf("failed asserting the initiation of Database")
// 	}
// }

// func TestUseMiddleWares(t *testing.T) {
// 	middlewares := []gin.HandlerFunc{
// 		func(c *gin.Context) {
// 			c.Set("VAR1", "VAL1")
// 		},
// 		func(c *gin.Context) {
// 			c.Set("VAR2", "VAL2")
// 		},
// 	}

// 	g := gin.New()
// 	app := New()
// 	g = app.UseMiddlewares(middlewares, g)

// 	g.GET("/", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"VAR1": c.MustGet("VAR1"),
// 			"VAR2": c.MustGet("VAR2"),
// 		})
// 	})
// 	s := httptest.NewServer(g)
// 	defer s.Close()
// 	res, _ := s.Client().Get(s.URL)
// 	body, _ := ioutil.ReadAll(res.Body)

// 	type ResponseStruct struct {
// 		VAR1 string `json:"VAR1"`
// 		VAR2 string `json:"VAR2"`
// 	}

// 	var response ResponseStruct
// 	json.Unmarshal(body, &response)

// 	if response.VAR1 != "VAL1" || response.VAR2 != "VAL2" {
// 		t.Errorf("failed asserting middlewares registering")
// 	}
// }

// func TestGetHTTPSHost(t *testing.T) {
// 	app := New()
// 	host := app.GetHTTPSHost()
// 	if host != "localhost" {
// 		t.Errorf("failed getting https host")
// 	}

// 	os.Setenv("APP_HTTPS_HOST", "testserver.com")
// 	host = app.GetHTTPSHost()
// 	if host != "testserver.com" {
// 		t.Errorf("failed getting https host")
// 	}

// }

// func TestGetHTTPHost(t *testing.T) {
// 	app := New()
// 	host := app.GetHTTPHost()
// 	if host != "localhost" {
// 		t.Errorf("failed getting http host")
// 	}

// 	os.Setenv("APP_HTTP_HOST", "testserver.com")
// 	host = app.GetHTTPHost()
// 	if host != "testserver.com" {
// 		t.Errorf("failed getting http host")
// 	}

// }

func TestMakeHTTPHandlerFunc(t *testing.T) {
	app := New()
	tmpFile := filepath.Join(t.TempDir(), uuid.NewString())
	app.SetLogsFilePath(filepath.Join(t.TempDir(), uuid.NewString()))
	hs := []Handler{
		func(c *Context) {
			f, _ := os.Create(tmpFile)
			f.WriteString("DFT2V56H")
		},
	}
	h := app.makeHTTPRouterHandlerFunc(hs)
	r := httptest.NewRequest("GET", "http://localhost", nil)
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
	m := &methodNotAllowed{}
	r := httptest.NewRequest("GET", "http://localhost", nil)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)
	rsp := w.Result()
	if rsp.StatusCode != 405 {
		t.Errorf("failed testing method not allowed")
	}
}

func TestNotFoundHandler(t *testing.T) {
	n := &notFoundHandler{}
	r := httptest.NewRequest("GET", "http://localhost", nil)
	w := httptest.NewRecorder()
	n.ServeHTTP(w, r)
	rsp := w.Result()
	if rsp.StatusCode != 404 {
		t.Errorf("failed testing not found handler")
	}

}

func TestUseMiddleware(t *testing.T) {
	app := New()
	UseMiddleware(func(c *Context) { c.Logger.Info("Testing!") })
	if len(app.middlewares.GetMiddlewares()) != 1 {
		t.Errorf("failed testing use middleware")
	}
}

func TestChainReset(t *testing.T) {
	c := &chain{}
	c.nodes = []Handler{
		func(c *Context) { c.Logger.Info("Testing1!") }, func(c *Context) { c.Logger.Info("Testing2!") },
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
		func(c *Context) { c.Logger.Info("testing!") },
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
	UseMiddleware(func(c *Context) { c.Logger.Info("Testing!") })
	hs := []Handler{
		func(c *Context) { c.Logger.Info("testing1!") },
		func(c *Context) { c.Logger.Info("testing2!") },
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
			httpRequest:    httptest.NewRequest("GET", "http://localhost", nil),
			httpPathParams: nil,
		},
		Response: &Response{
			headers:        []header{},
			textBody:       "",
			jsonBody:       "",
			responseWriter: httptest.NewRecorder(),
		},
		Logger: NewLogger(lgsPath),
	}
}

func TestRevHAndlers(t *testing.T) {
	app := New()
	t1 := func(c *Context) { c.Logger.Info("Testing1!") }
	t2 := func(c *Context) { c.Logger.Info("Testing2!") }

	handlers := []Handler{t1, t2}
	reved := app.revHandlers(handlers)
	if reflect.ValueOf(handlers[0]).Pointer() != reflect.ValueOf(reved[1]).Pointer() {
		t.Errorf("failed testing reverse handlers")
	}

	if reflect.ValueOf(handlers[1]).Pointer() != reflect.ValueOf(reved[0]).Pointer() {
		t.Errorf("failed testing reverse handlers")
	}
}
