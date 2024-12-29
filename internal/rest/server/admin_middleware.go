package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/commongo/errs"
	"github.com/nmarsollier/commongo/log"
)

// Gin middleware to validate user token and Admin Access
func IsAdminMiddleware(c *gin.Context) {
	payload, err := loadTokenFromHeader(c)

	if err != nil {
		c.Error(errs.Unauthorized)
		c.Abort()
		return
	}

	di := GinDi(c)
	di.Logger().WithField(log.LOG_FIELD_USER_ID, payload.UserID.Hex())

	if !di.UserService().Granted(payload.UserID.Hex(), "admin") {
		di.Logger().Warn("Unauthorized")
		c.Error(errs.Unauthorized)
		c.Abort()
	}
}
