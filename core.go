// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/gocondor/core/auth"
	"github.com/gocondor/core/cache"
	"github.com/gocondor/core/database"
	"github.com/gocondor/core/jwt"
	"github.com/gocondor/core/middlewares"
	"github.com/gocondor/core/routing"
	"github.com/gocondor/core/sessions"
	"github.com/unrolled/secure"
)

// App struct
type App struct {
	Features *Features
}

// GORM is a const represents gorm variable name
const GORM = "gorm"

// CACHE a cache engine variable
const CACHE = "cache"

// logs file path
const logsFilePath = "logs/app.log"

// logs file
var logsFile *os.File

// sessions middleware
var sesMiddleware gin.HandlerFunc

// New initiates the app struct
func New() *App {
	return &App{
		Features: &Features{},
	}
}

// SetEnv sets environment varialbes
func (app *App) SetEnv(env map[string]string) {
	for key, val := range env {
		os.Setenv(strings.TrimSpace(key), strings.TrimSpace(val))
	}
}

//Bootstrap initiate app
func (app *App) Bootstrap() {
	//initiate middlewares engine varialbe
	middlewares.New()

	//initiate routing engine varialbe
	routing.New()
	routing.NewGroupsHolder()

	//initiate data base varialb
	if app.Features.Database == true {
		database.New()
	}

	// initiate the cache varialbe
	if app.Features.Cache == true {
		cache.New(app.Features.Cache)
	}

	// initiate sessions
	sesMiddleware = initSessions(app.Features.Sessions)

}

// Run execute the app
func (app *App) Run(portNumber string) {
	// fallback to port number to 80 if not set
	if portNumber == "" {
		portNumber = "80"
	}

	// Log to file
	logsFile, err := os.OpenFile(logsFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer logsFile.Close()
	gin.DefaultWriter = io.MultiWriter(logsFile, os.Stdout)

	//initiate gin engines
	httpGinEngine := gin.Default()
	httpsGinEngine := gin.Default()

	// use sessions
	if app.Features.Sessions == true {
		httpGinEngine.Use(sesMiddleware)
		httpsGinEngine.Use(sesMiddleware)
	}

	// init auth
	auth.New(sessions.Resolve(), jwt.Resolve())

	httpsOn, _ := strconv.ParseBool(os.Getenv("APP_HTTPS_ON"))
	redirectToHTTPS, _ := strconv.ParseBool(os.Getenv("APP_REDIRECT_HTTP_TO_HTTPS"))
	letsencryptOn, _ := strconv.ParseBool(os.Getenv("APP_HTTPS_USE_LETSENCRYPT"))

	if httpsOn {
		//serve the https
		certFile := os.Getenv("APP_HTTPS_CERT_FILE_PATH")
		keyFile := os.Getenv("APP_HTTPS_KEY_FILE_PATH")
		host := app.GetHTTPSHost() + ":443"
		httpsGinEngine = app.UseMiddlewares(middlewares.Resolve().GetMiddlewares(), httpsGinEngine)
		httpsGinEngine = app.RegisterRoutes(routing.Resolve().GetRoutes(), httpsGinEngine)
		// register the groups routes
		httpsGinEngine = app.RegisterRoutes(routing.ResolveGroupsHolder().GetGroupsRoutes(), httpsGinEngine)

		// use let's encrypt
		if letsencryptOn {
			go log.Fatal(autotls.Run(httpsGinEngine, app.GetHTTPSHost()))
			return
		}

		go httpsGinEngine.RunTLS(host, certFile, keyFile)
	}

	//redirect http to https
	if httpsOn && redirectToHTTPS {
		secureFunc := func() gin.HandlerFunc {
			return func(c *gin.Context) {
				secureMiddleware := secure.New(secure.Options{
					SSLRedirect: true,
					SSLHost:     app.GetHTTPSHost() + ":443",
				})
				err := secureMiddleware.Process(c.Writer, c.Request)
				if err != nil {
					return
				}
				c.Next()
			}
		}()
		redirectEngine := gin.New()
		redirectEngine.Use(secureFunc)
		host := fmt.Sprintf("%s:%s", app.GetHTTPHost(), portNumber)
		redirectEngine.Run(host)
	}

	//serve the http version
	httpGinEngine = app.UseMiddlewares(middlewares.Resolve().GetMiddlewares(), httpGinEngine)
	httpGinEngine = app.RegisterRoutes(routing.Resolve().GetRoutes(), httpGinEngine)
	// register the groups routes
	httpGinEngine = app.RegisterRoutes(routing.ResolveGroupsHolder().GetGroupsRoutes(), httpGinEngine)

	host := fmt.Sprintf("%s:%s", app.GetHTTPHost(), portNumber)
	httpGinEngine.Run(host)
}

// SetAppMode set the mode if the app (debug|test|release)
func (app *App) SetAppMode(mode string) {
	if mode == gin.ReleaseMode || mode == gin.TestMode || mode == gin.DebugMode {
		gin.SetMode(mode)
	} else {
		gin.SetMode(gin.TestMode)
	}
}

// IntegratePackages helps with attaching packages to gin context
func (app *App) IntegratePackages(handlerFuncs []gin.HandlerFunc, engine *gin.Engine) *gin.Engine {
	for _, pkgIntegration := range handlerFuncs {
		engine.Use(pkgIntegration)
	}

	return engine
}

//SetEnabledFeatures to control what features to turn on or off
func (app *App) SetEnabledFeatures(features *Features) {
	app.Features = features
}

// UseMiddlewares use middlewares by gin engine
func (app *App) UseMiddlewares(middlewares []gin.HandlerFunc, engine *gin.Engine) *gin.Engine {
	for _, middleware := range middlewares {
		engine.Use(middleware)
	}

	return engine
}

// RegisterRoutes register routes on gin engine
func (app *App) RegisterRoutes(routers []routing.Route, engine *gin.Engine) *gin.Engine {
	for _, route := range routers {
		switch route.Method {
		case "get":
			engine.GET(route.Path, route.Handlers...)
		case "post":
			engine.POST(route.Path, route.Handlers...)
		case "delete":
			engine.DELETE(route.Path, route.Handlers...)
		case "patch":
			engine.PATCH(route.Path, route.Handlers...)
		case "put":
			engine.PUT(route.Path, route.Handlers...)
		case "options":
			engine.OPTIONS(route.Path, route.Handlers...)
		case "head":
			engine.HEAD(route.Path, route.Handlers...)
		}
	}

	return engine
}

// GetHTTPSHost returns https host name
func (app *App) GetHTTPSHost() string {
	host := os.Getenv("APP_HTTPS_HOST")
	//if not set get http instead
	if host == "" {
		host = os.Getenv("APP_HTTP_HOST")
	}
	//if both not set use local host
	if host == "" {
		host = "localhost"
	}
	return host
}

// GetHTTPHost returns http host name
func (app *App) GetHTTPHost() string {
	host := os.Getenv("APP_HTTP_HOST")
	//if both not set use local host
	if host == "" {
		host = "localhost"
	}
	return host
}

// initiate sessions
func initSessions(sessionsFeatureFlag bool) gin.HandlerFunc {
	ses := sessions.New(sessionsFeatureFlag)
	d := os.Getenv("SESSION_DRIVER")
	switch d {
	case "redis":
		return ses.InitiateRedistore("mysecret", "mysession")
	case "cookie":
		return ses.InitiateCookieStore("mysecret", "mysession")
	case "memstore":
		return ses.InitiateMemstoreStore("mysecret", "mysession")
	default:
		return ses.InitiateMemstoreStore("mysecret", "mysession")
	}
}
