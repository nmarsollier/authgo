package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/graph/tools"
)

func Grant(ctx context.Context, userID string, permissions []string) (bool, error) {
	if err := tools.ValidateAdmin(ctx); err != nil {
		return false, err
	}

	di := tools.GqlDi(ctx)

	if err := di.UserService().Grant(userID, permissions); err != nil {
		return false, err
	}

	return true, nil
}
