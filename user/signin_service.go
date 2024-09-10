package user

import (
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/errs"
)

type SignInRequest struct {
	Password string `json:"password" binding:"required"`
	Login    string `json:"login" binding:"required"`
}

// SignIn is the controller to sign in users
func SignIn(data SignInRequest, ctx ...interface{}) (string, error) {
	user, err := findByLogin(data.Login, ctx...)
	if err != nil {
		return "", err
	}

	if !user.Enabled {
		return "", errs.Unauthorized
	}

	if err = user.validatePassword(data.Password); err != nil {
		return "", err
	}

	newToken, err := token.Create(user.ID, ctx...)
	if err != nil {
		return "", errs.Unauthorized
	}

	return token.Encode(newToken)
}
