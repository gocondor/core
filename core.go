// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gocondor/core/cache"
	"github.com/gocondor/core/database"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/html"
)

var filePath string

var logger *Logger

// App struct
type App struct {
	Features    *Features
	t           int // for trancking middlewares
	chain       *chain
	middlewares *Middlewares
}

var app *App

func New() *App {
	app = &App{
		Features:    &Features{},
		chain:       &chain{},
		middlewares: NewMiddlewares(),
	}

	return app
}

func ResolveApp() *App {
	return app
}

func (app *App) SetEnv(env map[string]string) {
	for key, val := range env {
		os.Setenv(strings.TrimSpace(key), strings.TrimSpace(val))
	}
}

func (app *App) SetLogsFilePath(f string) {
	filePath = f
}

func (app *App) GetLogsFile() *os.File {
	return logsFile
}

func (app *App) Bootstrap() {
	NewRouter()
	if app.Features.Database == true {
		database.New()
	}
	if app.Features.Cache == true {
		cache.New(app.Features.Cache)
	}
	logger = NewLogger(filePath)
}

func (app *App) Run(portNumber string, router *httprouter.Router) {
	router = app.RegisterRoutes(ResolveRouter().GetRoutes(), router)
	ee, _ := strconv.ParseInt("0x1F985", 0, 64)
	fmt.Printf("Welcome to GoCondor %v \n", html.UnescapeString(strconv.FormatInt(ee, 10)))
	fmt.Printf("Listening on port %s\nWaiting for requests...\n", portNumber)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", portNumber), router))
}

func (app *App) SetEnabledFeatures(features *Features) {
	app.Features = features
}

func (app *App) RegisterRoutes(routes []Route, router *httprouter.Router) *httprouter.Router {
	router.PanicHandler = panicHandler
	router.NotFound = notFoundHandler{}
	router.MethodNotAllowed = methodNotAllowed{}
	for _, route := range routes {
		switch route.Method {
		case GET:
			router.GET(route.Path, app.makeHTTPRouterHandlerFunc(route.Handlers))
		case POST:
			router.POST(route.Path, app.makeHTTPRouterHandlerFunc(route.Handlers))
		case DELETE:
			router.DELETE(route.Path, app.makeHTTPRouterHandlerFunc(route.Handlers))
		case PATCH:
			router.PATCH(route.Path, app.makeHTTPRouterHandlerFunc(route.Handlers))
		case PUT:
			router.PUT(route.Path, app.makeHTTPRouterHandlerFunc(route.Handlers))
		case OPTIONS:
			router.OPTIONS(route.Path, app.makeHTTPRouterHandlerFunc(route.Handlers))
		case HEAD:
			router.HEAD(route.Path, app.makeHTTPRouterHandlerFunc(route.Handlers))
		}
	}

	return router
}

func (app *App) makeHTTPRouterHandlerFunc(hs []Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := &Context{
			Request: &Request{
				HttpRequest:    r,
				httpPathParams: ps,
			},
			Response: &Response{
				headers:            []header{},
				textBody:           "",
				jsonBody:           []byte(""),
				HttpResponseWriter: w,
			},
			logger: logger,
		}
		ctx.prepare(ctx)
		rhs := app.revHandlers(hs)
		app.prepareChain(rhs)
		app.t = 0
		app.chain.execute(ctx)

		for _, header := range ctx.Response.getHeaders() {
			w.Header().Add(header.key, header.val)
		}
		defer logsFile.Close()
		if ctx.Response.getTextBody() != "" {
			w.Header().Add(CONTENT_TYPE, CONTENT_TYPE_HTML)
			w.Write([]byte(ctx.Response.getTextBody()))
		}
		if string(ctx.Response.getJsonBody()) != "" {
			w.Header().Add(CONTENT_TYPE, CONTENT_TYPE_JSON)
			w.Write(ctx.Response.getJsonBody())
		}
		app.t = 0
		ctx.Response.reset()
		app.chain.reset()
	}
}

type notFoundHandler struct{}
type methodNotAllowed struct{}

func (n notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO handle debug mode flag
	w.WriteHeader(http.StatusNotFound)
	res := "{\"message\": \"Not Found\"}"
	logger.Error("Not Found")
	logger.Error(debug.Stack())
	w.Header().Add(CONTENT_TYPE, CONTENT_TYPE_JSON)
	w.Write([]byte(res))
}

func (n methodNotAllowed) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO handle debug mode flag
	w.WriteHeader(http.StatusMethodNotAllowed)
	res := "{\"message\": \"Method not allowed\"}"
	logger.Error("Method not allowed")
	logger.Error(debug.Stack())
	w.Header().Add(CONTENT_TYPE, CONTENT_TYPE_JSON)
	w.Write([]byte(res))
}

var panicHandler = func(w http.ResponseWriter, r *http.Request, e interface{}) {
	// TODO handle debug mode flag
	w.WriteHeader(http.StatusInternalServerError)
	res := fmt.Sprintf("{\"message\": \"[internal error]: %v\", \"stack trace\": \"%v\"}", e, string(debug.Stack()))
	logger.Error(fmt.Sprintf("[internal error]: %v", e))
	logger.Error(string(debug.Stack()))
	debug.PrintStack()
	w.Header().Add(CONTENT_TYPE, CONTENT_TYPE_JSON)
	w.Write([]byte(res))
}

func UseMiddleware(mw func(c *Context)) {
	ResolveMiddlewares().Attach(mw)
}

func (app *App) Next(c *Context) {
	app.t = app.t + 1
	n := app.chain.getByIndex(app.t)
	if n != nil {
		n(c)
	}
}

type chain struct {
	nodes []Handler
}

func (cn *chain) reset() {
	cn.nodes = []Handler{}
}

func (c *chain) getByIndex(i int) Handler {
	for k := range c.nodes {
		if k == i {
			return c.nodes[i]
		}
	}

	return nil
}

func (app *App) prepareChain(hs []Handler) {
	mw := app.middlewares.GetMiddlewares()
	mw = append(mw, hs...)
	app.chain.nodes = append(app.chain.nodes, mw...)
}

func (cn *chain) execute(ctx *Context) {
	if cn.getByIndex(0) != nil {
		cn.getByIndex(0)(ctx)
	}
}

func (app *App) revHandlers(hs []Handler) []Handler {
	var rev []Handler
	for i := range hs {
		rev = append(rev, hs[(len(hs)-1)-i])
	}
	return rev
}
