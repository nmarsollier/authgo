package user

import (
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/app_errors"
	"github.com/nmarsollier/authgo/tools/db"
)

// SignUpRequest es un nuevo usuario
type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Login    string `json:"login" binding:"required"`
}

// SignUp is the controller to signup new users
func SignUp(user *SignUpRequest, props ...interface{}) (string, error) {
	newUser := NewUser()
	newUser.Login = user.Login
	newUser.Name = user.Name
	newUser.SetPasswordText(user.Password)

	newUser, err := insert(newUser, props...)
	if err != nil {
		if db.IsUniqueKeyError(err) {
			return "", ErrLoginExist
		}
		return "", err
	}

	newToken, err := token.Create(newUser.ID, props...)
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
func SignIn(data SignInRequest, props ...interface{}) (string, error) {
	user, err := findByLogin(data.Login, props...)
	if err != nil {
		return "", err
	}

	if !user.Enabled {
		return "", app_errors.Unauthorized
	}

	if err = user.ValidatePassword(data.Password); err != nil {
		return "", err
	}

	newToken, err := token.Create(user.ID, props...)
	if err != nil {
		return "", app_errors.Unauthorized
	}

	return token.Encode(newToken)
}

// Get wrapper para obtener un usuario
func Get(userID string, props ...interface{}) (*User, error) {
	user, err := findByID(userID, props...)
	if err != nil {
		return nil, err
	}

	if !user.Enabled {
		return nil, app_errors.NotFound
	}

	return user, err
}

// ChangePassword cambiar la contraseña del usuario indicado
func ChangePassword(userID string, current string, newPassword string, props ...interface{}) error {
	user, err := findByID(userID, props...)
	if err != nil {
		return err
	}

	if err = user.ValidatePassword(current); err != nil {
		return err
	}

	if err = user.SetPasswordText(newPassword); err != nil {
		return err
	}

	_, err = update(user, props...)

	return err
}

// Grant Le habilita los permisos enviados por parametros
func Grant(userID string, permissions []string, props ...interface{}) error {
	user, err := findByID(userID, props...)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.Grant(value)
	}
	_, err = update(user, props...)

	return err
}

// Revoke Le revoca los permisos enviados por parametros
func Revoke(userID string, permissions []string, props ...interface{}) error {
	user, err := findByID(userID, props...)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.Revoke(value)
	}
	_, err = update(user, props...)

	return err
}

// Granted verifica si el usuario tiene el permiso
func Granted(userID string, permission string, props ...interface{}) bool {
	usr, err := findByID(userID, props...)
	if err != nil {
		return false
	}

	return usr.Granted(permission)
}

// Disable deshabilita un usuario
func Disable(userID string, props ...interface{}) error {
	usr, err := findByID(userID, props...)
	if err != nil {
		return err
	}

	usr.Enabled = false

	_, err = update(usr, props...)

	return err
}

// Enable habilita un usuario
func Enable(userID string, props ...interface{}) error {
	usr, err := findByID(userID, props...)
	if err != nil {
		return err
	}

	usr.Enabled = true
	_, err = update(usr, props...)

	return err
}

// Users wrapper para obtener todos los usuarios
func Users(props ...interface{}) ([]*User, error) {
	return findAll(props...)
}
