package usecases

import (
	"github.com/nmarsollier/authgo/internal/engine/errs"
	"github.com/nmarsollier/authgo/internal/token"
	"github.com/nmarsollier/authgo/internal/user"
)

type SignUpUseCase interface {
	SignUp(request *SignUpRequest) (*TokenResponse, error)
}

func NewSignUpUseCase(userService user.UserService, tokenService token.TokenService) SignUpUseCase {
	return &signUpUseCase{
		userService:  userService,
		tokenService: tokenService,
	}
}

type signUpUseCase struct {
	userService  user.UserService
	tokenService token.TokenService
}

type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Login    string `json:"login" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (s *signUpUseCase) SignUp(request *SignUpRequest) (*TokenResponse, error) {
	user, err := s.userService.New(request.Login, request.Name, request.Password)
	if err != nil {
		return nil, err
	}

	newToken, err := s.tokenService.Create(user.Id)

	if err != nil {
		return nil, errs.Internal
	}

	tokenString, err := token.Encode(newToken)
	if err != nil {
		return nil, errs.Unauthorized
	}

	return &TokenResponse{Token: tokenString}, nil
}
