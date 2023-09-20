// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"syscall"

	"github.com/gocondor/core/env"
	"github.com/gocondor/core/logger"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/acme/autocert"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var loggr *logger.Logger
var logsDriver *logger.LogsDriver
var requestC RequestConfig
var jwtC JWTConfig
var gormC GormConfig
var cacheC CacheConfig
var db *gorm.DB
var mailer *Mailer
var basePath string
var disableEvents bool = false

type configContainer struct {
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
			Request: requestC,
		},
	}
	return app
}

func ResolveApp() *App {
	return app
}

func (app *App) SetLogsDriver(d logger.LogsDriver) {
	logsDriver = &d
}

func (app *App) Bootstrap() {
	loggr = logger.NewLogger(*logsDriver)
	NewRouter()
	NewEventsManager()
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
		CertFile := os.Getenv("App_CERT_FILE_PATH")
		if CertFile == "" {
			CertFile = "tls/server.crt"
		}
		KeyFile := os.Getenv("App_KEY_FILE_PATH")
		if KeyFile == "" {
			KeyFile = "tls/server.key"
		}
		certFilePath := filepath.Join(basePath, CertFile)
		KeyFilePath := filepath.Join(basePath, KeyFile)
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
			router.GET(route.Path, app.makeHTTPRouterHandlerFunc(route.Handler, route.Middlewares))
		case POST:
			router.POST(route.Path, app.makeHTTPRouterHandlerFunc(route.Handler, route.Middlewares))
		case DELETE:
			router.DELETE(route.Path, app.makeHTTPRouterHandlerFunc(route.Handler, route.Middlewares))
		case PATCH:
			router.PATCH(route.Path, app.makeHTTPRouterHandlerFunc(route.Handler, route.Middlewares))
		case PUT:
			router.PUT(route.Path, app.makeHTTPRouterHandlerFunc(route.Handler, route.Middlewares))
		case OPTIONS:
			router.OPTIONS(route.Path, app.makeHTTPRouterHandlerFunc(route.Handler, route.Middlewares))
		case HEAD:
			router.HEAD(route.Path, app.makeHTTPRouterHandlerFunc(route.Handler, route.Middlewares))
		}
	}
	return router
}

func (app *App) makeHTTPRouterHandlerFunc(h Handler, ms []Middleware) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := &Context{
			Request: &Request{
				HttpRequest:    r,
				httpPathParams: ps,
			},
			Response: &Response{
				headers:             []header{},
				body:                nil,
				contentType:         "",
				overrideContentType: "",
				HttpResponseWriter:  w,
				isTerminated:        false,
				redirectTo:          "",
			},
			GetValidator:     getValidator(),
			GetJWT:           getJWT(),
			GetGorm:          GetGormFunc(),
			GetCache:         resolveCache(),
			GetHashing:       resloveHashing(),
			GetMailer:        resolveMailer(),
			GetEventsManager: resolveEventsManager(),
			GetLogger:        resolveLogger(),
		}
		ctx.prepare(ctx)
		rhs := app.combHandlers(h, ms)
		app.prepareChain(rhs)
		app.t = 0
		app.chain.execute(ctx)
		for _, header := range ctx.Response.headers {
			w.Header().Add(header.key, header.val)
		}
		logger.CloseLogsFile()
		var ct string
		if ctx.Response.overrideContentType != "" {
			ct = ctx.Response.overrideContentType
		} else if ctx.Response.contentType != "" {
			ct = ctx.Response.contentType
		} else {
			ct = CONTENT_TYPE_HTML
		}
		w.Header().Add(CONTENT_TYPE, ct)
		if ctx.Response.statusCode != 0 {
			w.WriteHeader(ctx.Response.statusCode)
		}
		if ctx.Response.redirectTo != "" {
			http.Redirect(w, r, ctx.Response.redirectTo, http.StatusPermanentRedirect)
		} else {
			w.Write(ctx.Response.body)
		}
		e := ResolveEventsManager()
		if e != nil {
			e.setContext(ctx).processFiredEvents()
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
	isDebugModeStr := os.Getenv("APP_DEBUG_MODE")
	isDebugMode, err := strconv.ParseBool(isDebugModeStr)
	if err != nil {
		errStr := "error parsing env var APP_DEBUG_MODE"
		loggr.Error(errStr)
		fmt.Sprintln(errStr)
		w.Write([]byte(errStr))
		return
	}
	if !isDebugMode {
		errStr := "internal error"
		loggr.Error(errStr)
		fmt.Sprintln(errStr)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add(CONTENT_TYPE, CONTENT_TYPE_JSON)
		w.Write([]byte(fmt.Sprintf("{\"message\": \"%v\"}", errStr)))
		return
	}
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

func UseMiddleware(mw Middleware) {
	ResolveMiddlewares().Attach(mw)
}

func (app *App) Next(c *Context) {
	app.t = app.t + 1
	n := app.chain.getByIndex(app.t)
	if n != nil {
		f, ok := n.(Middleware)
		if ok {
			f(c)
		} else {
			ff, ok := n.(Handler)
			if ok {
				ff(c)
			}
		}
	}
}

type chain struct {
	nodes []interface{}
}

func (cn *chain) reset() {
	cn.nodes = []interface{}{}
}

func (c *chain) getByIndex(i int) interface{} {
	for k := range c.nodes {
		if k == i {
			return c.nodes[i]
		}
	}

	return nil
}

func (app *App) prepareChain(hs []interface{}) {
	mw := app.middlewares.GetMiddlewares()
	for _, v := range mw {
		app.chain.nodes = append(app.chain.nodes, v)
	}
	for _, v := range hs {
		app.chain.nodes = append(app.chain.nodes, v)
	}
}

func (cn *chain) execute(ctx *Context) {
	i := cn.getByIndex(0)
	if i != nil {
		f, ok := i.(Middleware)
		if ok {
			f(ctx)
		} else {
			ff, ok := i.(Handler)
			if ok {
				ff(ctx)
			}
		}
	}
}

func (app *App) combHandlers(h Handler, mw []Middleware) []interface{} {
	var rev []interface{}
	for _, k := range mw {
		rev = append(rev, k)
	}
	rev = append(rev, h)
	return rev
}

func GetGormFunc() func() *gorm.DB {
	f := func() *gorm.DB {
		if !gormC.EnableGorm {
			panic("you are trying to use gorm but it's not enabled, you can enable it in the file config/gorm.go")
		}
		return ResolveGorm()
	}
	return f
}

func NewGorm() *gorm.DB {
	var err error
	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		db, err = mysqlConnect()
	case "postgres":
		db, err = postgresConnect()
	case "sqlite":
		sqlitePath := os.Getenv("SQLITE_DB_PATH")
		fullSqlitePath := path.Join(basePath, sqlitePath)
		_, err := os.Stat(fullSqlitePath)
		if err != nil {
			panic(fmt.Sprintf("error locating sqlite file: %v", err.Error()))
		}
		db, err = gorm.Open(sqlite.Open(fullSqlitePath), &gorm.Config{})
	default:
		panic("database driver not selected")
	}
	if gormC.EnableGorm && err != nil {
		panic(fmt.Sprintf("gorm has problem connecting to %v, (if it's not needed you can disable it in config/gorm.go): %v", os.Getenv("DB_DRIVER"), err))
	}
	return db
}

func ResolveGorm() *gorm.DB {
	if db != nil {
		return db
	}
	db = NewGorm()
	return db
}

func resolveCache() func() *Cache {
	f := func() *Cache {
		if !cacheC.EnableCache {
			panic("you are trying to use cache but it's not enabled, you can enable it in the file config/cache.go")
		}
		return NewCache(cacheC)
	}
	return f
}

func postgresConnect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB_NAME"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_SSL_MODE"),
		os.Getenv("POSTGRES_TIMEZONE"),
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func mysqlConnect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USERNAME"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DB_NAME"),
		os.Getenv("MYSQL_CHARSET"),
	)
	return gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{})
}

func getJWT() func() *JWT {
	f := func() *JWT {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			panic("jwt secret key is not set")
		}
		lifetimeStr := os.Getenv("JWT_LIFESPAN_MINUTES")
		if lifetimeStr == "" {
			lifetimeStr = "10080" // 7 days
		}
		lifetime64, err := strconv.ParseInt(lifetimeStr, 10, 32)
		if err != nil {
			panic(err)
		}
		lifetime := int(lifetime64)
		return newJWT(JWTOptions{
			SigningKey:      secret,
			LifetimeMinutes: lifetime,
		})
	}
	return f
}

func getValidator() func() *Validator {
	f := func() *Validator {
		return &Validator{}
	}
	return f
}

func resloveHashing() func() *Hashing {
	f := func() *Hashing {
		return &Hashing{}
	}
	return f
}

func resolveMailer() func() *Mailer {
	f := func() *Mailer {
		if mailer != nil {
			return mailer
		}
		var m *Mailer
		var emailsDriver string
		if os.Getenv("EMAILS_DRIVER") == "" {
			emailsDriver = "SMTP"
		}
		switch emailsDriver {
		case "SMTP":
			m = initiateMailerWithSMTP()
		case "sparkpost":
			m = initiateMailerWithSparkPost()
		case "sendgrid":
			m = initiateMailerWithSendGrid()
		case "mailgun":
			return initiateMailerWithMailGun()
		default:
			m = initiateMailerWithSMTP()
		}
		mailer = m
		return mailer
	}
	return f
}

func resolveEventsManager() func() *EventsManager {
	f := func() *EventsManager {
		return ResolveEventsManager()
	}
	return f
}

func resolveLogger() func() *logger.Logger {
	f := func() *logger.Logger {
		return loggr
	}
	return f
}

func (app *App) MakeDirs(dirs ...string) {
	o := syscall.Umask(0)
	defer syscall.Umask(o)
	for _, dir := range dirs {
		os.MkdirAll(path.Join(basePath, dir), 0766)
	}
}

func (app *App) SetRequestConfig(r RequestConfig) {
	requestC = r
}

func (app *App) SetGormConfig(g GormConfig) {
	gormC = g
}

func (app *App) SetCacheConfig(c CacheConfig) {
	cacheC = c
}

func (app *App) SetBasePath(path string) {
	basePath = path
}

func DisableEvents() {
	disableEvents = true
}

func EnableEvents() {
	disableEvents = false
}
