// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"reflect"
)

type Context struct {
	Request  *Request
	Response *Response
	logger   *Logger
}

func (c *Context) DebugAny(variable interface{}) {
	m := reflect.ValueOf(variable)
	if m.Kind() == reflect.Pointer {
		m = m.Elem()
	}
	formatted := fmt.Sprintf("Type: %T (%v) | value: %v", variable, m.Kind(), variable)
	fmt.Println(formatted)
	c.Response.HttpResponseWriter.Write([]byte(formatted))
}

func (c *Context) Next() {
	ResolveApp().Next(c)
}

func (c *Context) prepare(ctx *Context) {
	ctx.Request.HttpRequest.ParseMultipartForm(20000000)
}

func (c *Context) LogInfo(msg interface{}) {
	ResolveLogger().Info(msg)
}

func (c *Context) LogError(msg interface{}) {
	ResolveLogger().Error(msg)
}

func (c *Context) LogWarning(msg interface{}) {
	ResolveLogger().Warning(msg)
}

func (c *Context) LogDebug(msg interface{}) {
	ResolveLogger().Debug(msg)
}

func (c *Context) GetPathParam(key string) string {
	return c.Request.httpPathParams.ByName(key)
}

func (c *Context) GetRequestParam(key string) string {
	return c.Request.HttpRequest.FormValue(key)
}

func (c *Context) RequestParamExists(key string) bool {
	return c.Request.HttpRequest.Form.Has(key)
}

// TODO implement
func (c *Context) GetRequestFile(name string) {
	file, fileHeader, _ := c.Request.HttpRequest.FormFile(name)
	c.LogInfo(fileHeader.Filename)
	fmt.Println(file)
}
