package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/graph/tools"
)

func ChangePassword(ctx context.Context, oldPassword string, newPassword string) (bool, error) {
	token, err := tools.HeaderToken(ctx)
	if err != nil {
		return false, err
	}

	di := tools.GqlDi(ctx)
	if err := di.UserService().ChangePassword(token.UserID.Hex(), oldPassword, newPassword); err != nil {
		return false, err
	}

	return true, nil
}
