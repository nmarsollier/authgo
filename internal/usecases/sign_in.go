package usecases

import (
	"github.com/nmarsollier/authgo/internal/token"
	"github.com/nmarsollier/authgo/internal/user"
	"github.com/nmarsollier/commongo/errs"
)

type SignInUseCase interface {
	SignIn(request *SignInRequest) (*TokenResponse, error)
}

func NewSignInUseCase(userService user.UserService, tokenService token.TokenService) SignInUseCase {
	return &signInUseCase{
		userService:  userService,
		tokenService: tokenService,
	}
}

type signInUseCase struct {
	userService  user.UserService
	tokenService token.TokenService
}

type SignInRequest struct {
	Password string `json:"password" binding:"required"`
	Login    string `json:"login" binding:"required"`
}

func (s *signInUseCase) SignIn(request *SignInRequest) (*TokenResponse, error) {
	user, err := s.userService.SignIn(request.Login, request.Password)
	if err != nil {
		return nil, err
	}

	newToken, err := s.tokenService.Create(user.Id)
	if err != nil {
		return nil, errs.Unauthorized
	}

	tokenString, err := token.Encode(newToken)
	if err != nil {
		return nil, errs.Unauthorized
	}

	return &TokenResponse{Token: tokenString}, nil
}
