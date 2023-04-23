package core

import "fmt"

type Context struct {
	request  *Request
	response *Response
}

func (c *Context) Check() {
	fmt.Println("context check!!!")
}
