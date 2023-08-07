// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"net/http"
)

type Response struct {
	headers             []header
	body                []byte
	statusCode          int
	contentType         string
	overrideContentType string
	isTerminated        bool
	HttpResponseWriter  http.ResponseWriter
}

type header struct {
	key string
	val string
}

// TODO add doc
func (rs *Response) Any(body any) *Response {
	if rs.isTerminated == false {
		rs.contentType = CONTENT_TYPE_HTML
		rs.body = []byte(rs.castBasicVarsToString(body))
	}

	return rs
}

// TODO add doc
func (rs *Response) Json(body string) *Response {
	if rs.isTerminated == false {
		rs.contentType = CONTENT_TYPE_JSON
		rs.body = []byte(body)
	}
	return rs
}

// TODO add doc
func (rs *Response) Text(body string) *Response {
	if rs.isTerminated == false {
		rs.contentType = CONTENT_TYPE_TEXT
		rs.body = []byte(body)
	}
	return rs
}

// TODO add doc
func (rs *Response) HTML(body string) *Response {
	if rs.isTerminated == false {
		rs.contentType = CONTENT_TYPE_HTML
		rs.body = []byte(body)
	}
	return rs
}

// TODO add doc
func (rs *Response) SetStatusCode(code int) *Response {
	if rs.isTerminated == false {
		rs.statusCode = code
	}

	return rs
}

// TODO add doc
func (rs *Response) SetContentType(c string) *Response {
	if rs.isTerminated == false {
		rs.overrideContentType = c
	}

	return rs
}

// TODO add doc
func (rs *Response) SetHeader(key string, val string) *Response {
	if rs.isTerminated == false {
		h := header{
			key: key,
			val: val,
		}
		rs.headers = append(rs.headers, h)
	}
	return rs
}

func (rs *Response) ForceSendResponse() {
	rs.isTerminated = true
}

func (rs *Response) castBasicVarsToString(data interface{}) string {
	switch dataType := data.(type) {
	case string:
		return fmt.Sprintf("%v", data)
	case []byte:
		d := data.(string)
		return fmt.Sprintf("%v", d)
	case int:
		intVar, _ := data.(int)
		return fmt.Sprintf("%v", intVar)
	case int8:
		int8Var := data.(int8)
		return fmt.Sprintf("%v", int8Var)
	case int16:
		int16Var := data.(int16)
		return fmt.Sprintf("%v", int16Var)
	case int32:
		int32Var := data.(int32)
		return fmt.Sprintf("%v", int32Var)
	case int64:
		int64Var := data.(int64)
		return fmt.Sprintf("%v", int64Var)
	case uint:
		uintVar, _ := data.(uint)
		return fmt.Sprintf("%v", uintVar)
	case uint8:
		uint8Var := data.(uint8)
		return fmt.Sprintf("%v", uint8Var)
	case uint16:
		uint16Var := data.(uint16)
		return fmt.Sprintf("%v", uint16Var)
	case uint32:
		uint32Var := data.(uint32)
		return fmt.Sprintf("%v", uint32Var)
	case uint64:
		uint64Var := data.(uint64)
		return fmt.Sprintf("%v", uint64Var)
	case float32:
		float32Var := data.(float32)
		return fmt.Sprintf("%v", float32Var)
	case float64:
		float64Var := data.(float64)
		return fmt.Sprintf("%v", float64Var)
	case complex64:
		complex64Var := data.(complex64)
		return fmt.Sprintf("%v", complex64Var)
	case complex128:
		complex128Var := data.(complex128)
		return fmt.Sprintf("%v", complex128Var)
	case bool:
		boolVar := data.(bool)
		return fmt.Sprintf("%v", boolVar)
	default:
		panic(fmt.Sprintf("unsupported response data type %v!", dataType))
	}
}

func (rs *Response) reset() {
	rs.body = nil
	rs.statusCode = http.StatusOK
	rs.contentType = CONTENT_TYPE_HTML
	rs.overrideContentType = ""
	rs.isTerminated = false
}
