package user

import (
	"github.com/nmarsollier/authgo/tools/errs"
)

// Get wrapper para obtener un usuario
func Get(userID string, ctx ...interface{}) (*User, error) {
	user, err := findByID(userID, ctx...)
	if err != nil {
		return nil, err
	}

	if !user.Enabled {
		return nil, errs.NotFound
	}

	return user, err
}
