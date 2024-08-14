package user

import "github.com/nmarsollier/authgo/tools/errs"

// ErrID el id del documento es invalido
var ErrID = errs.NewValidation().Add("id", "Invalid")

// ErrLogin el login es invalido
var ErrLogin = errs.NewValidation().Add("login", "invalid")

// ErrLoginExist el login ya existe
var ErrLoginExist = errs.NewValidation().Add("login", "exist")

// ErrPassword el password es invalido
var ErrPassword = errs.NewValidation().Add("password", "invalid")
