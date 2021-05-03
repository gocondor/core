package auth

import (
	"errors"
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

// Logout logs the user out by id
func (a *Auth) Logout(userId uint64, c *gin.Context) error {
	// get the user id from the session and convert it to uint64
	userIdS, err := strconv.ParseUint(fmt.Sprintf("%v", a.Ses.Get(USER_ID, c)), 10, 64)
	if err != nil {
		return err
	}
	// check if the user id in args matches the user id in the session
	if userId == userIdS {
		// delete the user id from the session
		a.Ses.Delete(USER_ID, c)
		return nil
	}

	return errors.New("trying to logout different user")
}

// Check checks if a user is logged in
func (a *Auth) Check(userId uint64, c *gin.Context) bool {
	// if session doesn't have user id return false
	if !a.Ses.Has(USER_ID, c) {
		return false
	}

	// get the user id from the session and convert it to uint64
	userIdS, err := strconv.ParseUint(fmt.Sprintf("%v", a.Ses.Get(USER_ID, c)), 10, 64)
	if err != nil {
		return false
	}

	// if the arg user id matches the session's user id he is authenticated
	if userId == userIdS {
		return true
	}

	// not authenticated
	return false
}
