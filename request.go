package core

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Request struct {
	httpRequest    *http.Request
	httpPathParams httprouter.Params
}
