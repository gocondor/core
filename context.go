package core

import "fmt"

type Context struct {
	request  *Request
	response *Response
}

// TODO remote when things are ready
func (c *Context) Check() {
	fmt.Println("context check!!!")
}
