package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/model/token"
	"github.com/nmarsollier/authgo/model/user"
	"github.com/nmarsollier/authgo/tools/errors"
)

// ValidateAdmin check admin user is logged in
func ValidateAdmin(c *gin.Context) {
	payload, err := fetchAuthHeader(c)
	if err != nil {
		c.Error(errors.Unauthorized)
		c.Abort()
	}

	if !user.Granted(payload.UserID.Hex(), "admin") {
		c.Error(errors.Unauthorized)
		c.Abort()
	}
}

// ValidateLoggedIn user
func ValidateLoggedIn(c *gin.Context) {
	_, err := fetchAuthHeader(c)
	if err != nil {
		c.Error(errors.Unauthorized)
		c.Abort()
	}
}

// HeaderToken token from Authorization header
func HeaderToken(c *gin.Context) *token.Token {
	return c.MustGet("auth_header").(*token.Token)
}

func fetchAuthHeader(c *gin.Context) (*token.Token, error) {
	tokenString, err := RAWHeaderToken(c)
	if err != nil {
		return nil, err
	}

	payload, err := token.Validate(tokenString)
	if err != nil {
		return nil, err
	}

	c.Set("auth_header", payload)

	return payload, nil
}

// RAWHeaderToken token from Authorization header
func RAWHeaderToken(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")
	if strings.Index(tokenString, "bearer ") != 0 {
		return "", errors.Unauthorized
	}
	return tokenString[7:], nil
}
