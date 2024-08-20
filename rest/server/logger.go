package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/log"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func newGinLogger(c *gin.Context) *logrus.Entry {
	return log.Get().
		WithField("CorrelationId", getCorrelationId(c)).
		WithField("Controller", "Rest").
		WithField("Method", c.Request.Method).
		WithField("Path", c.Request.URL.Path)
}

func GinLoggerMiddleware(c *gin.Context) {
	logger := newGinLogger(c)

	c.Set("logger", logger)

	c.Next()

	if c.Request.Method != "OPTIONS" {
		logger.WithField("status", c.Writer.Status()).Info("Completed")
	}
}

func ginLogger(c *gin.Context) *logrus.Entry {
	logger, exist := c.Get("logger")
	if !exist {
		return newGinLogger(c)
	}
	return logger.(*logrus.Entry)
}

func getCorrelationId(c *gin.Context) string {
	value := c.GetHeader("CorrelationId")

	if len(value) == 0 {
		value = uuid.NewV4().String()
	}

	return value
}
