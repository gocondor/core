// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

// Route ais struct describes a specific route
type Route struct {
	Method  string
	Path    string
	Handler Handler
}

// Router handles routing
type Router struct {
	Routes []Route
}

var router *Router

// New initiates new router
func NewRouter() *Router {
	router = &Router{
		[]Route{},
	}
	return router
}

// Resolve resolves an already initiated router
func ResolveRouter() *Router {
	return router
}

// Get is a definition for get request
func (r *Router) Get(path string, handler Handler) *Router {
	r.Routes = append(r.Routes, Route{
		Method:  "get",
		Path:    path,
		Handler: handler,
	})

	return r
}

// Post is a definition for post request
func (r *Router) Post(path string, handler Handler) *Router {
	r.Routes = append(r.Routes, Route{
		Method:  "post",
		Path:    path,
		Handler: handler,
	})

	return r
}

// Delete is a definition for delete request
func (r *Router) Delete(path string, handler Handler) *Router {
	r.Routes = append(r.Routes, Route{
		Method:  "delete",
		Path:    path,
		Handler: handler,
	})

	return r
}

// Put is a definition for put request
func (r *Router) Put(path string, handler Handler) *Router {
	r.Routes = append(r.Routes, Route{
		Method:  "put",
		Path:    path,
		Handler: handler,
	})

	return r
}

// Options is a definition for options request
func (r *Router) Options(path string, handler Handler) *Router {
	r.Routes = append(r.Routes, Route{
		Method:  "options",
		Path:    path,
		Handler: handler,
	})

	return r
}

// Head is a definition for head request
func (r *Router) Head(path string, handler Handler) *Router {
	r.Routes = append(r.Routes, Route{
		Method:  "head",
		Path:    path,
		Handler: handler,
	})

	return r
}

// GetRoutes returns all Defined routes
func (r *Router) GetRoutes() []Route {
	// get the routing groups
	// combine the routing groups routes
	return r.Routes
}
