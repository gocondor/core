package core

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct{}

var j *JWT

func newJWT() *JWT {
	j = &JWT{}
	return j
}
func resolveJWT() *JWT {
	return j
}

type claims struct {
	J []byte
	jwt.RegisteredClaims
}

var signingKey = []byte("thesigningstring") // TODO extract to config
func (j *JWT) GenerateToken(payload map[string]interface{}, expiresAt time.Time) (string, error) {
	claims, err := mapClaims(payload, expiresAt)
	if err != nil {
		return "", err
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ := t.SignedString(signingKey)
	return token, nil
}

func (j *JWT) DecodeToken(token string) (payload map[string]interface{}, err error) {
	t, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	c, ok := t.Claims.(*claims)
	if !ok {
		return nil, errors.New("error decoding token")
	}
	expiresAt := time.Unix(c.ExpiresAt.Unix(), 0)
	et := time.Now().Compare(expiresAt)
	if et != -1 {
		return nil, errors.New("token has expired")
	}
	err = json.Unmarshal(c.J, &payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (j *JWT) isValid(token string) bool {

	return false
}

func (j *JWT) hasExpired(token string) bool {
	return false
}

func mapClaims(data map[string]interface{}, expiresAt time.Time) (jwt.Claims, error) {
	j, err := json.Marshal(data)
	if err != nil {
		return claims{}, err
	}
	r := claims{
		j,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	return r, nil
}
