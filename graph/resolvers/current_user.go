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

	deps := tools.GqlDeps(ctx)
	user, err := user.Get(token.UserID, deps...)
	if err != nil {
		return nil, err
	}

	return ToUser(user), nil
}
