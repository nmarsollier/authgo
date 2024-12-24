package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/graph/model"
	"github.com/nmarsollier/authgo/graph/tools"
	"github.com/nmarsollier/authgo/usecases"
)

func SignUp(ctx context.Context, name string, login string, password string) (*model.Token, error) {
	env := tools.GqlDi(ctx)
	token, err := env.SignUpUseCase().SignUp(&usecases.SignUpRequest{
		Name:     name,
		Login:    login,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return &model.Token{
		Token: token.Token,
	}, nil
}
