package user

import "github.com/nmarsollier/authgo/tools/errors"

// ErrID el id de usuario es invalido
var ErrID = errors.ErrInvalidField("id", "Invalid")

// ErrLogin el login es invalido
var ErrLogin = errors.ErrInvalidField("login", "Invalid")

// ErrLoginExist el login ya existe
var ErrLoginExist = errors.ErrInvalidField("login", "Ya existe")

// ErrPassword el password es invalido
var ErrPassword = errors.ErrInvalidField("password", "Invalid")
