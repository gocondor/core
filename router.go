// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

type Route struct {
	Method   string
	Path     string
	Handlers []Handler
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

func (r *Router) Get(path string, handlers ...Handler) *Router {
	r.Routes = append(r.Routes, Route{
		Method:   "get",
		Path:     path,
		Handlers: handlers,
	})
	return r
}

func (r *Router) Post(path string, handlers ...Handler) *Router {
	r.Routes = append(r.Routes, Route{
		Method:   "post",
		Path:     path,
		Handlers: handlers,
	})
	return r
}

func (r *Router) Delete(path string, handlers ...Handler) *Router {
	r.Routes = append(r.Routes, Route{
		Method:   "delete",
		Path:     path,
		Handlers: handlers,
	})
	return r
}

func (r *Router) Put(path string, handlers ...Handler) *Router {
	r.Routes = append(r.Routes, Route{
		Method:   "put",
		Path:     path,
		Handlers: handlers,
	})
	return r
}

func (r *Router) Options(path string, handlers ...Handler) *Router {
	r.Routes = append(r.Routes, Route{
		Method:   "options",
		Path:     path,
		Handlers: handlers,
	})
	return r
}

func (r *Router) Head(path string, handlers ...Handler) *Router {
	r.Routes = append(r.Routes, Route{
		Method:   "head",
		Path:     path,
		Handlers: handlers,
	})
	return r
}

func (r *Router) GetRoutes() []Route {
	return r.Routes
}
