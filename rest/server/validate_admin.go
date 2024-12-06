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

	deps := GinDeps(c)
	c.Set("logger", log.Get(deps...).WithField(log.LOG_FIELD_USER_ID, payload.UserID))
	deps = GinDeps(c)

	if !user.Granted(payload.UserID, "admin", deps...) {
		log.Get(deps...).Warn("Unauthorized")
		c.Error(errs.Unauthorized)
		c.Abort()
	}
}
