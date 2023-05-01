package core

import (
	"testing"
)

func TestWriteText(t *testing.T) {
	res := Response{}
	v := "test-text"
	res.WriteText(v)

	if res.getTextBody() != v {
		t.Errorf("failed writing text")
	}
}

func TestWriteJson(t *testing.T) {
	res := Response{}
	j := "{\"name\": \"test\"}"
	res.WriteJson(j)

	if res.getJsonBody() != j {
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
	res.WriteText("test text")
	if res.getTextBody() == "" {
		t.Errorf("expecting textBody to not be empty, found empty")
	}
	j := "{\"name\": \"test\"}"
	res.WriteJson(j)
	if res.getJsonBody() == "" {
		t.Errorf("expecting JsonBody to not be empty, found empty")
	}

	res.reset()

	if !(res.getTextBody() == "" && res.getJsonBody() == "") {
		t.Errorf("failed testing response reset()")
	}
}
