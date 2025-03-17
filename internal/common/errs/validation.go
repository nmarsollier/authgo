package errs

import (
	"encoding/json"
)

// Validation represents an interface for handling validation errors.
// It provides methods to add validation errors and retrieve the error message.
//
// Add adds a validation error with the specified path and message.
// It returns the updated Validation instance.
//
// Error returns a string representation of the validation errors.
type Validation interface {
	Add(path string, message string) Validation
	Error() string
}

// NewValidation creates a new instance of ValidationErr with an empty list of error messages.
// It returns the Validation interface implemented by ValidationErr.
func NewValidation() Validation {
	return &ValidationErr{
		Messages: []errField{},
	}
}

// ValidationErr represents a validation error that contains a list of error messages.
type ValidationErr struct {
	Messages []errField `json:"messages"`
}

// Error implements the error interface for ValidationErr.
func (e *ValidationErr) Error() string {
	body, err := json.Marshal(e)
	if err != nil {
		return "ErrValidation invalid."
	}
	return string(body)
}

// Add appends a new validation error message to the ValidationErr instance.
// It takes a path and a message as parameters, creates an errField with these values,
// and adds it to the Messages slice of the ValidationErr instance.
// It returns the updated ValidationErr instance.
//
// Parameters:
//   - path: The path or field name where the validation error occurred.
//   - message: The validation error message.
//
// Returns:
//   - Validation: The updated ValidationErr instance with the new error message added.
func (e *ValidationErr) Add(path string, message string) Validation {
	err := errField{
		Path:    path,
		Message: message,
	}
	e.Messages = append(e.Messages, err)
	return e
}

// errField represents a validation error for a specific field.
// Path indicates the location of the field that caused the error.
// Message provides a description of the validation error.
type errField struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}
