package middlewares_test

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/gocondor/core/middlewares"
)

func TestNew(t *testing.T) {
	m := New()
	if fmt.Sprintf("%T", m) != "*middlewares.MiddlewaresUtil" {
		t.Errorf("Failed asserting middleware util var initiation")
	}
}

func TestResolve(t *testing.T) {
	m := Resolve()
	if fmt.Sprintf("%T", m) != "*middlewares.MiddlewaresUtil" {
		t.Errorf("Failed asserting middleware util var resolve")
	}
}

func TestAttach(t *testing.T) {
	m := New()
	f := func(c *gin.Context) {
	}
	m.Attach(f)
	if len(m.GetMiddlewares()) != 1 {
		t.Error("failed attach middleware")
	}
}

func TestGetMiddlewares(t *testing.T) {
	m := New()
	middlewares := []gin.HandlerFunc{
		func(c *gin.Context) {},
		func(c *gin.Context) {},
		func(c *gin.Context) {},
	}
	for _, val := range middlewares {
		m.Attach(val)
	}
	if len(m.GetMiddlewares()) != 3 {
		t.Error("failed to get attached middlewares")
	}
}
