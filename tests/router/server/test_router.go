package server

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	_ "github.com/nmarsollier/authgo/docs"
	"github.com/nmarsollier/authgo/engine/di"
	"github.com/nmarsollier/authgo/rest/server"
)

func TestRouter(ctrl *gomock.Controller, deps di.Injector) *gin.Engine {
	engine := gin.Default()
	engine.Use(func(c *gin.Context) {
		c.Set("di", deps)
		c.Next()
	})

	engine.Use(gzip.Gzip(gzip.DefaultCompression))
	engine.Use(server.DiInjectorMiddleware())
	engine.Use(server.ErrorHandler)

	return engine
}
