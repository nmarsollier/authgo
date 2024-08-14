package errs

// Unauthorized el usuario no esta autorizado al recurso
var Unauthorized = NewRestError(401, "Unauthorized")

// NotFound cuando un registro no se encuentra en la db
var NotFound = NewRestError(404, "Document not found")

// AlreadyExist cuando no se puede ingresar un registro a la db
var AlreadyExist = NewRestError(400, "Already exist")

// Internal esta aplicación no sabe como manejar el error
var Internal = NewRestError(500, "Internal server error")

// - Creación de errors -
// NewRestError creates a new errCustom
func NewRestError(status int, message string) RestError {
	return &restError{
		status:  status,
		Message: message,
	}
}

//  - Algunas definiciones necesarias -

// RestError es una interfaz para definir errores custom
type RestError interface {
	Status() int
	Error() string
}

// restError es un error personalizado para http
type restError struct {
	status  int
	Message string `json:"error"`
}

func (e *restError) Error() string {
	return e.Message
}

// Status http status code
func (e *restError) Status() int {
	return e.status
}
