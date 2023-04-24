package core

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/fatih/structs"
)

type Context struct {
	Request     *Request
	ResponseBag *Response
	Logger      *Logger
}

func (c *Context) Response(data interface{}) *Response {
	dataMeta := reflect.ValueOf(data)
	if dataMeta.Kind() == reflect.Pointer {
		dataMeta = dataMeta.Elem()
	}

	if dataMeta.Kind() == reflect.Struct {
		str, _ := json.Marshal(structs.Map(data))
		c.ResponseBag.setResponseBody(string(str))

		return c.ResponseBag
	}

	str := c.CastBasicVarsToString(data)
	c.ResponseBag.setResponseBody(str)
	return c.ResponseBag
}

func (c *Context) CastBasicVarsToString(data interface{}) string {
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
