package server

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/tools/log"
)

// HeaderAuthorization token string from Authorization header
func HeaderAuthorization(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")

	if strings.Index(tokenString, "bearer ") == 0 {
		return tokenString[7:], nil
	}

	if strings.Index(tokenString, "Bearer ") == 0 {
		return tokenString[7:], nil
	}

	return "", errs.Unauthorized
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
