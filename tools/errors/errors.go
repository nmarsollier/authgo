package errors

import (
	"encoding/json"
	"fmt"
)

// - Algunos errors comunes en el sistema -

// ErrID el id del documento es invalido
var ErrID = NewValidationField("id", "Invalid")

// Unauthorized el usuario no esta autorizado al recurso
var Unauthorized = NewCustom(401, "Unauthorized")

// AccessLevel es el error de seguridad, el usuario no esta autorizado para acceder al recurso
var AccessLevel = NewCustom(401, "Accesos Insuficientes")

// NotFound cuando un registro no se encuentra en la db
var NotFound = NewCustom(400, "Document not found")

// AlreadyExist cuando no se puede ingresar un registro a la db
var AlreadyExist = NewCustom(400, "Already exist")

// Internal esta aplicación no sabe como manejar el error
var Internal = NewCustom(500, "Internal server error")

// - Creación de errors -

// NewValidationField crea un error de validación para un solo campo
func NewValidationField(field string, err string) Validation {
	return &ErrValidation{
		Messages: []ErrField{
			{
				Path:    field,
				Message: err,
			},
		},
	}
}

// NewValidation crea un error de validación para un solo campo
func NewValidation() Validation {
	return &ErrValidation{
		Messages: []ErrField{},
	}
}

// NewCustom creates a new errCustom
func NewCustom(status int, message string) *ErrCustom {
	return &ErrCustom{
		status:  status,
		Message: message,
	}
}

//  - Algunas definiciones necesarias -

// Custom es una interfaz para definir errores custom
type Custom interface {
	Status() int
	Error() string
}

// ErrCustom es un error personalizado para http
type ErrCustom struct {
	status  int
	Message string `json:"error"`
}

func (e *ErrCustom) Error() string {
	return fmt.Sprintf(e.Message)
}

// Status http status code
func (e *ErrCustom) Status() int {
	return e.status
}

// Validation es una interfaz para definir errores custom
// Validation es un error de validaciones de parameteros o de campos
type Validation interface {
	Add(path string, message string) Validation
	Size() int
	Error() string
}

// ErrField define un campo inválido. path y mensaje de error
type ErrField struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

// ErrValidation es un error de validaciones de parameteros o de campos
type ErrValidation struct {
	Messages []ErrField `json:"messages"`
}

func (e *ErrValidation) Error() string {
	body, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ErrValidation que no se puede pasar a json.")
	}
	return fmt.Sprintf(string(body))
}

// Add agrega errores a un validation error
func (e *ErrValidation) Add(path string, message string) Validation {
	err := ErrField{
		Path:    path,
		Message: message,
	}
	e.Messages = append(e.Messages, err)
	return e
}

// Size devuelve la cantidad de errores
func (e *ErrValidation) Size() int {
	return len(e.Messages)
}
