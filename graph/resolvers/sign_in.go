package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/graph/model"
	"github.com/nmarsollier/authgo/graph/tools"
	"github.com/nmarsollier/authgo/usecases"
)

// SignIn is the resolver for the signIn field.
func SignIn(ctx context.Context, login string, password string) (*model.Token, error) {
	env := tools.GqlDi(ctx)

	tokenString, err := env.SignInUseCase().SignIn(&usecases.SignInRequest{Login: login, Password: password})
	if err != nil {
		return nil, err
	}

	return &model.Token{
		Token: tokenString.Token,
	}, nil
}
