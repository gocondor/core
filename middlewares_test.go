// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewMiddlewares(t *testing.T) {
	mw := NewMiddlewares()
	if fmt.Sprintf("%T", mw) != "*core.Middlewares" {
		t.Errorf("failed testing new middleware")
	}
}

func TestResloveMiddleWares(t *testing.T) {
	NewMiddlewares()
	mw := ResolveMiddlewares()
	if fmt.Sprintf("%T", mw) != "*core.Middlewares" {
		t.Errorf("failed resolve middlewares")
	}
}

func TestAttach(t *testing.T) {
	mw := NewMiddlewares()
	tmw := func(c *Context) {
		c.LogInfo("Testing!")
	}
	mw.Attach(tmw)
	mws := mw.getByIndex(0)
	if reflect.ValueOf(tmw).Pointer() != reflect.ValueOf(mws).Pointer() {
		t.Errorf("Failed testing attach middleware")
	}
}

func TestGetMiddleWares(t *testing.T) {
	mw := NewMiddlewares()
	t1 := func(c *Context) {
		c.LogInfo("testing1!")
	}
	t2 := func(c *Context) {
		c.LogInfo("testing2!")
	}
	mw.Attach(t1)
	mw.Attach(t2)
	if len(mw.GetMiddlewares()) != 2 {
		t.Errorf("failed testing get middlewares")
	}
}

func TestMiddlewareGetByIndex(t *testing.T) {
	mw := NewMiddlewares()
	t1 := func(c *Context) {
		c.LogInfo("testing!")
	}
	mw.Attach(t1)
	if reflect.ValueOf(mw.getByIndex(0)).Pointer() != reflect.ValueOf(t1).Pointer() {
		t.Errorf("failed testing get by index")
	}
}
