package user

import (
	"github.com/nmarsollier/authgo/tools/errs"
)

// Get wrapper para obtener un usuario
func Get(userID string, ctx ...interface{}) (*UserResponse, error) {
	user, err := findByID(userID, ctx...)
	if err != nil {
		return nil, err
	}

	if !user.Enabled {
		return nil, errs.NotFound
	}

	return &UserResponse{
		Id:          user.ID.Hex(),
		Name:        user.Name,
		Permissions: user.Permissions,
		Login:       user.Login,
		Enabled:     user.Enabled,
	}, err
}
