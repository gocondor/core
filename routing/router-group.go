package routing

import (
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// GroupRouter represents a router group definition
type GroupRouter struct {
	base   string
	Routes []Route
}

// Groups represents a holder for routing groups
type Groups struct {
	GroupsRouters map[string]*GroupRouter
}

// GroupsHolder is the var for routing groups
var GroupsHolder *Groups

// NewGroupsHolder initiates new routing groups holder
func NewGroupsHolder() {
	GroupsHolder = &Groups{
		GroupsRouters: map[string]*GroupRouter{},
	}
}

// ResolveGroupsHolder resolves new routing groups holder
func ResolveGroupsHolder() *Groups {
	return GroupsHolder
}

// Group initiates a new routing group
func (r *Router) Group(name string) *GroupRouter {
	rg := &GroupRouter{
		base:   name,
		Routes: []Route{},
	}

	GroupsHolder.GroupsRouters[strings.ReplaceAll(name, "/", "")] = rg
	return rg
}

// Get is a definition for get request
func (r *GroupRouter) Get(path string, handlers ...gin.HandlerFunc) *GroupRouter {
	r.Routes = append(r.Routes, Route{
		Method:   "get",
		Path:     path,
		Handlers: handlers,
	})

	return r
}

// Post is a definition for post request
func (r *GroupRouter) Post(path string, handlers ...gin.HandlerFunc) *GroupRouter {
	r.Routes = append(r.Routes, Route{
		Method:   "post",
		Path:     path,
		Handlers: handlers,
	})

	return r
}

// Delete is a definition for delete request
func (r *GroupRouter) Delete(path string, handlers ...gin.HandlerFunc) *GroupRouter {
	r.Routes = append(r.Routes, Route{
		Method:   "delete",
		Path:     path,
		Handlers: handlers,
	})

	return r
}

// Put is a definition for put request
func (r *GroupRouter) Put(path string, handlers ...gin.HandlerFunc) *GroupRouter {
	r.Routes = append(r.Routes, Route{
		Method:   "put",
		Path:     path,
		Handlers: handlers,
	})

	return r
}

// Options is a definition for options request
func (r *GroupRouter) Options(path string, handlers ...gin.HandlerFunc) *GroupRouter {
	r.Routes = append(r.Routes, Route{
		Method:   "options",
		Path:     path,
		Handlers: handlers,
	})

	return r
}

// Head is a definition for head request
func (r *GroupRouter) Head(path string, handlers ...gin.HandlerFunc) *GroupRouter {
	r.Routes = append(r.Routes, Route{
		Method:   "head",
		Path:     path,
		Handlers: handlers,
	})

	return r
}

//GetRoutes returns all Defined routes of a routing group
func (r *GroupRouter) GetRoutes() []Route {
	// join the routes with the group base
	for key, route := range r.Routes {
		route.Path = path.Join(r.base, route.Path)
		r.Routes[key] = route
	}

	return r.Routes
}

//GetGroupsRoutes returns all Defined routes of all routing groups
func (r *Groups) GetGroupsRoutes() (routes []Route) {
	for _, group := range r.GroupsRouters {

		for _, val := range group.GetRoutes() {
			routes = append(routes, val)
		}

	}

	return routes
}
