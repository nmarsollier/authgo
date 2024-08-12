package engine

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/nmarsollier/authgo/tools/app_errors"
	"github.com/nmarsollier/authgo/tools/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

func ErrorHandler(c *gin.Context) {
	c.Next()
	handleErrorIfNeeded(c)
}

func AbortWithError(c *gin.Context, err error) {
	c.Error(err)
	c.Abort()
}

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
		handleCustom(c, app_errors.Internal)
		return
	case mongo.ErrNoDocuments:
		handleCustom(c, app_errors.NotFound)
		return
	}

	// Compruebo tipos de errores conocidos
	switch value := err.Err.(type) {
	case app_errors.Custom:
		// Son validaciones hechas con NewCustom
		handleCustom(c, value)
	case app_errors.Validation:
		// Son validaciones hechas con NewValidation
		c.JSON(400, err)
	case validator.ValidationErrors:
		// Son las validaciones de validator usadas en validaciones de estructuras
		handleValidationError(c, value)
	case mongo.WriteException:
		// Errores de mongo
		if db.IsUniqueKeyError(value) {
			handleCustom(c, app_errors.AlreadyExist)
		} else {
			handleCustom(c, app_errors.Internal)
		}
	case error:
		// Otros errores
		c.JSON(500, app_errors.OtherErrors{
			Error: value.Error(),
		})
	default:
		// No se sabe que es, devolvemos internal
		handleCustom(c, app_errors.Internal)
	}
}

func handleValidationError(c *gin.Context, validationErrors validator.ValidationErrors) {
	err := app_errors.NewValidation()

	for _, e := range validationErrors {
		err.Add(strings.ToLower(e.Field()), e.Tag())
	}

	c.JSON(400, err)
}

func handleCustom(c *gin.Context, err app_errors.Custom) {
	c.JSON(err.Status(), err)
}
