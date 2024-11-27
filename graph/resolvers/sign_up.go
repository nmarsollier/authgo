package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/graph/model"
	"github.com/nmarsollier/authgo/graph/tools"
	"github.com/nmarsollier/authgo/user"
)

func SignUp(ctx context.Context, name string, login string, password string) (*model.Token, error) {
	env := tools.GqlDeps(ctx)
	token, err := user.SignUp(&user.SignUpRequest{
		Name:     name,
		Login:    login,
		Password: password,
	}, env...)
	if err != nil {
		return nil, err
	}
	return &model.Token{
		Token: token.Token,
	}, nil
}
