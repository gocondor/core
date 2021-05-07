package sessions

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

type Sessions struct{}

var ses *Sessions
var sessionsOn bool

// New initiate sessions var
func New(sessionsFeatureOn bool) *Sessions {
	sessionsOn = sessionsFeatureOn

	ses = &Sessions{}
	return ses
}

// InitiateMemStore initiates memstore store of sessions
func (s *Sessions) InitiateMemstoreStore(secret string, name string) gin.HandlerFunc {
	store := memstore.NewStore([]byte(secret))
	return sessions.Sessions(name, store)
}

// InitiateCookieStore initiates cookie store of sessions
func (s *Sessions) InitiateCookieStore(secret string, name string) gin.HandlerFunc {
	store := cookie.NewStore([]byte(secret))
	return sessions.Sessions(name, store)
}

// InitiateRedistore initiates redis store of sessions
func (s *Sessions) InitiateRedistore(secret string, name string) gin.HandlerFunc {
	host := os.Getenv("localhost")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	store, _ := redis.NewStore(10, "tcp", fmt.Sprintf("%s:%s", host, port), password, []byte(secret))
	return sessions.Sessions(name, store)
}

// Resolve return sessions var
func Resolve() *Sessions {
	return ses
}

// Set sets key, val in session
func (s *Sessions) Set(key interface{}, val interface{}, c *gin.Context) {
	if !sessionsOn {
		exitWithlog()
	}
	session := sessions.Default(c)
	session.Set(key, val)
	session.Save()
}

// Get retrieves the key's value
func (s *Sessions) Get(key interface{}, c *gin.Context) interface{} {
	if !sessionsOn {
		exitWithlog()
	}
	session := sessions.Default(c)
	return session.Get(key)
}

// Get retrieves the key's value
func (s *Sessions) Has(key interface{}, c *gin.Context) bool {
	if !sessionsOn {
		exitWithlog()
	}
	session := sessions.Default(c)
	res := session.Get(key)
	if res == nil {
		return false
	}

	return true
}

// Get retrieves the key's value and delete it
func (s *Sessions) Pull(key interface{}, c *gin.Context) interface{} {
	if !sessionsOn {
		exitWithlog()
	}
	session := sessions.Default(c)
	val := session.Get(key)
	session.Delete(key)
	session.Save()
	return val
}

// Delete deletes session entry matching the given key
func (s *Sessions) Delete(key interface{}, c *gin.Context) {
	if !sessionsOn {
		exitWithlog()
	}
	session := sessions.Default(c)
	session.Delete(key)
	session.Save()
}

// Clear deletes all values in the session.
func (s *Sessions) Clear(c *gin.Context) {
	if !sessionsOn {
		exitWithlog()
	}
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}

func exitWithlog() {
	log.Fatal("please turn on the sessions feature before using it.")
}
