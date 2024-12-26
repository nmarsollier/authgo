package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/internal/graph/tools"
)

func Revoke(ctx context.Context, userID string, permissions []string) (bool, error) {
	if err := tools.ValidateAdmin(ctx); err != nil {
		return false, err
	}

	di := tools.GqlDi(ctx)

	if err := di.UserService().Revoke(userID, permissions); err != nil {
		return false, err
	}

	return true, nil
}
