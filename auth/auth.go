package auth

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gocondor/core/jwt"
	"github.com/gocondor/core/sessions"
)

// Auth is authentication management struct
type Auth struct {
	Ses *sessions.Sessions
	JWT *jwt.JWTUtil
}

var auth *Auth

// user id key
const USER_ID = "userId"

// New initiates new Auth
func New(ses *sessions.Sessions, jwt *jwt.JWTUtil) *Auth {
	auth = &Auth{
		ses,
		jwt,
	}

	return auth
}

// Resolve returns the initiated Auth variable
func Resolve() *Auth {
	return auth
}

// Login logs the user in by id
func (a *Auth) Login(userId uint, c *gin.Context) error {
	// store the user in the session
	a.Ses.Set(USER_ID, userId, c)

	return nil
}

// UserID returns authenticated user id
func (a *Auth) UserID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseInt(fmt.Sprintf("%v", a.Ses.Get(USER_ID, c)), 10, 32)
	if err != nil {
		return 0, err
	}
	userID := uint(id)

	return userID, nil
}

// Logout logs the user out by id
func (a *Auth) Logout(c *gin.Context) error {
	// delete the user id from the session
	a.Ses.Delete(USER_ID, c)

	return nil
}

// Check checks if a user is logged in
func (a *Auth) Check(c *gin.Context) (bool, error) {
	// if session doesn't have user id return false
	if !a.Ses.Has(USER_ID, c) {
		return false, nil
	}

	id, err := strconv.ParseInt(fmt.Sprintf("%v", a.Ses.Get(USER_ID, c)), 10, 32)
	if err != nil {
		return false, err
	}
	sessionUserId := uint(id)

	// extract the token
	token, err := a.JWT.ExtractToken(c)
	if err != nil {
		return false, err
	}
	// extract the token's user id
	tokenUserId, err := a.JWT.DecodeToken(token)
	if err != nil {
		return false, err
	}
	// if the session's user id matches the token user id he is authenticated
	if sessionUserId == tokenUserId {
		return true, nil
	}
	return false, nil
}
