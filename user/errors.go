package user

import "github.com/nmarsollier/authgo/tools/errors"

// ErrLogin el login es invalido
var ErrLogin = errors.NewInvalidField("login", "Invalid")

// ErrLoginExist el login ya existe
var ErrLoginExist = errors.NewInvalidField("login", "Ya existe")

// ErrPassword el password es invalido
var ErrPassword = errors.NewInvalidField("password", "Invalid")
