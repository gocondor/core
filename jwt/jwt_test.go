package jwt_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/gocondor/core/jwt"
)

func TestCreateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "qwertyuio")
	os.Setenv("JWT_LIFESPAN_MINUTES", "15")

	var userID uint = 333
	j := New()

	token, err := j.CreateToken(userID)
	if err != nil {
		t.Error("failed create new jwt token.", err)
	}

	ok, err := j.ValidateToken(token)
	if !ok {
		t.Error("failed create new jwt token.", err)
	}
}

func TestCreateRefreshToken(t *testing.T) {
	os.Setenv("JWT_REFRESH_TOKEN_SECRET", "qwertyuio")
	os.Setenv("JWT_REFRESH_TOKEN_LIFESPAN_HOURS", "15")

	var userID uint = 333
	j := New()

	token, err := j.CreateRefreshToken(userID)
	if err != nil {
		t.Error("failed create new jwt refresh token.", err)
	}

	ok, err := j.ValidateToken(token)
	if !ok {
		t.Error("failed create new jwt refresh token.", err)
	}
}

func TestExtractToken(t *testing.T) {
	j := New()
	g := gin.New()

	g.GET("/", func(c *gin.Context) {
		token, err := j.ExtractToken(c)
		if err != nil {
			t.Error("failed to extract jwt token, expect the token to be in format 'bearer: [token]'. ", err)
		}
		_, err = j.ValidateToken(token)
		if err != nil {
			t.Error("failed to extract jwt token, expect the token to be in format 'bearer: [token]'. ", err)
		}
	})
	s := httptest.NewServer(g)
	var userID uint = 333
	token, _ := j.CreateToken(userID)
	rq, _ := http.NewRequest("GET", s.URL, nil)
	rq.Header.Set("Authorization", "bear: "+token)
	_, err := s.Client().Do(rq)
	if err != nil {
		t.Error("failed to extract jwt token.", err)
	}
}

func TestValidateToken(t *testing.T) {
	j := New()
	var userID uint = 333
	token, _ := j.CreateToken(userID)
	_, err := j.ValidateToken(token)
	if err != nil {
		t.Error("failed assert validate jwt token")
	}

}

func TestDecodeToken(t *testing.T) {
	j := New()
	var userID uint = 333
	token, _ := j.CreateToken(userID)

	id, err := j.DecodeToken(token)
	if err != nil {
		t.Error("failed decoding the token.", err)
	}

	if id != userID {
		t.Error("failed decoding the token")
	}
}
