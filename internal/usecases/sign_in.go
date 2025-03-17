package usecases

import (
	"github.com/nmarsollier/authgo/internal/common/errs"
	"github.com/nmarsollier/authgo/internal/common/log"
	"github.com/nmarsollier/authgo/internal/token"
	"github.com/nmarsollier/authgo/internal/user"
)

type SignInRequest struct {
	Password string `json:"password" binding:"required"`
	Login    string `json:"login" binding:"required"`
}

func SignIn(
	log log.LogRusEntry,
	request *SignInRequest,
) (*TokenResponse, error) {
	user, err := user.SignIn(log, request.Login, request.Password)
	if err != nil {
		return nil, err
	}

	newToken, err := token.Create(log, user.Id)
	if err != nil {
		return nil, errs.Unauthorized
	}

	tokenString, err := token.Encode(newToken)
	if err != nil {
		return nil, errs.Unauthorized
	}

	return &TokenResponse{Token: tokenString}, nil
}
