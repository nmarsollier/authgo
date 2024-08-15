package user

import (
	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/tools/db"
)

// Mock Data
func TestUser() (*User, string) {
	password := "123"
	userData := NewUser()
	userData.Login = "Login"
	userData.Name = "Name"
	userData.Permissions = []string{"user", "other"}
	userData.SetPasswordText(password)
	return userData, password
}

func TestAdminUser() (*User, string) {
	password := "123"
	userData := NewUser()
	userData.Login = "Login"
	userData.Name = "Name"
	userData.Permissions = []string{"user", "admin"}
	userData.SetPasswordText(password)
	return userData, password
}

func ExpectUserFindOne(coll *db.MockMongoCollection, userData *User, times int) {
	coll.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, params DbUserIdFilter, update *User) error {
			*update = *userData
			return nil
		},
	).Times(times)
}
