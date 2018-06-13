package errors

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/core/topology"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/nmarsollier/authgo/tools/db"
	validator "gopkg.in/go-playground/validator.v8"
)

var alreadyExistError = gin.H{"error": "Already exist"}
var internalServerError = gin.H{"error": "Internal server error"}

// Handled es un error que se serializa utilizando Handle
// Si se desea definir un error personalizado como se serializa
// Ver unauthorized
type Handled interface {
	Handle(c *gin.Context)
}

// ErrInvalidField crea un error de validación para un solo campo
func ErrInvalidField(field string, err string) error {
	result := make(validator.ValidationErrors)

	result[field] = &validator.FieldError{
		Field: field,
		Tag:   err,
	}

	return result
}

// Handle maneja cualquier error para serializarlo como JSON al cliente
/**
 * @apiDefine OtherErrors
 *
 * @apiSuccessExample {json} 404 Not Found
 *     HTTP/1.1 404 Not Found
 *     {
 *        "url" : "{Url no encontrada}",
 *        "error" : "Not Found"
 *     }
 *
 * @apiSuccessExample {json} 500 Server Error
 *     HTTP/1.1 500 Internal Server Error
 *     {
 *        "error" : "Not Found"
 *     }
 *
 */
func Handle(c *gin.Context, err interface{}) {
	// Compruebo errores bien conocidos
	switch err {
	case topology.ErrServerSelectionTimeout, topology.ErrTopologyClosed:
		// Errores de conexión con MongoDB
		db.CheckError(err)
		c.JSON(500, internalServerError)
		return
	case mongo.ErrNoDocuments:
		c.JSON(400, gin.H{
			"error": "Document not found",
		})
		return
	}

	// Compruebo tipos de errores conocidos
	switch value := err.(type) {
	case validator.ValidationErrors:
		handleValidationError(c, value)
	case Handled:
		value.Handle(c)
	case mongo.WriteErrors:
		if IsUniqueKeyError(value) {
			c.JSON(400, alreadyExistError)
		} else {
			log.Output(1, fmt.Sprintf("Error DB : %s", value.Error()))
			c.JSON(500, internalServerError)
		}
	case error:
		c.JSON(500, gin.H{
			"error": value.Error(),
		})
	default:
		c.JSON(500, internalServerError)
	}
}

// IsUniqueKeyError retorna true si el error es de indice único
func IsUniqueKeyError(err error) bool {
	if wErr, ok := err.(mongo.WriteErrors); ok {
		for i := 0; i < len(wErr); i++ {
			if wErr[i].Code == 11000 {
				return true
			}
		}
	}
	return false
}

/**
 * @apiDefine ParamValidationErrors
 *
 * @apiSuccessExample {json} 400 Bad Request
 *     HTTP/1.1 400 Bad Request
 *     {
 *        "messages" : [
 *          {
 *            "path" : "{Nombre de la propiedad}",
 *            "message" : "{Motivo del error}"
 *          },
 *          ...
 *       ]
 *     }
 */
func handleValidationError(c *gin.Context, validationErrors validator.ValidationErrors) {
	var result []gin.H

	for _, err := range validationErrors {
		result = append(result, gin.H{
			"path":    strings.ToLower(err.Field),
			"message": err.Tag,
		})
	}

	c.JSON(400, gin.H{
		"messages": result,
	})
}
