package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/graph/model"
	"github.com/nmarsollier/authgo/graph/tools"
)

func CurrentUser(ctx context.Context) (*model.User, error) {
	token, err := tools.HeaderToken(ctx)
	if err != nil {
		return nil, err
	}

	di := tools.GqlDi(ctx)
	user, err := di.UserService().FindById(token.UserID.Hex())
	if err != nil {
		return nil, err
	}

	return ToUser(user), nil
}
