// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"

	"github.com/gocondor/core/env"
	"github.com/gocondor/core/logger"
	"github.com/harranali/mailing"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/acme/autocert"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var logsDriver logger.LogsDriver
var loggr *logger.Logger
var requestC RequestConfig
var jwtC JWTConfig
var gormC GormConfig
var cacheC CacheConfig
var db *gorm.DB
var mailer *mailing.Mailer

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
	logsDriver = d
}

func (app *App) Bootstrap() {
	NewRouter()
	// database.New()
	// cache.New(app.Features.Cache)
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
			logger:       loggr,
			GetValidator: getValidator(),
			GetJWT:       getJWT(),
			GetGorm:      GetGormFunc(),
			GetCache:     resolveCache(),
			GetHashing:   resloveHashing(),
			GetMailer:    resolveMailer(),
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
			code := ctx.Response.getStatusCode()
			if code == 0 {
				code = http.StatusOK
			}
			w.WriteHeader(code)
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
		if sqlitePath == "" {
			panic("wrong path to sqlite file")
		}
		db, err = gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
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
			SigningKey: secret,
			Lifetime:   lifetime,
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

// TODO implement the mail wrapper
func resolveMailer() func() *mailing.Mailer {
	f := func() *mailing.Mailer {
		var mailer *mailing.Mailer
		var emailsDriver string
		if os.Getenv("EMAILS_DRIVER") == "" {
			emailsDriver = "SMTP"
		}
		switch emailsDriver {
		case "SMTP":
			mailer = initiateMailerWithSMTP()
		case "sparkpost":
			mailer = initiateMailerWithSparkPost()
		case "sendgrid":
			mailer = initiateMailerWithSendGrid()
		case "mailgun":
			return initiateMailerWithMailGun()
		default:
			mailer = initiateMailerWithSMTP()
		}
		return mailer
	}
	return f
}

func initiateMailerWithSMTP() *mailing.Mailer {
	portStr := os.Getenv("SMTP_PORT")
	if portStr == "" {
		panic("error reading smtp port env var")
	}
	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("error parsing smtp port env var: %v", err))
	}
	skipTlsVerifyStr := os.Getenv("SMTP_TLS_SKIP_VERIFY_HOST")
	if skipTlsVerifyStr == "" {
		panic("error reading smtp tls verify env var")
	}
	skipTlsVerify, err := strconv.ParseBool(skipTlsVerifyStr)
	if err != nil {
		panic(fmt.Sprintf("error parsing smtp tls verify env var: %v", err))
	}
	return mailing.NewMailerWithSMTP(&mailing.SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     int(port),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		TLSConfig: tls.Config{
			ServerName:         os.Getenv("SMTP_HOST"),
			InsecureSkipVerify: skipTlsVerify,
		},
	})
}

func initiateMailerWithSparkPost() *mailing.Mailer {
	apiVersionStr := os.Getenv("SPARKPOST_API_VERSION")
	if apiVersionStr == "" {
		panic("error reading sparkpost base url env var")
	}
	apiVersion, err := strconv.ParseInt(apiVersionStr, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("error parsing sparkpost base url env var: %v", apiVersion))
	}
	return mailing.NewMailerWithSparkPost(&mailing.SparkPostConfig{
		BaseUrl:    os.Getenv("SPARKPOST_BASE_URL"),
		ApiKey:     os.Getenv("SPARKPOST_API_KEY"),
		ApiVersion: int(apiVersion),
	})
}

func initiateMailerWithSendGrid() *mailing.Mailer {
	return mailing.NewMailerWithSendGrid(&mailing.SendGridConfig{
		Host:     os.Getenv("SENDGRID_HOST"),
		Endpoint: os.Getenv("SENDGRID_ENDPOINT"),
		ApiKey:   os.Getenv("SENDGRID_API_KEY"),
	})
}

func initiateMailerWithMailGun() *mailing.Mailer {
	skipTlsVerifyStr := os.Getenv("MAILGUN_TLS_SKIP_VERIFY_HOST")
	if skipTlsVerifyStr == "" {
		panic("error reading mailgun tls verify env var")
	}
	skipTlsVerify, err := strconv.ParseBool(skipTlsVerifyStr)
	if err != nil {
		panic(fmt.Sprintf("error parsing mailgun tls verify env var: %v", err))
	}
	return mailing.NewMailerWithMailGun(&mailing.MailGunConfig{
		Domain:              os.Getenv("MAILGUN_DOMAIN"),
		APIKey:              os.Getenv("MAILGUN_API_KEY"),
		SkipTLSVerification: skipTlsVerify,
	})
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
