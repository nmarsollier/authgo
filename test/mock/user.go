package mock

import (
	"github.com/nmarsollier/authgo/internal/user"
)

// Mock Data
func TestUser() (*user.User, string) {
	password := "123"
	userData := user.NewUser()
	userData.Login = "Login"
	userData.Name = "Name"
	userData.Permissions = []string{"user", "other"}
	userData.SetPasswordText(password)
	return userData, password
}

func TestAdminUser() (*user.User, string) {
	password := "123"
	userData := user.NewUser()
	userData.Login = "Login"
	userData.Name = "Name"
	userData.Permissions = []string{"user", "admin"}
	userData.SetPasswordText(password)
	return userData, password
}
