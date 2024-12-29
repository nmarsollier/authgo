package rest

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/internal/di"
	"github.com/nmarsollier/authgo/internal/rest/server"
	"github.com/nmarsollier/commongo/rst"
)

func TestRouter(ctrl *gomock.Controller, deps di.Injector) *gin.Engine {
	engine := gin.Default()
	engine.Use(func(c *gin.Context) {
		c.Set("di", deps)
		c.Next()
	})

	engine.Use(gzip.Gzip(gzip.DefaultCompression))
	engine.Use(server.DiInjectorMiddleware())
	engine.Use(rst.ErrorHandler)

	return engine
}
