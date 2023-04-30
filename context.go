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
	Logger   *Logger
}

func (c *Context) Debug(d interface{}) {
	m := reflect.ValueOf(d)
	if m.Kind() == reflect.Pointer {
		m = m.Elem()
	}
	formatted := fmt.Sprintf("Type: %T | underlaying type: %v | value: %v", d, m.Kind(), d)
	fmt.Println(formatted)
	c.Response.responseWriter.Write([]byte(formatted))
}

func (c *Context) Next() {
	ResolveApp().Next(c)
}
