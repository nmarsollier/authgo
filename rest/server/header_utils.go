package server

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/engine/errs"
	"github.com/nmarsollier/authgo/token"
)

// HeaderAuthorization token string from Authorization header
func HeaderAuthorization(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")

	if strings.Index(strings.ToUpper(tokenString), "BEARER ") == 0 {
		return tokenString[7:], nil
	}

	return "", errs.Unauthorized
}

// HeaderToken Token data from Authorization header
func HeaderToken(c *gin.Context) *token.Token {
	return c.MustGet("auth_header").(*token.Token)
}

func fetchAuthHeader(c *gin.Context) (*token.Token, error) {
	di := GinDi(c)
	tokenString, err := HeaderAuthorization(c)
	if err != nil {
		di.Logger().Error(err)
		return nil, err
	}

	payload, err := di.TokenService().Validate(tokenString)
	if err != nil {
		di.Logger().Error(err)
		return nil, err
	}

	c.Set("auth_header", payload)

	return payload, nil
}
