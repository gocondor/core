package core

import (
	"fmt"
	"testing"
)

func TestWrite(t *testing.T) {
	res := Response{}
	v := "test-text"
	res.Write(v)

	if res.getBody() != v {
		t.Errorf("failed writing text")
	}
}

func TestWriteJson(t *testing.T) {
	res := Response{}
	j := "{\"name\": \"test\"}"
	res.WriteJson([]byte(j))

	if string(res.getJsonBody()) != j {
		t.Errorf("failed wrting jsom")
	}
}

func TestSetHeaders(t *testing.T) {
	res := Response{}
	res.SetHeader("testkey", "testval")

	headers := res.getHeaders()
	if len(headers) < 1 {
		t.Errorf("testing set header failed")
	}
}

func TestReset(t *testing.T) {
	res := Response{}
	res.Write("test text")
	if res.getBody() == "" {
		t.Errorf("expecting textBody to not be empty, found empty")
	}
	j := "{\"name\": \"test\"}"
	res.WriteJson([]byte(j))
	if string(res.getJsonBody()) == "" {
		t.Errorf("expecting JsonBody to not be empty, found empty")
	}

	res.reset()

	if !(res.getBody() == "" && string(res.getJsonBody()) == "") {
		t.Errorf("failed testing response reset()")
	}
}

func TestCastBasicVarToString(t *testing.T) {
	s := "test str"
	r := Response{}
	c := r.castBasicVarsToString(s)
	if fmt.Sprintf("%T", c) != "string" {
		t.Errorf("failed test cast basic var to string")
	}
	var i int = 3
	r = Response{}
	c = r.castBasicVarsToString(i)
	if fmt.Sprintf("%T", c) != "string" {
		t.Errorf("failed test cast basic var to string")
	}
	var i8 int8 = 3
	r = Response{}
	c = r.castBasicVarsToString(i8)
	if fmt.Sprintf("%T", c) != "string" {
		t.Errorf("failed test cast basic var to string")
	}
	var i16 int16 = 3
	r = Response{}
	c = r.castBasicVarsToString(i16)
	if fmt.Sprintf("%T", c) != "string" {
		t.Errorf("failed test cast basic var to string")
	}
	var i32 int32 = 3
	r = Response{}
	c = r.castBasicVarsToString(i32)
	if fmt.Sprintf("%T", c) != "string" {
		t.Errorf("failed test cast basic var to string")
	}
	var i64 int64 = 3
	r = Response{}
	c = r.castBasicVarsToString(i64)
	if fmt.Sprintf("%T", c) != "string" {
		t.Errorf("failed test cast basic var to string")
	}
	var ui uint = 3
	r = Response{}
	c = r.castBasicVarsToString(ui)
	if fmt.Sprintf("%T", c) != "string" {
		t.Errorf("failed test cast basic var to string")
	}
	var ui8 uint8 = 3
	r = Response{}
	c = r.castBasicVarsToString(ui8)
	if fmt.Sprintf("%T", c) != "string" {
		t.Errorf("failed test cast basic var to string")
	}
	var ui16 uint16 = 3
	r = Response{}
	c = r.castBasicVarsToString(ui16)
	if fmt.Sprintf("%T", c) != "string" {
		t.Errorf("failed test cast basic var to string")
	}
	var ui32 uint32 = 3
	r = Response{}
	c = r.castBasicVarsToString(ui32)
	if fmt.Sprintf("%T", c) != "string" {
		t.Errorf("failed test cast basic var to string")
	}
	var ui64 uint64 = 3
	r = Response{}
	c = r.castBasicVarsToString(ui64)
	if fmt.Sprintf("%T", c) != "string" {
		t.Errorf("failed test cast basic var to string")
	}
	var f32 float32 = 3
	r = Response{}
	c = r.castBasicVarsToString(f32)
	if fmt.Sprintf("%T", c) != "string" {
		t.Errorf("failed test cast basic var to string")
	}
	var f64 float64 = 3
	r = Response{}
	c = r.castBasicVarsToString(f64)
	if fmt.Sprintf("%T", c) != "string" {
		t.Errorf("failed test cast basic var to string")
	}
	var c64 complex64 = 3
	r = Response{}
	c = r.castBasicVarsToString(c64)
	if fmt.Sprintf("%T", c) != "string" {
		t.Errorf("failed test cast basic var to string")
	}
	var c128 complex128 = 3
	r = Response{}
	c = r.castBasicVarsToString(c128)
	if fmt.Sprintf("%T", c) != "string" {
		t.Errorf("failed test cast basic var to string")
	}
	var b bool = true
	r = Response{}
	c = r.castBasicVarsToString(b)
	if fmt.Sprintf("%T", c) != "string" {
		t.Errorf("failed test cast basic var to string")
	}
}
