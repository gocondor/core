package core

import (
	"fmt"
	"net/http"
)

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

func (rs *Response) WriteText(body interface{}) {
	if rs.textBody == "" {
		rs.textBody = rs.castBasicVarsToString(body)
	}
}

func (rs *Response) WriteJson(body string) {
	if rs.textBody == "" {
		rs.jsonBody = body
	}
}

func (rs *Response) getTextBody() string {
	return rs.textBody
}

func (rs *Response) getJsonBody() string {
	return rs.jsonBody
}

func (rs *Response) SetHeader(key string, val string) {
	h := header{
		key: key,
		val: val,
	}
	rs.headers = append(rs.headers, h)
}

func (rs *Response) getHeaders() []header {
	return rs.headers
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
		panic(fmt.Sprintf("unsupported data type %v!", dataType))
	}
}

func (rs *Response) reset() {
	rs.textBody = ""
	rs.jsonBody = ""
}
