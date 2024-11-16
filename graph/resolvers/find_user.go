package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/graph/model"
	"github.com/nmarsollier/authgo/graph/tools"
	"github.com/nmarsollier/authgo/user"
)

func FindUserByID(ctx context.Context, id string) (*model.User, error) {
	/*_, err := tools.HeaderToken(ctx)
	if err != nil {
		return nil, err
	}*/

	env := tools.GqlCtx(ctx)
	user, err := user.Get(id, env...)
	if err != nil {
		return nil, err
	}

	return ToUser(user), nil
}
