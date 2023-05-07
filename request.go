// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Request struct {
	HttpRequest    *http.Request
	httpPathParams httprouter.Params
}
