package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/tools/log"
	"github.com/nmarsollier/authgo/user"
)

// Gin middleware to validate user token and Admin Access
func ValidateAdmin(c *gin.Context) {
	payload, err := fetchAuthHeader(c)

	if err != nil {
		c.Error(errs.Unauthorized)
		c.Abort()
		return
	}

	ctx := GinCtx(c)
	c.Set("logger", log.Get(ctx...).WithField(log.LOG_FIELD_USER_ID, payload.UserID.Hex()))
	ctx = GinCtx(c)

	if !user.Granted(payload.UserID.Hex(), "admin", ctx...) {
		log.Get(ctx...).Warn("Unauthorized")
		c.Error(errs.Unauthorized)
		c.Abort()
	}
}
