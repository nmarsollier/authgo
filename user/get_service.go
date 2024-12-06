package user

import (
	"github.com/nmarsollier/authgo/tools/errs"
)

// Get wrapper para obtener un usuario
func Get(userID string, deps ...interface{}) (*UserData, error) {
	user, err := findByID(userID, deps...)
	if err != nil {
		return nil, err
	}

	if !user.Enabled {
		return nil, errs.NotFound
	}

	return &UserData{
		Id:          user.ID,
		Name:        user.Name,
		Permissions: user.Permissions,
		Login:       user.Login,
		Enabled:     user.Enabled,
	}, err
}
