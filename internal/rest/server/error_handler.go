package server

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nmarsollier/authgo/internal/common/errs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

// ErrorHandler a middleware to handle errors
func ErrorHandler(c *gin.Context) {
	c.Next()

	handleErrorIfNeeded(c)
}

func handleErrorIfNeeded(c *gin.Context) {
	err := c.Errors.Last()
	if err == nil {
		return
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		setError(c, errs.NotFound)
		return
	}

	if errors.Is(err, topology.ErrServerSelectionTimeout) || errors.Is(err, topology.ErrTopologyClosed) {
		setError(c, errs.Internal)
		return
	}

	handleErrorByType(c, err.Err)
}

// handleErrorByType handles any error to serialize it as JSON to the client
func handleErrorByType(c *gin.Context, err interface{}) {
	// Check for known error types
	switch value := err.(type) {
	case errs.RestError:
		// These are validations made with NewCustom
		setError(c, value)
	case errs.Validation:
		// These are validations made with NewValidation
		c.JSON(400, err)
	case mongo.WriteException:
		if IsDbUniqueKeyError(value) {
			setError(c, errs.AlreadyExist)
		} else {
			setError(c, errs.Internal)
		}
	case validator.ValidationErrors:
		// These are validator validations used in structure validations
		handleValidationError(c, value)
	case error:
		// Other errors
		c.JSON(500, ErrorData{
			Error: value.Error(),
		})
	default:
		// Unknown error type, return internal error
		setError(c, errs.Internal)
	}
}

func handleValidationError(c *gin.Context, validationErrors validator.ValidationErrors) {
	err := errs.NewValidation()

	for _, e := range validationErrors {
		err.Add(strings.ToLower(e.Field()), e.Tag())
	}

	c.JSON(400, err)
}

func setError(c *gin.Context, err errs.RestError) {
	c.JSON(err.Status(), err)
}

// IsDbUniqueKeyError retorna true si el error es de indice único
func IsDbUniqueKeyError(err error) bool {
	if wErr, ok := err.(mongo.WriteException); ok {
		for i := 0; i < len(wErr.WriteErrors); i++ {
			if wErr.WriteErrors[i].Code == 11000 {
				return true
			}
		}
	}
	return false
}

type ErrorData struct {
	Error string `json:"error"`
}
