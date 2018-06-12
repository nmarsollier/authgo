package user

import (
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/errors"
)

// NewUser es un nuevo usuario
type NewUser struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Login    string `json:"login" binding:"required"`
}

// SignUp is the controller to signup new users
func SignUp(user *NewUser) (string, error) {
	newUser := newUser()
	newUser.Login = user.Login
	newUser.Name = user.Name
	newUser.Roles = []string{"user"}
	newUser.setPasswordText(user.Password)

	newUser, err := save(newUser)
	if err != nil {
		if errors.IsUniqueKeyError(err) {
			return "", ErrLoginExist
		} else {
			return "", err
		}
	}

	tokenString, err := token.Create(newUser.ID())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// SignIn is the controller to sign in users
func SignIn(login string, password string) (string, error) {
	user, err := findByLogin(login)
	if err != nil {
		return "", err
	}

	err = user.validatePassword(password)
	if err != nil {
		return "", err
	}

	tokenString, err := token.Create(user.ID())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// CurrentUser is the controller to get the current logged in user
func CurrentUser(userID string) (*User, error) {
	return findByID(userID)
}

// ChangePassword Change Password Controller
func ChangePassword(userID string, current string, newPassword string) error {
	user, err := findByID(userID)
	if err != nil {
		return err
	}

	err = user.validatePassword(current)
	if err != nil {
		return err
	}

	err = user.setPasswordText(newPassword)
	if err != nil {
		return err
	}

	_, err = save(*user)
	if err != nil {
		return err
	}

	return nil
}
