package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/internal/graph/model"
	"github.com/nmarsollier/authgo/internal/graph/tools"
)

func FindUserByID(ctx context.Context, id string) (*model.User, error) {
	di := tools.GqlDi(ctx)
	user, err := di.UserService().FindById(id)
	if err != nil {
		return nil, err
	}

	return ToUser(user), nil
}
