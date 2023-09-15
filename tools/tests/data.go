package tests

import (
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func TestToken() (*token.Token, string) {
	tokenData := &token.Token{
		ID:      primitive.NewObjectID(),
		UserID:  primitive.NewObjectID(),
		Enabled: true,
	}

	tokenString, _ := token.Encode(tokenData)

	return tokenData, tokenString
}
