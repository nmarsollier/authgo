package server

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nmarsollier/authgo/internal/engine/db"
	"github.com/nmarsollier/authgo/internal/engine/di"
	"github.com/nmarsollier/authgo/internal/engine/errs"
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

	if errors.Is(err, mongo.ErrNoDocuments) {
		setError(c, errs.NotFound)
		return
	}

	if errors.Is(err, topology.ErrServerSelectionTimeout) || errors.Is(err, topology.ErrTopologyClosed) {
		di.IsDbTimeoutError(err)
		setError(c, errs.Internal)
		return
	}

	switch value := err.Err.(type) {
	case errs.RestError:
		setError(c, value)
	case errs.Validation:
		c.JSON(400, err)
	case validator.ValidationErrors:
		handleValidationError(c, value)
	case mongo.WriteException:
		if db.IsDbUniqueKeyError(value) {
			setError(c, errs.AlreadyExist)
		} else {
			setError(c, errs.Internal)
		}
	case error:
		c.JSON(500, ErrorData{
			Error: value.Error(),
		})
	default:
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

type ErrorData struct {
	Error string `json:"error"`
}
