package user

import "github.com/nmarsollier/authgo/tools/errors"

// ErrLogin el login es invalido
var ErrLogin = errors.NewValidationField("login", "invalid")

// ErrLoginExist el login ya existe
var ErrLoginExist = errors.NewValidationField("login", "exist")

// ErrPassword el password es invalido
var ErrPassword = errors.NewValidationField("password", "invalid")
