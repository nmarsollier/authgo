package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/graph/tools"
)

func Disable(ctx context.Context, userID string) (bool, error) {
	if err := tools.ValidateAdmin(ctx); err != nil {
		return false, err
	}

	di := tools.GqlDi(ctx)

	if err := di.UserService().Disable(userID); err != nil {
		return false, err
	}

	return true, nil
}
