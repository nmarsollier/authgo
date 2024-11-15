package graph

import (
	"context"

	"github.com/nmarsollier/authgo/graph/tools"
	"github.com/nmarsollier/authgo/token"
)

func signOutResolver(ctx context.Context) (bool, error) {
	tokenString, err := tools.TokenString(ctx)
	if err != nil {
		return false, err
	}

	env := tools.GqlCtx(ctx)

	if err := token.Invalidate(tokenString, env...); err != nil {
		return false, err
	}

	return true, nil
}
