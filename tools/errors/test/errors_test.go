package test

import (
	"testing"

	"github.com/gin-gonic/gin/json"
	"github.com/nmarsollier/authgo/tools/errors"
	"github.com/stretchr/testify/assert"
)

// Test unitario de constantes predefinidas
func TestConstants(t *testing.T) {

	jsonErr, _ := json.Marshal(errors.ErrID)
	assert.Equal(t, "{\"messages\":[{\"path\":\"id\",\"message\":\"Invalid\"}]}", string(jsonErr))

	jsonErr, _ = json.Marshal(errors.Unauthorized)
	assert.Equal(t, "{\"error\":\"Unauthorized\"}", string(jsonErr))

	jsonErr, _ = json.Marshal(errors.AccessLevel)
	assert.Equal(t, "{\"error\":\"Accesos Insuficientes\"}", string(jsonErr))

	jsonErr, _ = json.Marshal(errors.NotFound)
	assert.Equal(t, "{\"error\":\"Document not found\"}", string(jsonErr))

	jsonErr, _ = json.Marshal(errors.AlreadyExist)
	assert.Equal(t, "{\"error\":\"Already exist\"}", string(jsonErr))

	jsonErr, _ = json.Marshal(errors.Internal)
	assert.Equal(t, "{\"error\":\"Internal server error\"}", string(jsonErr))
}

func TestNewValidationField(t *testing.T) {
	jsonErr, _ := json.Marshal(errors.NewValidationField("test", "Error Text"))
	assert.Equal(t, "{\"messages\":[{\"path\":\"test\",\"message\":\"Error Text\"}]}", string(jsonErr))

	validation := errors.NewValidationField("f1", "Ef1")
	validation.Add("f2", "Ef2")
	jsonErr, _ = json.Marshal(validation)
	assert.Equal(t, "{\"messages\":[{\"path\":\"f1\",\"message\":\"Ef1\"},{\"path\":\"f2\",\"message\":\"Ef2\"}]}", string(jsonErr))
}

func TestNewValidation(t *testing.T) {
	jsonErr, _ := json.Marshal(errors.NewValidation())
	assert.Equal(t, "{\"messages\":[]}", string(jsonErr))

	validation := errors.NewValidation()
	validation.Add("f1", "Ef1")
	jsonErr, _ = json.Marshal(validation)
	assert.Equal(t, validation.Size(), 1)
	assert.Equal(t, "{\"messages\":[{\"path\":\"f1\",\"message\":\"Ef1\"}]}", string(jsonErr))

	validation = errors.NewValidation()
	validation.Add("f1", "Ef1")
	validation.Add("f2", "Ef2")
	jsonErr, _ = json.Marshal(validation)
	assert.Equal(t, validation.Size(), 2)
	assert.Equal(t, "{\"messages\":[{\"path\":\"f1\",\"message\":\"Ef1\"},{\"path\":\"f2\",\"message\":\"Ef2\"}]}", string(jsonErr))
}

func TestNewCustom(t *testing.T) {
	err := errors.NewCustom(400, "Err400")
	jsonErr, _ := json.Marshal(err)
	assert.Equal(t, err.Status(), 400)
	assert.Equal(t, "{\"error\":\"Err400\"}", string(jsonErr))

	err = errors.NewCustom(100, "Another")
	jsonErr, _ = json.Marshal(err)
	assert.Equal(t, err.Status(), 100)
	assert.Equal(t, "{\"error\":\"Another\"}", string(jsonErr))
}
