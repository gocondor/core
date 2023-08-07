// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

type Route struct {
	Method      string
	Path        string
	Handler     Handler
	Middlewares []Middleware
}

type Router struct {
	Routes []Route
}

var router *Router

func NewRouter() *Router {
	router = &Router{
		[]Route{},
	}
	return router
}

func ResolveRouter() *Router {
	return router
}

func (r *Router) Get(path string, handler Handler, middlewares ...Middleware) *Router {
	r.Routes = append(r.Routes, Route{
		Method:      GET,
		Path:        path,
		Handler:     handler,
		Middlewares: middlewares,
	})
	return r
}

func (r *Router) Post(path string, handler Handler, middlewares ...Middleware) *Router {
	r.Routes = append(r.Routes, Route{
		Method:      POST,
		Path:        path,
		Handler:     handler,
		Middlewares: middlewares,
	})
	return r
}

func (r *Router) Delete(path string, handler Handler, middlewares ...Middleware) *Router {
	r.Routes = append(r.Routes, Route{
		Method:      DELETE,
		Path:        path,
		Handler:     handler,
		Middlewares: middlewares,
	})
	return r
}

func (r *Router) Patch(path string, handler Handler, middlewares ...Middleware) *Router {
	r.Routes = append(r.Routes, Route{
		Method:      PATCH,
		Path:        path,
		Handler:     handler,
		Middlewares: middlewares,
	})
	return r
}

func (r *Router) Put(path string, handler Handler, middlewares ...Middleware) *Router {
	r.Routes = append(r.Routes, Route{
		Method:      PUT,
		Path:        path,
		Handler:     handler,
		Middlewares: middlewares,
	})
	return r
}

func (r *Router) Options(path string, handler Handler, middlewares ...Middleware) *Router {
	r.Routes = append(r.Routes, Route{
		Method:      OPTIONS,
		Path:        path,
		Handler:     handler,
		Middlewares: middlewares,
	})
	return r
}

func (r *Router) Head(path string, handler Handler, middlewares ...Middleware) *Router {
	r.Routes = append(r.Routes, Route{
		Method:      HEAD,
		Path:        path,
		Handler:     handler,
		Middlewares: middlewares,
	})
	return r
}

func (r *Router) GetRoutes() []Route {
	return r.Routes
}
