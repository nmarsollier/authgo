package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/graph/model"
	"github.com/nmarsollier/authgo/graph/tools"
	"github.com/nmarsollier/authgo/user"
)

// SignIn is the resolver for the signIn field.
func SignIn(ctx context.Context, login string, password string) (*model.Token, error) {
	env := tools.GqlDeps(ctx)

	tokenString, err := user.SignIn(user.SignInRequest{Login: login, Password: password}, env...)
	if err != nil {
		return nil, err
	}

	return &model.Token{
		Token: tokenString.Token,
	}, nil
}
