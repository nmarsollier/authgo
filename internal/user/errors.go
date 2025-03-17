package user

import (
	"github.com/nmarsollier/authgo/internal/common/errs"
)

// ErrID el id del documento es invalido
var ErrID = errs.NewValidation().Add("id", "Invalid")

// ErrLogin el login es invalido
var ErrLogin = errs.NewValidation().Add("login", "invalid")

// ErrPassword el password es invalido
var ErrPassword = errs.NewValidation().Add("password", "invalid")
