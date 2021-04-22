// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package middlewares

import (
	"github.com/gin-gonic/gin"
)

// MiddlewaresUtil handles middlewares registration
type MiddlewaresUtil struct {
	middlewares []gin.HandlerFunc
}

//Middleware a function defines a middleware
type Middleware func(c *gin.Context)

var middlewaresUtil *MiddlewaresUtil

//New initiates a new middlware util
func New() *MiddlewaresUtil {
	middlewaresUtil = &MiddlewaresUtil{}
	return middlewaresUtil
}

// Resolve returns an already initated middleware util
func Resolve() *MiddlewaresUtil {
	return middlewaresUtil
}

// Attach attach a middleware globally to the app
func (m *MiddlewaresUtil) Attach(mw gin.HandlerFunc) *MiddlewaresUtil {
	m.middlewares = append(m.middlewares, mw)

	return middlewaresUtil
}

// GetMiddlewares get all attached middlewares
func (m *MiddlewaresUtil) GetMiddlewares() []gin.HandlerFunc {
	return m.middlewares
}
