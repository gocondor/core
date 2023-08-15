package core

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	signingKey []byte
	expiresAt  time.Time
}
type JWTOptions struct {
	SigningKey      string
	LifetimeMinutes int
}

var j *JWT

func newJWT(opts JWTOptions) *JWT {
	d := time.Duration(opts.LifetimeMinutes)
	j = &JWT{
		signingKey: []byte(opts.SigningKey),
		expiresAt:  time.Now().Add(d * time.Minute),
	}
	return j
}
func resolveJWT() *JWT {
	return j
}

type claims struct {
	J []byte
	jwt.RegisteredClaims
}

func (j *JWT) GenerateToken(payload map[string]interface{}) (string, error) {
	claims, err := mapClaims(payload, j.expiresAt)
	if err != nil {
		return "", err
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString(j.signingKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (j *JWT) DecodeToken(token string) (payload map[string]interface{}, err error) {
	t, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.signingKey, nil
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

func (j *JWT) HasExpired(token string) (bool, error) {
	_, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.signingKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return true, nil
		}
		return true, err
	}
	return false, nil
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
