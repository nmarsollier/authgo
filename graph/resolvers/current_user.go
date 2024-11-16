package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/graph/model"
	"github.com/nmarsollier/authgo/graph/tools"
	"github.com/nmarsollier/authgo/user"
)

func CurrentUser(ctx context.Context) (*model.User, error) {
	token, err := tools.HeaderToken(ctx)
	if err != nil {
		return nil, err
	}

	env := tools.GqlCtx(ctx)
	user, err := user.Get(token.UserID.Hex(), env...)
	if err != nil {
		return nil, err
	}

	return ToUser(user), nil
}