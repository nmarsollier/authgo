package resolvers

import (
	"github.com/nmarsollier/authgo/internal/graph/model"
	"github.com/nmarsollier/authgo/internal/user"
)

func ToUser(user *user.UserData) (result *model.User) {
	return &model.User{
		ID:          user.Id,
		Name:        user.Name,
		Permissions: user.Permissions,
		Login:       user.Login,
		Enabled:     user.Enabled,
	}
}
