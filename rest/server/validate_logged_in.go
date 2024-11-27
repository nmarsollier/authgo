package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/tools/log"
)

// Gin middleware to validate logged in user token
func ValidateLoggedIn(c *gin.Context) {
	token, err := fetchAuthHeader(c)
	if err != nil {
		c.Error(errs.Unauthorized)
		c.Abort()
		return
	}

	deps := GinDeps(c)
	c.Set("logger", log.Get(deps...).WithField(log.LOG_FIELD_USER_ID, token.UserID.Hex()))
}
