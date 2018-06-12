package user

import "github.com/nmarsollier/authgo/tools/errors"

// ErrID el id de usuario es invalido
var ErrID = errors.ErrValidation("id", "Invalid")

// ErrLogin el login es invalido
var ErrLogin = errors.ErrValidation("login", "Invalid")

// ErrLoginExist el login ya existe
var ErrLoginExist = errors.ErrValidation("login", "Ya existe")

// ErrPassword el password es invalido
var ErrPassword = errors.ErrValidation("password", "Invalid")
