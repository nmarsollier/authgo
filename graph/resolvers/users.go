package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/graph/model"
	"github.com/nmarsollier/authgo/graph/tools"
	"github.com/nmarsollier/authgo/user"
)

func Users(ctx context.Context) ([]*model.User, error) {
	if err := tools.ValidateAdmin(ctx); err != nil {
		return nil, err
	}

	env := tools.GqlDeps(ctx)
	users, err := user.FindAllUsers(env...)

	if err != nil {
		return nil, err
	}

	result := make([]*model.User, len(users))
	for i := range users {
		result[i] = ToUser(users[i])
	}

	return result, nil
}
