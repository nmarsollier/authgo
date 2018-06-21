package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	user := newUser()
	assert.Equal(t, user.Enabled, true)
	assert.NotEqual(t, user.ID.Hex(), "000000000000000000000000")
	assert.NotNil(t, user.Created)
	assert.NotNil(t, user.Updated)
	assert.Equal(t, user.Permissions[0], "user")

	user.setPasswordText("password")
	assert.NotEmpty(t, user.Password)
	assert.NotEqual(t, user.Password, "password")
	assert.Nil(t, user.validatePassword("password"))

	user.grant("otro")
	assert.Equal(t, user.Permissions[0], "user")
	assert.Equal(t, user.Permissions[1], "otro")
	assert.Equal(t, len(user.Permissions), 2)

	assert.Equal(t, user.granted("otro"), true)
	user.revoke("otro")
	assert.Equal(t, user.Permissions[0], "user")
	assert.Equal(t, len(user.Permissions), 1)
}
