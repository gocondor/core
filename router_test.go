// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core_test

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/gocondor/core/routing"
)

func TestNewRouter(t *testing.T) {
	r := New()

	if fmt.Sprintf("%T", r) != "*routing.Router" {
		t.Error("failed asserting initiation of new router")
	}
}

func TestResolveRouter(t *testing.T) {
	r := Resolve()
	if fmt.Sprintf("%T", r) != "*routing.Router" {
		t.Error("failed resolve router variable")
	}
}

func TestGetRequest(t *testing.T) {
	r := New()
	handler := func(c *gin.Context) {}
	r.Get("/", handler)

	route := r.GetRoutes()[0]
	if route.Method != "get" || route.Path != "/" {
		t.Errorf("failed adding route with get http method")
	}
}

func TestPostRequest(t *testing.T) {
	r := New()
	handler := func(c *gin.Context) {}
	r.Post("/", handler)

	route := r.GetRoutes()[0]
	if route.Method != "post" || route.Path != "/" {
		t.Errorf("failed adding route with post http method")
	}
}

func TestDeleteRequest(t *testing.T) {
	r := New()
	handler := func(c *gin.Context) {}
	r.Delete("/", handler)

	route := r.GetRoutes()[0]
	if route.Method != "delete" || route.Path != "/" {
		t.Errorf("failed adding route with delete http method")
	}
}

func TestPutRequest(t *testing.T) {
	r := New()
	handler := func(c *gin.Context) {}
	r.Put("/", handler)

	route := r.GetRoutes()[0]
	if route.Method != "put" || route.Path != "/" {
		t.Errorf("failed adding route with put http method")
	}
}

func TestOptionsRequest(t *testing.T) {
	r := New()
	handler := func(c *gin.Context) {}
	r.Options("/", handler)

	route := r.GetRoutes()[0]
	if route.Method != "options" || route.Path != "/" {
		t.Errorf("failed adding route with options http method")
	}
}

func TestHeadRequest(t *testing.T) {
	r := New()
	handler := func(c *gin.Context) {}
	r.Head("/", handler)

	route := r.GetRoutes()[0]
	if route.Method != "head" || route.Path != "/" {
		t.Errorf("failed adding route with head http method")
	}
}

func TestGetRoutesRequest(t *testing.T) {
	r := New()
	r.Get("/", func(c *gin.Context) {})
	r.Post("/", func(c *gin.Context) {})

	if len(r.GetRoutes()) != 2 {
		t.Errorf("failed getting added routes")
	}
}
