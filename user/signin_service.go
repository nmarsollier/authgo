package user

import (
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/errs"
)

type SignInRequest struct {
	Password string `json:"password" binding:"required"`
	Login    string `json:"login" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// SignIn is the controller to sign in users
func SignIn(data SignInRequest, deps ...interface{}) (*TokenResponse, error) {
	user, err := findByLogin(data.Login, deps...)
	if err != nil {
		return nil, err
	}

	if !user.Enabled {
		return nil, errs.Unauthorized
	}

	if err = user.validatePassword(data.Password); err != nil {
		return nil, err
	}

	newToken, err := token.Create(user.ID, deps...)
	if err != nil {
		return nil, errs.Unauthorized
	}

	tokenString, err := token.Encode(newToken)
	if err != nil {
		return nil, errs.Unauthorized
	}

	return &TokenResponse{Token: tokenString}, nil
}
