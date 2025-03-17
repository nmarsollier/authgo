package usecases

import (
	"github.com/nmarsollier/authgo/internal/common/log"
	"github.com/nmarsollier/authgo/internal/token"
	"github.com/nmarsollier/authgo/internal/user"
)

type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Login    string `json:"login" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func SignUp(
	log log.LogRusEntry,
	request *SignUpRequest,
) (*TokenResponse, error) {
	user, err := user.New(log, request.Login, request.Name, request.Password)
	if err != nil {
		return nil, err
	}

	newToken, err := token.Create(log, user.Id)
	if err != nil {
		return nil, err
	}

	tokenString, err := token.Encode(newToken)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{Token: tokenString}, nil
}
