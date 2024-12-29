package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/di"
	"github.com/nmarsollier/authgo/internal/env"
	"github.com/nmarsollier/commongo/log"
	"github.com/nmarsollier/commongo/rst"
)

func DiInjectorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var deps di.Injector
		dep_param, exists := c.Get("di")

		if !exists {
			logger := rst.GinLogger(c, env.Get().FluentURL, "authgo")
			deps = di.NewInjector(logger)
			c.Set("di", deps)
		} else {
			deps = dep_param.(di.Injector)
		}

		c.Next()

		if c.Request.Method != "OPTIONS" {
			deps.Logger().WithField(log.LOG_FIELD_HTTP_STATUS, c.Writer.Status()).Info("Completed")
		}
	}
}

func GinDi(c *gin.Context) di.Injector {
	return c.MustGet("di").(di.Injector)
}
