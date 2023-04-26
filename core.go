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
}

func New() *App {
	return &App{
		Features: &Features{},
	}
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
}

func (app *App) Run(portNumber string, router *httprouter.Router) {
	router = app.RegisterRoutes(ResolveRouter().GetRoutes(), router)
	ee, _ := strconv.ParseInt("0x1F985", 0, 64)
	fmt.Printf("Welcome to GoCondor %v \n", html.UnescapeString(string(ee)))
	fmt.Printf("Listening on port %s\nWaiting for requests...", portNumber)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", portNumber), router))
}

func (app *App) SetEnabledFeatures(features *Features) {
	app.Features = features
}

func (app *App) RegisterRoutes(routes []Route, router *httprouter.Router) *httprouter.Router {
	for _, route := range routes {
		switch route.Method {
		case "get":
			router.GET(route.Path, makeHTTPRouterHandlerFunc(route.Handler))
		case "post":
			router.POST(route.Path, makeHTTPRouterHandlerFunc(route.Handler))
		case "delete":
			router.DELETE(route.Path, makeHTTPRouterHandlerFunc(route.Handler))
		case "patch":
			router.PATCH(route.Path, makeHTTPRouterHandlerFunc(route.Handler))
		case "put":
			router.PUT(route.Path, makeHTTPRouterHandlerFunc(route.Handler))
		case "options":
			router.OPTIONS(route.Path, makeHTTPRouterHandlerFunc(route.Handler))
		case "head":
			router.HEAD(route.Path, makeHTTPRouterHandlerFunc(route.Handler))
		}
	}

	return router
}

func makeHTTPRouterHandlerFunc(h Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := &Context{
			Request: &Request{
				httpRequest:    r,
				httpPathParams: ps,
			},
			Response: &Response{
				headers:        []header{},
				responseWriter: w,
			},
			Logger: NewLogger(filePath),
		}
		h(ctx)
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
	}
}
