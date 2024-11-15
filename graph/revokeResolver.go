package graph

import (
	"context"

	"github.com/nmarsollier/authgo/graph/tools"
	"github.com/nmarsollier/authgo/user"
)

func revokeResolver(ctx context.Context, userID string, permissions []string) (bool, error) {
	if err := tools.ValidateAdmin(ctx); err != nil {
		return false, err
	}

	env := tools.GqlCtx(ctx)

	if err := user.Revoke(userID, permissions, env...); err != nil {
		return false, err
	}

	return true, nil
}
