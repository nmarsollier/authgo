package user

import (
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/errors"
)

// SignUpRequest es un nuevo usuario
type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Login    string `json:"login" binding:"required"`
}

// SignUp is the controller to signup new users
func SignUp(user *SignUpRequest) (string, error) {
	newUser := newUser()
	newUser.Login = user.Login
	newUser.Name = user.Name
	newUser.setPasswordText(user.Password)

	newUser, err := insert(newUser)
	if err != nil {
		if errors.IsUniqueKeyError(err) {
			return "", ErrLoginExist
		}
		return "", err
	}

	tokenString, err := token.Create(newUser.ID)
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

	if !user.Enabled {
		return "", errors.Unauthorized
	}

	if err = user.validatePassword(password); err != nil {
		return "", err
	}

	tokenString, err := token.Create(user.ID)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// CurrentUser is the controller to get the current logged in user
func CurrentUser(userID string) (*User, error) {
	return findByID(userID)
}

// ChangePassword Permite cambiar contraseña
func ChangePassword(userID string, current string, newPassword string) error {
	user, err := findByID(userID)
	if err != nil {
		return err
	}

	if err = user.validatePassword(current); err != nil {
		return err
	}

	if err = user.setPasswordText(newPassword); err != nil {
		return err
	}

	if _, err = update(user); err != nil {
		return err
	}

	return nil
}

// Grant Le habilita los permisos enviados por parametros
func Grant(userID string, permissions []string) error {
	user, err := findByID(userID)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.Grant(value)
	}

	if _, err = update(user); err != nil {
		return err
	}

	return nil
}

// Revoke Le revoca los permisos enviados por parametros
func Revoke(userID string, permissions []string) error {
	user, err := findByID(userID)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.Revoke(value)
	}

	if _, err = update(user); err != nil {
		return err
	}

	return nil
}

//Granted verifica si el usuario tiene el permiso
func Granted(userID string, permission string) bool {

	usr, err := findByID(userID)
	if err != nil {
		return false
	}

	if !usr.Granted(permission) {
		return false
	}

	return true
}

//Disable deshabilita un usuario
func Disable(userID string) error {

	usr, err := findByID(userID)
	if err != nil {
		return err
	}

	usr.Enabled = false

	if _, err = update(usr); err != nil {
		return err
	}

	return nil
}

//Enable habilita un usuario
func Enable(userID string) error {

	usr, err := findByID(userID)
	if err != nil {
		return err
	}

	usr.Enabled = true

	if _, err = update(usr); err != nil {
		return err
	}

	return nil
}
