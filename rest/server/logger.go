package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/log"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func newGinLogger(c *gin.Context) *logrus.Entry {
	return log.Get().
		WithField(log.LOG_FIELD_CORRELATION_ID, getCorrelationId(c)).
		WithField(log.LOG_FIELD_CONTOROLLER, "Rest").
		WithField(log.LOG_FIELD_HTTP_METHOD, c.Request.Method).
		WithField(log.LOG_FIELD_HTTP_PATH, c.Request.URL.Path)
}

func GinLoggerMiddleware(c *gin.Context) {
	logger := newGinLogger(c)

	c.Set("logger", logger)

	c.Next()

	if c.Request.Method != "OPTIONS" {
		ctx := GinCtx(c)
		log.Get(ctx...).WithField(log.LOG_FIELD_HTTP_STATUS, c.Writer.Status()).Info("Completed")
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
	value := c.GetHeader(log.LOG_FIELD_CORRELATION_ID)

	if len(value) == 0 {
		value = uuid.NewV4().String()
	}

	return value
}
