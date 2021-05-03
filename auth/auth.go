package auth

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gocondor/core/jwt"
	"github.com/gocondor/core/sessions"
)

type Auth struct {
	Ses *sessions.Sessions
	JWT *jwt.JWTUtil
}

var auth *Auth

const USER_ID = "userId"

func New(ses *sessions.Sessions, jwt *jwt.JWTUtil) *Auth {
	auth = &Auth{
		ses,
		jwt,
	}

	return auth
}

func Resolve() *Auth {
	return auth
}

func (a *Auth) Login(userId uint, c *gin.Context) error {
	// store the user in the session
	a.Ses.Set(USER_ID, userId, c)

	return nil
}

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
