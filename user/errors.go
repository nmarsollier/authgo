package user

import "github.com/nmarsollier/authgo/tools/errors"

var InvalidUserIdError = errors.NewValidationErrorError("id", "Invalid")
var InvalidLoginError = errors.NewValidationErrorError("login", "Invalid")
var LoginAlreadyExistError = errors.NewValidationErrorError("login", "Ya existe")
var InvalidPasswordError = errors.NewValidationErrorError("password", "Invalid")
