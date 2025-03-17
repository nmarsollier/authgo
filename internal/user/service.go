package user

import (
	"github.com/nmarsollier/authgo/internal/common/errs"
	"github.com/nmarsollier/authgo/internal/common/log"
)

func ChangePassword(
	log log.LogRusEntry,
	userID string,
	current string,
	newPassword string,
) error {
	user, err := findByID(log, userID)
	if err != nil {
		return err
	}

	if err = user.validatePassword(current); err != nil {
		return err
	}

	if err = user.setPasswordText(newPassword); err != nil {
		return err
	}

	_, err = update(log, user)

	return err
}

func Disable(
	log log.LogRusEntry,
	userID string,
) error {
	usr, err := findByID(log, userID)
	if err != nil {
		return err
	}

	usr.Enabled = false

	_, err = update(log, usr)

	return err
}

func FindAllUsers(log log.LogRusEntry) (users []*UserData, err error) {
	user, err := findAll(log)

	if err != nil {
		return
	}

	for i := 0; i < len(user); i = i + 1 {
		users = append(users, newUserData(user[i]))
	}

	return
}

func New(log log.LogRusEntry, login string, name string, password string) (*UserData, error) {
	newUser := newUser()
	newUser.Login = login
	newUser.Name = name
	newUser.setPasswordText(password)

	result, err := insert(log, newUser)
	if err != nil {
		return nil, err
	}

	return newUserData(result), nil
}

func Get(log log.LogRusEntry, userID string) (*UserData, error) {
	user, err := findByID(log, userID)
	if err != nil {
		return nil, err
	}

	if !user.Enabled {
		return nil, errs.NotFound
	}

	return newUserData(user), err
}

func Grant(log log.LogRusEntry, userID string, permissions []string) error {
	user, err := findByID(log, userID)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.grant(value)
	}
	_, err = update(log, user)

	return err
}

func Granted(log log.LogRusEntry, userID string, permission string) bool {
	usr, err := findByID(log, userID)
	if err != nil {
		return false
	}

	return usr.granted(permission)
}

func Revoke(log log.LogRusEntry, userID string, permissions []string) error {
	user, err := findByID(log, userID)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.revoke(value)
	}
	_, err = update(log, user)

	return err
}

func SignIn(log log.LogRusEntry, login string, password string) (*UserData, error) {
	user, err := findByLogin(log, login)
	if err != nil {
		return nil, err
	}

	if !user.Enabled {
		return nil, errs.Unauthorized
	}

	if err = user.validatePassword(password); err != nil {
		return nil, err
	}

	return newUserData(user), nil
}

func Enable(log log.LogRusEntry, userID string) error {
	usr, err := findByID(log, userID)
	if err != nil {
		return err
	}

	usr.Enabled = true
	_, err = update(log, usr)

	return err
}

func FindById(log log.LogRusEntry, userID string) (*UserData, error) {
	user, err := findByID(log, userID)
	if err != nil {
		return nil, err
	}

	if !user.Enabled {
		return nil, errs.NotFound
	}

	return newUserData(user), err
}
