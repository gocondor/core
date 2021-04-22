// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package pkgintegrator

import "github.com/gin-gonic/gin"

// PKGIntegrator integrates packages to gin engine
type PKGIntegrator struct {
	integrations []gin.HandlerFunc
}

var integrator *PKGIntegrator

// New initiates the package integrator
func New() *PKGIntegrator {
	integrator = &PKGIntegrator{}

	return integrator
}

// Resolve returns intitiated PKGIntegrator
func Resolve() *PKGIntegrator {
	return integrator
}

// Integrate adds package to be integrated to gin context
func (i *PKGIntegrator) Integrate(pkgIntegration gin.HandlerFunc) {
	i.integrations = append(i.integrations, pkgIntegration)
}

// GetIntegrations returns a list of package to be integrated to gin context
func (i *PKGIntegrator) GetIntegrations() []gin.HandlerFunc {
	return i.integrations
}
