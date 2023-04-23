package core

import "net/http"

type Response struct {
	httpResponseWriter http.ResponseWriter
}
