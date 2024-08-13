package user

import "github.com/nmarsollier/authgo/tools/apperr"

// ErrLogin el login es invalido
var ErrLogin = apperr.NewValidationField("login", "invalid")

// ErrLoginExist el login ya existe
var ErrLoginExist = apperr.NewValidationField("login", "exist")

// ErrPassword el password es invalido
var ErrPassword = apperr.NewValidationField("password", "invalid")
