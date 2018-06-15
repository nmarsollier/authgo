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
		handleCustom(c, Internal)
		return
	case mongo.ErrNoDocuments:
		handleCustom(c, NotFound)
		return
	}

	// Compruebo tipos de errores conocidos
	switch value := err.(type) {
	case validator.ValidationErrors:
		handleValidationError(c, value)
	case Custom:
		handleCustom(c, value)
	case mongo.WriteErrors:
		if db.IsUniqueKeyError(value) {
			handleCustom(c, AlreadyExist)
		} else {
			log.Output(1, fmt.Sprintf("Error DB : %s", value.Error()))
			handleCustom(c, Internal)
		}
	case error:
		c.JSON(500, gin.H{
			"error": value.Error(),
		})
	default:
		handleCustom(c, Internal)
	}
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

func handleCustom(c *gin.Context, err Custom) {
	c.JSON(err.Status(), gin.H{
		"error": err.Message(),
	})
}
