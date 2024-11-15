package graph

import (
	"context"

	"github.com/nmarsollier/authgo/graph/tools"
	"github.com/nmarsollier/authgo/user"
)

// SignIn is the resolver for the signIn field.
func signInResolver(ctx context.Context, login string, password string) (*user.TokenResponse, error) {
	env := tools.GqlCtx(ctx)

	tokenString, err := user.SignIn(user.SignInRequest{Login: login, Password: password}, env...)
	if err != nil {
		return nil, err
	}

	return tokenString, nil
}
