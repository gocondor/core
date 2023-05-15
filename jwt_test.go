package core

import (
	"fmt"
	"testing"
	"time"
)

func TestNewJWT(t *testing.T) {
	j := newJWT(JWTOptions{
		SigningKey: "testsigning",
		Lifetime:   2,
	})

	if fmt.Sprintf("%T", j) != "*core.JWT" {
		t.Errorf("failed testing new jwt")
	}
}

func TestResolveJWT(t *testing.T) {
	initiateJWTHelper(t)
	j := resolveJWT()
	if fmt.Sprintf("%T", j) != "*core.JWT" {
		t.Errorf("failed testing resolve jwt")
	}
}

func TestGenerateToken(t *testing.T) {
	j := initiateJWTHelper(t)
	token, err := j.GenerateToken(map[string]interface{}{
		"testKey": "testVal",
	})
	if err != nil || token == "" {
		t.Errorf("error testing generate jwt token")
	}
	d, err := j.DecodeToken(token)
	if err != nil {
		t.Errorf("error testing generate jwt token")
	}
	if d["testKey"] != "testVal" {
		t.Errorf("error testing generate jwt token")
	}
}

func TestDecodeToken(t *testing.T) {
	j := initiateJWTHelper(t)
	token, err := j.GenerateToken(map[string]interface{}{
		"testKey": "testVal",
	})
	if err != nil || token == "" {
		t.Errorf("failed testing decode token")
	}
	d, err := j.DecodeToken(token)
	if err != nil {
		t.Errorf("error testing decode jwt token")
	}
	if d["testKey"] != "testVal" {
		t.Errorf("error testing decode jwt token")
	}

	d, err = j.DecodeToken("test-token")
	if err == nil {
		t.Errorf("error testing decode jwt token")
	}
}

func TestHasExpired(t *testing.T) {
	j := initiateJWTHelper(t)
	token, err := j.GenerateToken(map[string]interface{}{
		"testKey": "testVal",
	})
	if err != nil || token == "" {
		t.Errorf("failed testing decode token")
	}
	err = j.HasExpired(token)
	if err != nil {
		t.Errorf("error testing decode jwt token")
	}
}

func TestMapClaims(t *testing.T) {
	c, err := mapClaims(map[string]interface{}{
		"testKey": "testVal",
	}, time.Now())

	if err != nil {
		t.Errorf("failed testing map claims")
	}
	if fmt.Sprintf("%T", c) != "core.claims" {
		t.Errorf("failed testing map claims")
	}
}

func initiateJWTHelper(t *testing.T) *JWT {
	t.Helper()
	j := newJWT(JWTOptions{
		SigningKey: "testsigning",
		Lifetime:   2,
	})
	return j
}
