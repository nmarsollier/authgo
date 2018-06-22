package test

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/core/topology"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/nmarsollier/authgo/mocks"
	handler "github.com/nmarsollier/authgo/tools/errors"
	validator "gopkg.in/go-playground/validator.v9"
)

func TestHandleErrID(t *testing.T) {
	response := mocks.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)
	handler.Handle(context, handler.ErrID)
	response.Assert(400, "{\"messages\":[{\"path\":\"id\",\"message\":\"Invalid\"}]}")
}

func TestHandleUnauthorized(t *testing.T) {
	response := mocks.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)
	handler.Handle(context, handler.Unauthorized)
	response.Assert(401, "{\"error\":\"Unauthorized\"}")
}

func TestHandleAccessLevel(t *testing.T) {
	response := mocks.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)
	handler.Handle(context, handler.AccessLevel)
	response.Assert(401, "{\"error\":\"Accesos Insuficientes\"}")
}

func TestHandleNotFound(t *testing.T) {
	response := mocks.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)
	handler.Handle(context, handler.NotFound)
	response.Assert(400, "{\"error\":\"Document not found\"}")
}

func TestHandleAlreadyExist(t *testing.T) {
	response := mocks.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)
	handler.Handle(context, handler.AlreadyExist)
	response.Assert(400, "{\"error\":\"Already exist\"}")
}

func TestHandleInternal(t *testing.T) {
	response := mocks.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)
	handler.Handle(context, handler.Internal)
	response.Assert(500, "{\"error\":\"Internal server error\"}")
}

func TestHandleNewValidation(t *testing.T) {
	validation := handler.NewValidation()
	validation.Add("f1", "Ef1")
	validation.Add("f2", "Ef2")

	response := mocks.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)
	handler.Handle(context, validation)
	response.Assert(400, "{\"messages\":[{\"path\":\"f1\",\"message\":\"Ef1\"},{\"path\":\"f2\",\"message\":\"Ef2\"}]}")
}

func TestHandleErrServerSelectionTimeout(t *testing.T) {
	response := mocks.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)
	handler.Handle(context, topology.ErrServerSelectionTimeout)
	response.Assert(500, "{\"error\":\"Internal server error\"}")
}

func TestHandleErrNoDocuments(t *testing.T) {
	response := mocks.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)
	handler.Handle(context, mongo.ErrNoDocuments)
	response.Assert(400, "{\"error\":\"Document not found\"}")
}

func TestHandleValidationError(t *testing.T) {
	type validStruct struct {
		OkField  string
		Required string `validate:"required"`
		Min      string `validate:"min=5"`
		Max      string `validate:"max=1"`
	}

	e := &validStruct{
		Min: "a",
		Max: "ab",
	}

	err := validator.New().Struct(e)

	response := mocks.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)
	handler.Handle(context, err)
	response.Assert(400, "{\"messages\":[{\"path\":\"required\",\"message\":\"required\"},{\"path\":\"min\",\"message\":\"min\"},{\"path\":\"max\",\"message\":\"max\"}]}")
}

func TestHandleWriteErrorsUnique(t *testing.T) {
	we := mongo.WriteErrors{
		mongo.WriteError{
			Index:   1,
			Code:    11000,
			Message: "Index",
		},
	}

	response := mocks.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)
	handler.Handle(context, we)
	response.Assert(400, "{\"error\":\"Already exist\"}")
}

func TestHandleWriteErrorsOther(t *testing.T) {
	we := mongo.WriteErrors{
		mongo.WriteError{
			Index:   1,
			Code:    11001,
			Message: "Other",
		},
	}

	response := mocks.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)
	handler.Handle(context, we)
	response.Assert(500, "{\"error\":\"Internal server error\"}")
}

func TestHandleError(t *testing.T) {
	response := mocks.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)
	handler.Handle(context, errors.New("Test"))
	response.Assert(500, "{\"error\":\"Test\"}")
}

func TestHandleNotError(t *testing.T) {
	response := mocks.ResponseWriter(t)
	context, _ := gin.CreateTestContext(response)
	handler.Handle(context, "Test")
	response.Assert(500, "{\"error\":\"Internal server error\"}")
}
