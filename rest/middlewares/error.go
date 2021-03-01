package middlewares

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

/**
 * @apiDefine AuthHeader
 *
 * @apiExample {String} Header Autorización
 *    Authorization=bearer {token}
 *
 * @apiErrorExample 401 Unauthorized
 *    HTTP/1.1 401 Unauthorized
 *    {
 *       "error" : "Unauthorized"
 *    }
 */

// ErrorHandler a middleware to handle errors
func ErrorHandler(c *gin.Context) {
	c.Next()

	handleErrorIfNeeded(c)
}

func AbortWithError(c *gin.Context, err error) {
	c.Error(err)
	c.Abort()
}

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
func handleErrorIfNeeded(c *gin.Context) {
	err := c.Errors.Last()
	if err == nil {
		return
	}

	// Compruebo errores bien conocidos
	switch err {
	case topology.ErrServerSelectionTimeout, topology.ErrTopologyClosed:
		// Errores de conexión con MongoDB
		db.CheckError(err)
		handleCustom(c, errors.Internal)
		return
	case mongo.ErrNoDocuments:
		handleCustom(c, errors.NotFound)
		return
	}

	// Compruebo tipos de errores conocidos
	switch value := err.Err.(type) {
	case errors.Custom:
		// Son validaciones hechas con NewCustom
		handleCustom(c, value)
	case errors.Validation:
		// Son validaciones hechas con NewValidation
		c.JSON(400, err)
	case validator.ValidationErrors:
		// Son las validaciones de validator usadas en validaciones de estructuras
		handleValidationError(c, value)
	case mongo.WriteException:
		// Errores de mongo
		if db.IsUniqueKeyError(value) {
			handleCustom(c, errors.AlreadyExist)
		} else {
			log.Output(1, fmt.Sprintf("Error DB : %s", value.Error()))
			handleCustom(c, errors.Internal)
		}
	case error:
		// Otros errores
		c.JSON(500, gin.H{
			"error": value.Error(),
		})
	default:
		// No se sabe que es, devolvemos internal
		handleCustom(c, errors.Internal)
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
	err := errors.NewValidation()

	for _, e := range validationErrors {
		err.Add(strings.ToLower(e.Field()), e.Tag())
	}

	c.JSON(400, err)
}

func handleCustom(c *gin.Context, err errors.Custom) {
	c.JSON(err.Status(), err)
}
