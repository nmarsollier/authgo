package apperr

import (
	"encoding/json"

	"github.com/golang/glog"
)

// Validation es una interfaz para definir errores custom
// Validation es un error de validaciones de parameteros o de campos
type Validation interface {
	Add(path string, message string) Validation
	Size() int
	Error() string
}

func NewValidation() Validation {
	return &validation{
		Messages: []errField{},
	}
}

type validation struct {
	Messages []errField `json:"messages"`
}

func (e *validation) Error() string {
	body, err := json.Marshal(e)
	if err != nil {
		glog.Error(err)
		return "ErrValidation invalid."
	}
	return string(body)
}

// Add agrega errores a un validation error
func (e *validation) Add(path string, message string) Validation {
	err := errField{
		Path:    path,
		Message: message,
	}
	e.Messages = append(e.Messages, err)
	return e
}

// Size devuelve la cantidad de errores
func (e *validation) Size() int {
	return len(e.Messages)
}

// errField define un campo inválido. path y mensaje de error
type errField struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}
