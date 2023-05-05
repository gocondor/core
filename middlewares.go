// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

type Middlewares struct {
	middlewares []Handler
}

var m *Middlewares

func NewMiddlewares() *Middlewares {
	m = &Middlewares{}
	return m
}

func ResolveMiddlewares() *Middlewares {
	return m
}

func (m *Middlewares) Attach(mw Handler) *Middlewares {
	m.middlewares = append(m.middlewares, mw)

	return m
}

func (m *Middlewares) GetMiddlewares() []Handler {
	return m.middlewares
}

func (m *Middlewares) getByIndex(i int) Handler {
	for key, _ := range m.middlewares {
		if key == i {
			return m.middlewares[i]
		}
	}
	return nil
}