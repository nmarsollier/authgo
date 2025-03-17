package server

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/common/errs"
	"github.com/nmarsollier/authgo/internal/common/log"
	uuid "github.com/satori/go.uuid"
)

// get token from Authorization header
func GetHeaderToken(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")
	if strings.Index(strings.ToUpper(tokenString), "BEARER ") != 0 {
		return "", errs.Unauthorized
	}
	return tokenString[7:], nil
}

func NewLogger(c *gin.Context, fluentUrl string, serverName string) log.LogRusEntry {
	return log.Get(fluentUrl, serverName).
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

func AbortWithError(c *gin.Context, err error) {
	c.Error(err)
	c.Abort()
}
