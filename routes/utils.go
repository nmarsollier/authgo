package routes

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/tools/errors"
)

// get token from Authorization header
func getTokenHeader(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")
	if strings.Index(tokenString, "bearer ") != 0 {
		return "", errors.Unauthorized
	}
	return tokenString[7:], nil
}
