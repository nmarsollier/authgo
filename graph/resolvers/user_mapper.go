package resolvers

import (
	"github.com/nmarsollier/authgo/graph/model"
	"github.com/nmarsollier/authgo/user"
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
