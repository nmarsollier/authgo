package engine

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nmarsollier/authgo/tools/apperr"
	"github.com/nmarsollier/authgo/tools/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

func ErrorHandler(c *gin.Context) {
	c.Next()
	handleErrorIfNeeded(c)
}

func handleErrorIfNeeded(c *gin.Context) {
	err := c.Errors.Last()
	if err == nil {
		return
	}

	// Compruebo errores bien conocidos
	if errors.Is(err, mongo.ErrNoDocuments) {
		setError(c, apperr.NotFound)
		return
	}
	if errors.Is(err, topology.ErrServerSelectionTimeout) || errors.Is(err, topology.ErrTopologyClosed) {
		// Errores de conexión con MongoDB
		db.IsDbTimeoutError(err)
		setError(c, apperr.Internal)
		return
	}

	// Compruebo tipos de errores conocidos
	switch value := err.Err.(type) {
	case apperr.RestError:
		// Son validaciones hechas con NewCustom
		setError(c, value)
	case apperr.Validation:
		// Son validaciones hechas con NewValidation
		c.JSON(400, err)
	case validator.ValidationErrors:
		// Son las validaciones de validator usadas en validaciones de estructuras
		handleValidationError(c, value)
	case mongo.WriteException:
		// Errores de mongo
		if db.IsDbUniqueKeyError(value) {
			setError(c, apperr.AlreadyExist)
		} else {
			setError(c, apperr.Internal)
		}
	case error:
		// Otros errores
		c.JSON(500, ErrorData{
			Error: value.Error(),
		})
	default:
		// No se sabe que es, devolvemos internal
		setError(c, apperr.Internal)
	}
}

func handleValidationError(c *gin.Context, validationErrors validator.ValidationErrors) {
	err := apperr.NewValidation()

	for _, e := range validationErrors {
		err.Add(strings.ToLower(e.Field()), e.Tag())
	}

	c.JSON(400, err)
}

func setError(c *gin.Context, err apperr.RestError) {
	c.JSON(err.Status(), err)
}

type ErrorData struct {
	Error string `json:"error"`
}
