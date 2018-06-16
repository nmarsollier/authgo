package errors

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/core/topology"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/nmarsollier/authgo/tools/db"
	validator "gopkg.in/go-playground/validator.v9"
)

// Handle maneja cualquier error para serializarlo como JSON al cliente
/**
 * @apiDefine OtherErrors
 *
 * @apiErrorExample {json} 500 Server Error
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
		handleCustom(c, Internal)
		return
	case mongo.ErrNoDocuments:
		handleCustom(c, NotFound)
		return
	}

	// Compruebo tipos de errores conocidos
	switch value := err.(type) {
	case Custom:
		// Son validaciones hechas con NewCustom
		handleCustom(c, value)
	case Validation:
		// Son validaciones hechas con NewValidation
		c.JSON(400, err)
	case validator.ValidationErrors:
		// Son las validaciones de validator.v9 usadas en validaciones de estructuras
		handleValidationError(c, value)
	case mongo.WriteErrors:
		// Errores de mongo
		if db.IsUniqueKeyError(value) {
			handleCustom(c, AlreadyExist)
		} else {
			log.Output(1, fmt.Sprintf("Error DB : %s", value.Error()))
			handleCustom(c, Internal)
		}
	case error:
		// Otros errores
		c.JSON(500, gin.H{
			"error": value.Error(),
		})
	default:
		// No se sabe que es, devolvemos internal
		handleCustom(c, Internal)
	}
}

/**
 * @apiDefine ParamValidationErrors
 *
 * @apiErrorExample 400 Bad Request
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
	err := NewValidation()

	for _, e := range validationErrors {
		err.Add(strings.ToLower(e.Field()), e.Tag())
	}

	c.JSON(400, err)
}

func handleCustom(c *gin.Context, err Custom) {
	c.JSON(err.Status(), err)
}
