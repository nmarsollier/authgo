package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/engine/errs"
	"github.com/nmarsollier/authgo/engine/log"
)

// Gin middleware to validate logged in user token
func IsAuthenticatedMiddleware(c *gin.Context) {
	token, err := fetchAuthHeader(c)
	if err != nil {
		c.Error(errs.Unauthorized)
		c.Abort()
		return
	}

	GinDi(c).Logger().WithField(log.LOG_FIELD_USER_ID, token.UserID.Hex())
}
