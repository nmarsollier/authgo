package engine

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/app_errors"
	"github.com/nmarsollier/authgo/user"
)

// ValidateAdmin check admin user is logged in
func ValidateAdmin(c *gin.Context) {
	var extraParams []interface{}
	if mocks, ok := c.Get("mocks"); ok {
		extraParams = mocks.([]interface{})
	}

	payload, err := fetchAuthHeader(c)
	if err != nil {
		c.Error(app_errors.Unauthorized)
		c.Abort()
		return
	}

	if !user.Granted(payload.UserID.Hex(), "admin", extraParams...) {
		c.Error(app_errors.Unauthorized)
		c.Abort()
	}
}

// ValidateLoggedIn user
func ValidateLoggedIn(c *gin.Context) {
	_, err := fetchAuthHeader(c)
	if err != nil {
		c.Error(app_errors.Unauthorized)
		c.Abort()
	}
}

// HeaderToken token from Authorization header
func HeaderToken(c *gin.Context) *token.Token {
	return c.MustGet("auth_header").(*token.Token)
}

func fetchAuthHeader(c *gin.Context) (*token.Token, error) {
	var extraParams []interface{}
	if mocks, ok := c.Get("mocks"); ok {
		extraParams = mocks.([]interface{})
	}

	tokenString, err := RAWHeaderToken(c)
	if err != nil {
		return nil, err
	}

	payload, err := token.Validate(tokenString, extraParams...)
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
		return "", app_errors.Unauthorized
	}
	return tokenString[7:], nil
}
