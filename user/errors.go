package user

import "github.com/nmarsollier/authgo/tools/app_errors"

// ErrLogin el login es invalido
var ErrLogin = app_errors.NewValidationField("login", "invalid")

// ErrLoginExist el login ya existe
var ErrLoginExist = app_errors.NewValidationField("login", "exist")

// ErrPassword el password es invalido
var ErrPassword = app_errors.NewValidationField("password", "invalid")
