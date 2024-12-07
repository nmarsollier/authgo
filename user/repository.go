package user

import (
	"time"

	"github.com/nmarsollier/authgo/tools/log"
)

type DbUserUpdateDocumentBody struct {
	Name        string `validate:"required,min=1,max=100"`
	Password    string `validate:"required"`
	Permissions []string
	Enabled     bool
	Updated     time.Time
}

type DbUserUpdateDocument struct {
	Set DbUserUpdateDocumentBody `bson:"$set"`
}

type DbUserLoginFilter struct {
	Login string `bson:"login"`
}

func insert(user *User, deps ...interface{}) (*User, error) {
	if err := user.validateSchema(); err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	var conn, err = GetUserDao(deps...)
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	if err := conn.Insert(user); err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	return user, nil
}

func update(user *User, deps ...interface{}) (err error) {
	if err = user.validateSchema(); err != nil {
		log.Get(deps...).Error(err)
		return
	}

	conn, err := GetUserDao(deps...)
	if err != nil {
		log.Get(deps...).Error(err)
		return
	}

	user.Updated = time.Now()

	err = conn.Update(user)
	if err != nil {
		log.Get(deps...).Error(err)
		return
	}

	return
}

// FindAll devuelve todos los usuarios
func findAll(deps ...interface{}) (users []*User, err error) {
	conn, err := GetUserDao(deps...)
	if err != nil {
		log.Get(deps...).Error(err)
		return
	}

	users, err = conn.FindAll()
	if err != nil {
		log.Get(deps...).Error(err)
	}

	return
}

// FindByID lee un usuario desde la db
func findByID(userID string, deps ...interface{}) (user *User, err error) {
	conn, err := GetUserDao(deps...)

	if err != nil {
		log.Get(deps...).Error(err)
		return
	}

	if user, err = conn.FindById(userID); err != nil {
		log.Get(deps...).Error(err)
	}

	return
}

// FindByLogin lee un usuario desde la db
func findByLogin(login string, deps ...interface{}) (user *User, err error) {
	conn, err := GetUserDao(deps...)
	if err != nil {
		log.Get(deps...).Error(err)
		return
	}

	user, err = conn.FindByLogin(login)
	if err != nil {
		log.Get(deps...).Error(err)
		return
	}

	return
}
