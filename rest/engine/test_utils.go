package engine

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Obtiene Router engine con el contexto de testing adecuado
// mockeando interfaces a serivcios externos
func TestRouter(ctx ...interface{}) *gin.Engine {
	engine = nil
	Router()
	if len(ctx) > 0 {
		engine.Use(func(c *gin.Context) {
			c.Set("mock_ctx", ctx)
			c.Next()
		})
	}
	fmt.Println("Paso 0")

	return engine
}

// Obtiene el contexto de interfaces mockeadas a serivcios externos
// En prod este contexto esta vacio.
func TestCtx(c *gin.Context) []interface{} {
	var ctx []interface{}
	fmt.Println("Paso 1")
	if mocks, ok := c.Get("mock_ctx"); ok {
		ctx = mocks.([]interface{})
	}
	return ctx
}
