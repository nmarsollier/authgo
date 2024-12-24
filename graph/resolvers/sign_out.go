package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/graph/tools"
)

func SignOut(ctx context.Context) (bool, error) {
	tokenString, err := tools.TokenString(ctx)
	if err != nil {
		return false, err
	}

	env := tools.GqlDi(ctx)

	if err := env.InvalidateTokenUseCase().InvalidateToken(tokenString); err != nil {
		return false, err
	}

	return true, nil
}
