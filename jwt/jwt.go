// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package jwt

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTUtil pakcage struct
type JWTUtil struct{}

// DefaultTokenLifeSpan is the default ttl for jwt token
var DefaultTokenLifeSpan time.Duration = 15 * time.Minute //15 minutes

// DefaultRefreshTokenLifeSpanHours is the default ttl for jwt refresh token
var DefaultRefreshTokenLifeSpanHours time.Duration = 24 * time.Hour //24 hours

var USER_ID = "userID"

var JWT *JWTUtil

// New initiates Jwt struct
func New() *JWTUtil {
	JWT = &JWTUtil{}
	return JWT
}

//Resolve returns initiated jwt token
func Resolve() *JWTUtil {
	return JWT
}

// CreateToken generates new jwt token with the given user id
func (j *JWTUtil) CreateToken(userID uint) (string, error) {

	claims := jwt.MapClaims{}

	var duration time.Duration
	durationStr := os.Getenv("JWT_LIFESPAN_MINUTES")
	if durationStr == "" {
		duration = DefaultTokenLifeSpan
	} else {
		d, _ := strconv.ParseInt(durationStr, 10, 64)
		duration = time.Duration(d) * time.Minute
	}

	claims[USER_ID] = userID
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(duration).Unix()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("missing jwt token secret")
	}
	token, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

// CreateRefreshToken generates new jwt refresh token with the given user id
func (j *JWTUtil) CreateRefreshToken(userID uint) (string, error) {
	claims := jwt.MapClaims{}

	var duration time.Duration
	durationStrHours := os.Getenv("JWT_REFRESH_TOKEN_LIFESPAN_HOURS")
	if durationStrHours == "" {
		duration = DefaultRefreshTokenLifeSpanHours
	} else {
		d, _ := strconv.ParseInt(durationStrHours, 10, 64)
		duration = time.Duration(d) * time.Hour
	}
	claims[USER_ID] = userID
	claims["exp"] = time.Now().Add(duration).Unix()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_REFRESH_TOKEN_SECRET")
	if secret == "" {
		return "", errors.New("missing jwt token refresh secret")
	}
	token, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

//ExtractToken extracts the token from the request header
func (j *JWTUtil) ExtractToken(c *gin.Context) (token string, err error) {
	sentTokenSlice := c.Request.Header["Authorization"]
	if len(sentTokenSlice) == 0 {
		return "", errors.New("Missing authorization token")
	}
	sentTokenSlice = strings.Split(sentTokenSlice[0], " ")
	if len(sentTokenSlice) != 2 {
		return "", errors.New("Something wrong with the token")
	}

	return sentTokenSlice[1], nil
}

// DecodeToken decodes a given token and returns the user id
func (j *JWTUtil) DecodeToken(tokenString string) (userID uint, err error) {
	// validate the token
	_, err = j.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}

	//extract claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	claims := token.Claims.(jwt.MapClaims)
	delete(claims, "authorized")
	delete(claims, "exp")

	id, err := strconv.ParseInt(fmt.Sprintf("%v", claims[USER_ID]), 10, 32)
	if err != nil {
		return 0, err
	}
	userID = uint(id)
	return userID, nil
}

// ValidateToken makes sure the given token is valid
func (j *JWTUtil) ValidateToken(tokenString string) (bool, error) {
	// parse the token string
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %s", token.Method.Alg())
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

// RefreshToken generates a new token based on the refresh token
// TODO: implement
func RefreshToken(token string, refreshToken string) (newToken string, err error) {
	return
}
