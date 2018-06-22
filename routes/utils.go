package routes

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/security"
	"github.com/nmarsollier/authgo/tools/errors"
)

// get token from Authorization header
func getAuthHeader(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")
	if strings.Index(tokenString, "bearer ") != 0 {
		return "", errors.Unauthorized
	}
	return tokenString[7:], nil
}

func validateAuthHeader(c *gin.Context) (*security.Token, error) {
	tokenString, err := getAuthHeader(c)
	if err != nil {
		return nil, err
	}

	payloadRepository, err := security.NewService()
	if err != nil {
		return nil, err
	}

	payload, err := payloadRepository.Validate(tokenString)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
