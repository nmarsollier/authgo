package user

import (
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/app_errors"
)

// SignUpRequest es un nuevo usuario
type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Login    string `json:"login" binding:"required"`
}

// SignUp is the controller to signup new users
func SignUp(user *SignUpRequest, options ...interface{}) (string, error) {
	newUser := NewUser()
	newUser.Login = user.Login
	newUser.Name = user.Name
	newUser.SetPasswordText(user.Password)

	newUser, err := insert(newUser, options...)
	if err != nil {
		return "", err
	}

	newToken, err := token.Create(newUser.ID, options...)
	if err != nil {
		return "", app_errors.Internal
	}

	return token.Encode(newToken)
}

type SignInRequest struct {
	Password string `json:"password" binding:"required"`
	Login    string `json:"login" binding:"required"`
}

// SignIn is the controller to sign in users
func SignIn(data SignInRequest, options ...interface{}) (string, error) {
	user, err := findByLogin(data.Login, options...)
	if err != nil {
		return "", err
	}

	if !user.Enabled {
		return "", app_errors.Unauthorized
	}

	if err = user.ValidatePassword(data.Password); err != nil {
		return "", err
	}

	newToken, err := token.Create(user.ID, options...)
	if err != nil {
		return "", app_errors.Unauthorized
	}

	return token.Encode(newToken)
}

// Get wrapper para obtener un usuario
func Get(userID string, options ...interface{}) (*User, error) {
	user, err := findByID(userID, options...)
	if err != nil {
		return nil, err
	}

	if !user.Enabled {
		return nil, app_errors.NotFound
	}

	return user, err
}

// ChangePassword cambiar la contraseña del usuario indicado
func ChangePassword(userID string, current string, newPassword string, options ...interface{}) error {
	user, err := findByID(userID, options...)
	if err != nil {
		return err
	}

	if err = user.ValidatePassword(current); err != nil {
		return err
	}

	if err = user.SetPasswordText(newPassword); err != nil {
		return err
	}

	_, err = update(user, options...)

	return err
}

// Grant Le habilita los permisos enviados por parametros
func Grant(userID string, permissions []string, options ...interface{}) error {
	user, err := findByID(userID, options...)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.Grant(value)
	}
	_, err = update(user, options...)

	return err
}

// Revoke Le revoca los permisos enviados por parametros
func Revoke(userID string, permissions []string, options ...interface{}) error {
	user, err := findByID(userID, options...)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.Revoke(value)
	}
	_, err = update(user, options...)

	return err
}

// Granted verifica si el usuario tiene el permiso
func Granted(userID string, permission string, options ...interface{}) bool {
	usr, err := findByID(userID, options...)
	if err != nil {
		return false
	}

	return usr.Granted(permission)
}

// Disable deshabilita un usuario
func Disable(userID string, options ...interface{}) error {
	usr, err := findByID(userID, options...)
	if err != nil {
		return err
	}

	usr.Enabled = false

	_, err = update(usr, options...)

	return err
}

// Enable habilita un usuario
func Enable(userID string, options ...interface{}) error {
	usr, err := findByID(userID, options...)
	if err != nil {
		return err
	}

	usr.Enabled = true
	_, err = update(usr, options...)

	return err
}

// Users wrapper para obtener todos los usuarios
func Users(options ...interface{}) ([]*User, error) {
	return findAll(options...)
}
