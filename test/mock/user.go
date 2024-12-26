package mock

import (
	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/internal/user"
	"github.com/nmarsollier/authgo/test/mockgen"
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

func ExpectUserFindOne(coll *mockgen.MockCollection, userData *user.User, times int) {
	coll.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, params user.DbUserIdFilter, update *user.User) error {
			*update = *userData
			return nil
		},
	).Times(times)
}
