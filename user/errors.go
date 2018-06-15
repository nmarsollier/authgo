package user

import "github.com/nmarsollier/authgo/tools/errors"

// ErrLogin el login es invalido
var ErrLogin = errors.NewValidationField("login", "Invalid")

// ErrLoginExist el login ya existe
var ErrLoginExist = errors.NewValidationField("login", "Ya existe")

// ErrPassword el password es invalido
var ErrPassword = errors.NewValidationField("password", "Invalid")
