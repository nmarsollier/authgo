package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/common/errs"
	"github.com/nmarsollier/authgo/internal/common/log"
)

// Gin middleware to validate logged in user token
func IsAuthenticatedMiddleware(c *gin.Context) {
	token, err := loadTokenFromHeader(c)
	if err != nil {
		c.Error(errs.Unauthorized)
		c.Abort()
		return
	}

	GinLogger(c).WithField(log.LOG_FIELD_USER_ID, token.UserID.Hex())
}
