package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/graph/tools"
	"github.com/nmarsollier/authgo/user"
)

func SignUp(ctx context.Context, name string, login string, password string) (*user.TokenResponse, error) {
	env := tools.GqlCtx(ctx)
	token, err := user.SignUp(&user.SignUpRequest{
		Name:     name,
		Login:    login,
		Password: password,
	}, env...)
	if err != nil {
		return nil, err
	}
	return token, nil
}
