package core

type Context struct {
	Request  *Request
	Response *Response
	Logger   *Logger
}
