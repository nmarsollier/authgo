package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/graph/tools"
	"github.com/nmarsollier/authgo/user"
)

func ChangePassword(ctx context.Context, oldPassword string, newPassword string) (bool, error) {
	token, err := tools.HeaderToken(ctx)
	if err != nil {
		return false, err
	}

	env := tools.GqlDeps(ctx)
	if err := user.ChangePassword(token.UserID, oldPassword, newPassword, env...); err != nil {
		return false, err
	}

	return true, nil
}
