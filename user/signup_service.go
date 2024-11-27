package user

import (
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/errs"
)

// SignUpRequest es un nuevo usuario
type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Login    string `json:"login" binding:"required"`
}

// SignUp is the controller to signup new users
func SignUp(user *SignUpRequest, deps ...interface{}) (*TokenResponse, error) {
	newUser := NewUser()
	newUser.Login = user.Login
	newUser.Name = user.Name
	newUser.setPasswordText(user.Password)

	newUser, err := insert(newUser, deps...)
	if err != nil {
		return nil, err
	}

	newToken, err := token.Create(newUser.ID, deps...)
	if err != nil {
		return nil, errs.Internal
	}

	tokenString, err := token.Encode(newToken)
	if err != nil {
		return nil, errs.Unauthorized
	}

	return &TokenResponse{Token: tokenString}, nil
}
