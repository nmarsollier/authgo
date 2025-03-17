package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/common/log"
	"github.com/nmarsollier/authgo/internal/env"
)

func DiInjectorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var logger log.LogRusEntry
		logger_ref, exists := c.Get("logger")

		if !exists {
			logger = NewLogger(c, env.Get().FluentURL, env.Get().ServerName)
			c.Set("logger", logger)
		} else {
			logger = logger_ref.(log.LogRusEntry)
		}

		c.Next()

		if c.Request.Method != "OPTIONS" {
			logger.WithField(log.LOG_FIELD_HTTP_STATUS, c.Writer.Status()).Info("Completed")
		}
	}
}

func GinLogger(c *gin.Context) log.LogRusEntry {
	return c.MustGet("logger").(log.LogRusEntry)
}
