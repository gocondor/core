// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"github.com/gin-gonic/gin"
	"github.com/gincoat/core/jwtloader"
)

// RegisterJwt returns a gin handler func with jwt variable set in gin context
func RegisterJwt(jwt *jwtloader.JwtLoader) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("jwt", jwt)
		c.Next()
	}
}
