package errors

import (
	"fmt"

	validator "gopkg.in/go-playground/validator.v8"
)

// - Algunos errors comunes en el sistema -

// ErrID el id de usuario es invalido
var ErrID = NewInvalidField("id", "Invalid")

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

// NewInvalidField crea un error de validación para un solo campo
func NewInvalidField(field string, err string) error {
	result := make(validator.ValidationErrors)

	result[field] = &validator.FieldError{
		Field: field,
		Tag:   err,
	}

	return result
}

// NewCustom creates a new errCustom
func NewCustom(status int, message string) *ErrCustom {
	return &ErrCustom{
		status:  status,
		message: message,
	}
}

//  - Algunas definiciones necesarias -

// Custom es una interfaz para definir errores custom
type Custom interface {
	Status() int
	Message() string
}

// ErrCustom es un error personalizado para http
type ErrCustom struct {
	status  int
	message string
}

func (e *ErrCustom) Error() string {
	return fmt.Sprintf(e.message)
}

// Status http status code
func (e *ErrCustom) Status() int {
	return e.status
}

// Message mensage de error
func (e *ErrCustom) Message() string {
	return e.message
}
