package errors

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type unauthorized struct {
}

// Unauthorized es el error de seguridad, el usuario no esta autorizado para acceder al recurso
var Unauthorized = &unauthorized{}

func (e *unauthorized) Error() string {
	return fmt.Sprintf("Unauthorized")
}

// Handle es un error que se serializa como Json
func (e *unauthorized) Handle(c *gin.Context) {
	c.JSON(401, gin.H{
		"error": "Unauthorized",
	})
}
