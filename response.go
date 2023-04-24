package core

import "net/http"

type Response struct {
	headers        []header
	textBody       string
	jsonBody       string
	responseWriter http.ResponseWriter
}

type header struct {
	key string
	val string
}

func (rs *Response) setTextBody(body string) {
	rs.textBody = body
}

func (rs *Response) setJsonBody(body string) {
	rs.jsonBody = body
}

func (rs *Response) getTextBody() string {
	return rs.textBody
}

func (rs *Response) getJsonBody() string {
	return rs.jsonBody
}

func (rs *Response) setHeader(key string, val string) {
	h := header{
		key: key,
		val: val,
	}
	rs.headers = append(rs.headers, h)
}

func (rs *Response) getHeader(key string) string {

	return ""
}
