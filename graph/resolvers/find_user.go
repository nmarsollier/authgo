package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/graph/model"
	"github.com/nmarsollier/authgo/graph/tools"
)

func FindUserByID(ctx context.Context, id string) (*model.User, error) {
	di := tools.GqlDi(ctx)
	user, err := di.UserService().FindById(id)
	if err != nil {
		return nil, err
	}

	return ToUser(user), nil
}
