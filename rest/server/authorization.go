package server

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/log"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/errs"
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

// Gin middleware to validate logged in user token
func ValidateLoggedIn(c *gin.Context) {
	token, err := fetchAuthHeader(c)
	if err != nil {
		c.Error(errs.Unauthorized)
		c.Abort()
		return
	}

	ctx := GinCtx(c)
	c.Set("logger", log.Get(ctx...).WithField(log.LOG_FIELD_USER_ID, token.UserID.Hex()))
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
	ctx := GinCtx(c)
	tokenString, err := HeaderAuthorization(c)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	payload, err := token.Validate(tokenString, ctx...)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	c.Set("auth_header", payload)

	return payload, nil
}
