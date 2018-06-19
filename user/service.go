package user

import (
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/db"
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
		if db.IsUniqueKeyError(err) {
			return "", ErrLoginExist
		}
		return "", err
	}

	return token.Create(newUser.ID)
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

	return token.Create(user.ID)
}

// GetUser wrapper para obtener un usuario
func GetUser(userID string) (*User, error) {
	return findByID(userID)
}

// ChangePassword cambiar la contraseña del usuario indicado
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

	_, err = update(user)

	return err
}

// Grant Le habilita los permisos enviados por parametros
func Grant(userID string, permissions []string) error {
	user, err := findByID(userID)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.grant(value)
	}
	_, err = update(user)

	return err
}

// Revoke Le revoca los permisos enviados por parametros
func Revoke(userID string, permissions []string) error {
	user, err := findByID(userID)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.revoke(value)
	}
	_, err = update(user)

	return err
}

//Granted verifica si el usuario tiene el permiso
func Granted(userID string, permission string) bool {
	usr, err := findByID(userID)
	if err != nil {
		return false
	}

	return usr.granted(permission)
}

//Disable deshabilita un usuario
func Disable(userID string) error {
	usr, err := findByID(userID)
	if err != nil {
		return err
	}

	usr.Enabled = false

	_, err = update(usr)

	return err
}

//Enable habilita un usuario
func Enable(userID string) error {
	usr, err := findByID(userID)
	if err != nil {
		return err
	}

	usr.Enabled = true
	_, err = update(usr)

	return err
}

// Users wrapper para obtener todos los usuarios
func Users() ([]*User, error) {
	return findAll()
}
