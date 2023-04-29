// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gocondor/core/cache"
	"github.com/gocondor/core/database"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/html"
)

var filePath string

// App struct
type App struct {
	Features *Features
	t        int // for trancking middlewares
	chain    *chain
}

var app *App

func New() *App {
	app = &App{
		Features: &Features{},
		chain:    &chain{},
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
	NewMiddlewares()
	if app.Features.Database == true {
		database.New()
	}
	if app.Features.Cache == true {
		cache.New(app.Features.Cache)
	}
}

func (app *App) Run(portNumber string, router *httprouter.Router) {
	router = app.RegisterRoutes(ResolveRouter().GetRoutes(), router)
	ee, _ := strconv.ParseInt("0x1F985", 0, 64)
	fmt.Printf("Welcome to GoCondor %v \n", html.UnescapeString(string(ee)))
	fmt.Printf("Listening on port %s\nWaiting for requests...\n", portNumber)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", portNumber), router))
}

func (app *App) SetEnabledFeatures(features *Features) {
	app.Features = features
}

func (app *App) RegisterRoutes(routes []Route, router *httprouter.Router) *httprouter.Router {
	router.NotFound = notFoundHandler{}
	router.MethodNotAllowed = methodNotAllowed{}
	for _, route := range routes {
		switch route.Method {
		case "get":
			router.GET(route.Path, app.makeHTTPRouterHandlerFunc(route.Handler))
		case "post":
			router.POST(route.Path, app.makeHTTPRouterHandlerFunc(route.Handler))
		case "delete":
			router.DELETE(route.Path, app.makeHTTPRouterHandlerFunc(route.Handler))
		case "patch":
			router.PATCH(route.Path, app.makeHTTPRouterHandlerFunc(route.Handler))
		case "put":
			router.PUT(route.Path, app.makeHTTPRouterHandlerFunc(route.Handler))
		case "options":
			router.OPTIONS(route.Path, app.makeHTTPRouterHandlerFunc(route.Handler))
		case "head":
			router.HEAD(route.Path, app.makeHTTPRouterHandlerFunc(route.Handler))
		}
	}

	return router
}

func (app *App) makeHTTPRouterHandlerFunc(h Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := &Context{
			Request: &Request{
				httpRequest:    r,
				httpPathParams: ps,
			},
			Response: &Response{
				headers:        []header{},
				textBody:       "",
				jsonBody:       "",
				responseWriter: w,
			},
			Logger: NewLogger(filePath),
		}

		app.prepareChain(h)
		app.chain.execute(ctx)

		for _, header := range ctx.Response.getHeaders() {
			w.Header().Add(header.key, header.val)
		}
		defer logsFile.Close()
		if ctx.Response.getTextBody() != "" {
			w.Header().Add("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(ctx.Response.getTextBody()))
		}
		if ctx.Response.getJsonBody() != "" {
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(ctx.Response.getJsonBody()))
		}
		ctx.Response.reset()
	}
}

type notFoundHandler struct{}
type methodNotAllowed struct{}

func (n notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("Not Found"))
}

func (n methodNotAllowed) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(405)
	w.Write([]byte("Method not allowed"))
}

func UseMiddleware(mw func(C *Context)) {
	ResolveMiddlewares().Attach(mw)
}

func (app *App) handleMiddlewares(ctx *Context) {
	// check the index first
	app.t = 0
	if ResolveMiddlewares().getByIndex(0) != nil {
		ResolveMiddlewares().GetMiddlewares()[0](ctx)
	}
}

func (app *App) Next(c *Context) {
	n := app.chain.getByIndex(app.t + 1)
	if n != nil {
		n(c)
	}
}

type chain struct {
	nodes []Handler
}

func (c *chain) getByIndex(i int) Handler {
	for key, _ := range c.nodes {
		if key == i {
			return c.nodes[i]
		}
	}

	return nil
}

func (app *App) prepareChain(h Handler) {
	mw := ResolveMiddlewares().GetMiddlewares()
	mw = append(mw, h)
	app.chain.nodes = append(app.chain.nodes, mw...)
}

func (c *chain) execute(ctx *Context) {
	app.t = 0
	if c.getByIndex(0) != nil {
		c.getByIndex(0)(ctx)
	}
}
