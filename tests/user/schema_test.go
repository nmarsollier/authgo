package user

import (
	"testing"

	"github.com/nmarsollier/authgo/user"
	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	user := user.NewUser()
	assert.Equal(t, user.Enabled, true)
	assert.NotEqual(t, user.ID.Hex(), "000000000000000000000000")
	assert.NotNil(t, user.Created)
	assert.NotNil(t, user.Updated)
	assert.Equal(t, user.Permissions[0], "user")

	user.SetPasswordText("password")
	assert.NotEmpty(t, user.Password)
	assert.NotEqual(t, user.Password, "password")
	assert.Nil(t, user.ValidatePassword("password"))

	user.Grant("otro")
	assert.Equal(t, user.Permissions[0], "user")
	assert.Equal(t, user.Permissions[1], "otro")
	assert.Equal(t, len(user.Permissions), 2)

	assert.Equal(t, user.Granted("otro"), true)
	user.Revoke("otro")
	assert.Equal(t, user.Permissions[0], "user")
	assert.Equal(t, len(user.Permissions), 1)
}
