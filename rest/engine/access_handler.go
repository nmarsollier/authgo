package engine

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/user"
)

// Gin middleware to validate user token and Admin Access
func ValidateAdmin(c *gin.Context) {
	ctx := TestCtx(c)

	payload, err := fetchAuthHeader(c)
	if err != nil {
		c.Error(errs.Unauthorized)
		c.Abort()
		return
	}

	if !user.Granted(payload.UserID.Hex(), "admin", ctx...) {
		c.Error(errs.Unauthorized)
		c.Abort()
	}
}

// Gin middleware to validate logged in user token
func ValidateLoggedIn(c *gin.Context) {
	_, err := fetchAuthHeader(c)
	if err != nil {
		c.Error(errs.Unauthorized)
		c.Abort()
	}
}

// HeaderAuthorization token string from Authorization header
func HeaderAuthorization(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")
	if strings.Index(tokenString, "bearer ") != 0 {
		return "", errs.Unauthorized
	}
	return tokenString[7:], nil
}

// HeaderToken Token data from Authorization header
func HeaderToken(c *gin.Context) *token.Token {
	return c.MustGet("auth_header").(*token.Token)
}

func fetchAuthHeader(c *gin.Context) (*token.Token, error) {
	ctx := TestCtx(c)
	tokenString, err := HeaderAuthorization(c)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	payload, err := token.Validate(tokenString, ctx...)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	c.Set("auth_header", payload)

	return payload, nil
}
