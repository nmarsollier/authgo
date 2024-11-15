package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/graph/tools"
	"github.com/nmarsollier/authgo/user"
)

func Users(ctx context.Context) ([]*user.UserResponse, error) {
	if err := tools.ValidateAdmin(ctx); err != nil {
		return nil, err
	}

	env := tools.GqlCtx(ctx)
	result, err := user.FindAllUsers(env...)

	if err != nil {
		return nil, err
	}

	return result, nil
}
