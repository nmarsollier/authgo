package errors

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type accessLevel struct {
}

// AccessLevel es el error de seguridad, el usuario no esta autorizado para acceder al recurso
var AccessLevel = &accessLevel{}

func (e *accessLevel) Error() string {
	return fmt.Sprintf("Accesos Insuficientes")
}

// Handle define como se serializa el error como Json
func (e *accessLevel) Handle(c *gin.Context) {
	c.JSON(401, gin.H{
		"error": "Accesos Insuficientes",
	})
}
