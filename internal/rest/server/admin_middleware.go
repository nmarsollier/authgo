package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/common/errs"
	"github.com/nmarsollier/authgo/internal/common/log"
	"github.com/nmarsollier/authgo/internal/user"
)

// Gin middleware to validate user token and Admin Access
func IsAdminMiddleware(c *gin.Context) {
	payload, err := loadTokenFromHeader(c)

	if err != nil {
		c.Error(errs.Unauthorized)
		c.Abort()
		return
	}

	log := GinLogger(c).WithField(log.LOG_FIELD_USER_ID, payload.UserID.Hex())

	if !user.Granted(log, payload.UserID.Hex(), "admin") {
		log.Warn("Unauthorized")
		c.Error(errs.Unauthorized)
		c.Abort()
	}
}
