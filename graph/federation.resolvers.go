package graph

import (
	"context"

	"github.com/nmarsollier/authgo/graph/resolvers"
	"github.com/nmarsollier/authgo/user"
)

func (r *Resolver) FindUserByID(ctx context.Context, id string) (*user.UserResponse, error) {
	user, err := resolvers.FindUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
