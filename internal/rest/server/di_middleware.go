package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/di"
	"github.com/nmarsollier/authgo/internal/env"
	"github.com/nmarsollier/commongo/log"
	uuid "github.com/satori/go.uuid"
)

func DiInjectorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var deps di.Injector
		dep_param, exists := c.Get("di")

		if !exists {
			logger := ginLogger(c)
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

func ginLogger(c *gin.Context) log.LogRusEntry {
	return log.Get(env.Get().FluentURL, "authgo").
		WithField(log.LOG_FIELD_CORRELATION_ID, getCorrelationId(c)).
		WithField(log.LOG_FIELD_CONTROLLER, "Rest").
		WithField(log.LOG_FIELD_HTTP_METHOD, c.Request.Method).
		WithField(log.LOG_FIELD_HTTP_PATH, c.Request.URL.Path)
}

func getCorrelationId(c *gin.Context) string {
	value := c.GetHeader(log.LOG_FIELD_CORRELATION_ID)

	if len(value) == 0 {
		value = uuid.NewV4().String()
	}

	return value
}
