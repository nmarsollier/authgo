package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/graph/tools"
	"github.com/nmarsollier/authgo/token"
)

func SignOut(ctx context.Context) (bool, error) {
	tokenString, err := tools.TokenString(ctx)
	if err != nil {
		return false, err
	}

	env := tools.GqlDeps(ctx)

	if err := token.Invalidate(tokenString, env...); err != nil {
		return false, err
	}

	return true, nil
}
