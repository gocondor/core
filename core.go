// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"

	"github.com/gocondor/core/env"
	"github.com/gocondor/core/logger"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/acme/autocert"
)

var logsDriver logger.LogsDriver
var loggr *logger.Logger
var appC AppConfig
var requestC RequestConfig
var jwtC JWTConfig

type configContainer struct {
	App     AppConfig
	Request RequestConfig
}

type App struct {
	t           int // for trancking middlewares
	chain       *chain
	middlewares *Middlewares
	Config      *configContainer
}

var app *App

func New() *App {
	app = &App{
		chain:       &chain{},
		middlewares: NewMiddlewares(),
		Config: &configContainer{
			App:     appC,
			Request: requestC,
		},
	}
	return app
}

func ResolveApp() *App {
	return app
}

func (app *App) SetLogsDriver(d logger.LogsDriver) {
	logsDriver = d
}

func (app *App) Bootstrap() {
	NewRouter()
	// database.New()
	// cache.New(app.Features.Cache)
	newValidator()
	loggr = logger.NewLogger(logsDriver)
}

func (app *App) Run(router *httprouter.Router) {
	portNumber := os.Getenv("App_HTTP_PORT")
	if portNumber == "" {
		portNumber = "80"
	}
	router = app.RegisterRoutes(ResolveRouter().GetRoutes(), router)
	useHttpsStr := os.Getenv("App_USE_HTTPS")
	if useHttpsStr == "" {
		useHttpsStr = "false"
	}
	useHttps, _ := strconv.ParseBool(useHttpsStr)

	fmt.Printf("Welcome to GoCondor\n")
	if useHttps {
		fmt.Printf("Listening on https \nWaiting for requests...\n")
	} else {
		fmt.Printf("Listening on port %s\nWaiting for requests...\n", portNumber)
	}
	UseLetsEncryptStr := os.Getenv("App_USE_LETSENCRYPT")
	if UseLetsEncryptStr == "" {
		UseLetsEncryptStr = "false"
	}
	UseLetsEncrypt, _ := strconv.ParseBool(UseLetsEncryptStr)
	if useHttps && UseLetsEncrypt {
		m := &autocert.Manager{
			Cache:  autocert.DirCache("letsencrypt-certs-dir"),
			Prompt: autocert.AcceptTOS,
		}
		LetsEncryptEmail := os.Getenv("APP_LETSENCRYPT_EMAIL")
		if LetsEncryptEmail != "" {
			m.Email = LetsEncryptEmail
		}
		HttpsHosts := os.Getenv("App_HTTPS_HOSTS")
		if HttpsHosts != "" {
			m.HostPolicy = autocert.HostWhitelist(HttpsHosts)
		}
		log.Fatal(http.Serve(m.Listener(), router))
		return
	}
	if useHttps && !UseLetsEncrypt {
		wd, err := os.Getwd()
		if err != nil {
			panic("can not get the current working dir")
		}
		CertFile := os.Getenv("App_CERT_FILE_PATH")
		if CertFile == "" {
			CertFile = "ssl/server.crt"
		}
		KeyFile := os.Getenv("App_KEY_FILE_PATH")
		if KeyFile == "" {
			KeyFile = "ssl/server.key"
		}
		certFilePath := filepath.Join(wd, CertFile)
		KeyFilePath := filepath.Join(wd, KeyFile)
		log.Fatal(http.ListenAndServeTLS(":443", certFilePath, KeyFilePath, router))
		return
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", portNumber), router))
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
			logger:    loggr,
			Validator: newValidator(),
			JWT: newJWT(JWTOptions{
				SigningKey: jwtC.SecretKey,
				Lifetime:   jwtC.Lifetime,
			}),
		}
		ctx.prepare(ctx)
		rhs := app.revHandlers(hs)
		app.prepareChain(rhs)
		app.t = 0
		app.chain.execute(ctx)
		for _, header := range ctx.Response.getHeaders() {
			w.Header().Add(header.key, header.val)
		}
		logger.CloseLogsFile()
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
	w.WriteHeader(http.StatusNotFound)
	res := "{\"message\": \"Not Found\"}"
	loggr.Error("Not Found")
	loggr.Error(debug.Stack())
	w.Header().Add(CONTENT_TYPE, CONTENT_TYPE_JSON)
	w.Write([]byte(res))
}

func (n methodNotAllowed) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	res := "{\"message\": \"Method not allowed\"}"
	loggr.Error("Method not allowed")
	loggr.Error(debug.Stack())
	w.Header().Add(CONTENT_TYPE, CONTENT_TYPE_JSON)
	w.Write([]byte(res))
}

var panicHandler = func(w http.ResponseWriter, r *http.Request, e interface{}) {
	shrtMsg := fmt.Sprintf("%v", e)
	loggr.Error(shrtMsg)
	fmt.Println(shrtMsg)
	loggr.Error(string(debug.Stack()))
	var res string
	if env.GetVarOtherwiseDefault("APP_ENV", "local") == PRODUCTION {
		res = "{\"message\": \"internal error\"}"
	} else {
		res = fmt.Sprintf("{\"message\": \"%v\", \"stack trace\": \"%v\"}", e, string(debug.Stack()))
	}
	w.WriteHeader(http.StatusInternalServerError)
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

func (app *App) SetAppConfig(a AppConfig) {
	appC = a
}

func (app *App) SetRequestConfig(r RequestConfig) {
	requestC = r
}

func (app *App) SetJWTConfig(j JWTConfig) {
	jwtC = j
}
