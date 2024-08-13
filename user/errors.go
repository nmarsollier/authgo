package user

import "github.com/nmarsollier/authgo/tools/apperr"

// ErrID el id del documento es invalido
var ErrID = apperr.NewValidation().Add("id", "Invalid")

// ErrLogin el login es invalido
var ErrLogin = apperr.NewValidation().Add("login", "invalid")

// ErrLoginExist el login ya existe
var ErrLoginExist = apperr.NewValidation().Add("login", "exist")

// ErrPassword el password es invalido
var ErrPassword = apperr.NewValidation().Add("password", "invalid")
