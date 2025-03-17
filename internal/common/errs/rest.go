package errs

// Unauthorized is a predefined error for unauthorized access (HTTP 401).
var Unauthorized = NewRestError(401, "Unauthorized")

// NotFound is a predefined error for a document not found (HTTP 404).
var NotFound = NewRestError(404, "Document not found")

// AlreadyExist is a predefined error for a resource that already exists (HTTP 400).
var AlreadyExist = NewRestError(400, "Already exist")

// Internal is a predefined error for internal server errors (HTTP 500).
var Internal = NewRestError(500, "Internal server error")

var Invalid = NewRestError(400, "Invalid Document")

// NewRestError creates a new RestError with the given status code and message.
func NewRestError(status int, message string) RestError {
	return &restError{
		status:  status,
		Message: message,
	}
}

// RestError is an interface that defines methods for RESTful errors.
type RestError interface {
	Status() int
	Error() string
}

// restError is a struct that implements the RestError interface.
type restError struct {
	status  int
	Message string `json:"error"`
}

// Error returns the error message.
func (e *restError) Error() string {
	return e.Message
}

// Status returns the HTTP status code.
func (e *restError) Status() int {
	return e.status
}
