package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/internal/graph/model"
	"github.com/nmarsollier/authgo/internal/graph/tools"
)

func FindAllUsers(ctx context.Context) ([]*model.User, error) {
	if err := tools.ValidateAdmin(ctx); err != nil {
		return nil, err
	}

	env := tools.GqlDi(ctx)
	users, err := env.UserService().FindAllUsers()

	if err != nil {
		return nil, err
	}

	result := make([]*model.User, len(users))
	for i := range users {
		result[i] = ToUser(users[i])
	}

	return result, nil
}
