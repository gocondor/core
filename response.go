package core

import "net/http"

type Response struct {
	header         map[string]string
	body           string
	responseWriter http.ResponseWriter
}

func (rs *Response) setResponseBody(body string) {
	rs.body = body
}

func (rs *Response) GetResponseBody() string {
	return rs.body
}
