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

	env := tools.GqlCtx(ctx)
	if err := user.ChangePassword(token.UserID.Hex(), oldPassword, newPassword, env...); err != nil {
		return false, err
	}

	return true, nil
}
