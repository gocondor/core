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
