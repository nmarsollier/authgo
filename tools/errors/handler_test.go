package errors

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/core/topology"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/nmarsollier/authgo/test"
	validator "gopkg.in/go-playground/validator.v9"
)

func TestHandleErrID(t *testing.T) {
	response := test.NewFakeHttpResponse(t)
	context, _ := gin.CreateTestContext(response)
	Handle(context, ErrID)
	response.Assert(400, "{\"messages\":[{\"path\":\"id\",\"message\":\"Invalid\"}]}")
}

func TestHandleUnauthorized(t *testing.T) {
	response := test.NewFakeHttpResponse(t)
	context, _ := gin.CreateTestContext(response)
	Handle(context, Unauthorized)
	response.Assert(401, "{\"error\":\"Unauthorized\"}")
}

func TestHandleAccessLevel(t *testing.T) {
	response := test.NewFakeHttpResponse(t)
	context, _ := gin.CreateTestContext(response)
	Handle(context, AccessLevel)
	response.Assert(401, "{\"error\":\"Accesos Insuficientes\"}")
}

func TestHandleNotFound(t *testing.T) {
	response := test.NewFakeHttpResponse(t)
	context, _ := gin.CreateTestContext(response)
	Handle(context, NotFound)
	response.Assert(400, "{\"error\":\"Document not found\"}")
}

func TestHandleAlreadyExist(t *testing.T) {
	response := test.NewFakeHttpResponse(t)
	context, _ := gin.CreateTestContext(response)
	Handle(context, AlreadyExist)
	response.Assert(400, "{\"error\":\"Already exist\"}")
}

func TestHandleInternal(t *testing.T) {
	response := test.NewFakeHttpResponse(t)
	context, _ := gin.CreateTestContext(response)
	Handle(context, Internal)
	response.Assert(500, "{\"error\":\"Internal server error\"}")
}

func TestHandleNewValidation(t *testing.T) {
	validation := NewValidation()
	validation.Add("f1", "Ef1")
	validation.Add("f2", "Ef2")

	response := test.NewFakeHttpResponse(t)
	context, _ := gin.CreateTestContext(response)
	Handle(context, validation)
	response.Assert(400, "{\"messages\":[{\"path\":\"f1\",\"message\":\"Ef1\"},{\"path\":\"f2\",\"message\":\"Ef2\"}]}")
}

func TestHandleErrServerSelectionTimeout(t *testing.T) {
	response := test.NewFakeHttpResponse(t)
	context, _ := gin.CreateTestContext(response)
	Handle(context, topology.ErrServerSelectionTimeout)
	response.Assert(500, "{\"error\":\"Internal server error\"}")
}

func TestHandleErrNoDocuments(t *testing.T) {
	response := test.NewFakeHttpResponse(t)
	context, _ := gin.CreateTestContext(response)
	Handle(context, mongo.ErrNoDocuments)
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

	response := test.NewFakeHttpResponse(t)
	context, _ := gin.CreateTestContext(response)
	Handle(context, err)
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

	response := test.NewFakeHttpResponse(t)
	context, _ := gin.CreateTestContext(response)
	Handle(context, we)
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

	response := test.NewFakeHttpResponse(t)
	context, _ := gin.CreateTestContext(response)
	Handle(context, we)
	response.Assert(500, "{\"error\":\"Internal server error\"}")
}

func TestHandleError(t *testing.T) {
	response := test.NewFakeHttpResponse(t)
	context, _ := gin.CreateTestContext(response)
	Handle(context, errors.New("Test"))
	response.Assert(500, "{\"error\":\"Test\"}")
}

func TestHandleNotError(t *testing.T) {
	response := test.NewFakeHttpResponse(t)
	context, _ := gin.CreateTestContext(response)
	Handle(context, "Test")
	response.Assert(500, "{\"error\":\"Internal server error\"}")
}
