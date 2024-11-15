package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/graph/tools"
	"github.com/nmarsollier/authgo/user"
)

func FindUser(ctx context.Context, id string) (*user.UserData, error) {
	/*_, err := tools.HeaderToken(ctx)
	if err != nil {
		return nil, err
	}*/

	env := tools.GqlCtx(ctx)
	user, err := user.Get(id, env...)
	if err != nil {
		return nil, err
	}

	return user, nil
}
