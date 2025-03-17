package user

import (
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User data structure
type User struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name" validate:"required,min=1,max=100"`
	Login       string             `bson:"login" validate:"required,min=5,max=100"`
	Password    string             `bson:"password" validate:"required"`
	Permissions []string           `bson:"permissions"`
	Enabled     bool               `bson:"enabled"`
	Created     time.Time          `bson:"created"`
	Updated     time.Time          `bson:"updated"`
}

// newUser Nueva instancia de usuario
func newUser() *User {
	return &User{
		ID:          primitive.NewObjectID(),
		Enabled:     true,
		Created:     time.Now(),
		Updated:     time.Now(),
		Permissions: []string{"user"},
	}
}

// setPasswordText Asigna la contraseña en modo texto, la encripta
func (e *User) setPasswordText(pwd string) error {
	hash, err := encryptPassword(pwd)
	if err != nil {
		return ErrPassword
	}

	e.Password = hash
	return nil
}

// validatePassword Valida si la contraseña es la correcta
func (e *User) validatePassword(plainPwd string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(plainPwd)); err != nil {
		return ErrPassword
	}
	return nil
}

// granted verifica si el usuario tiene el permiso indicado
func (e *User) granted(permission string) bool {
	for _, p := range e.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// grant le otorga el permiso indicado al usuario
func (e *User) grant(permission string) {
	if !e.granted(permission) {
		e.Permissions = append(e.Permissions, permission)
	}
}

// revoke le revoca el permiso indicado al usuario
func (e *User) revoke(permission string) {
	if e.granted(permission) {
		var newPermissions []string
		for _, p := range e.Permissions {
			if p != permission {
				newPermissions = append(newPermissions, p)
			}
		}
		e.Permissions = newPermissions
	}
}

// validateSchema valida la estructura para ser insertada en la db
func (e *User) validateSchema() error {
	return validator.New().Struct(e)
}

func encryptPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", ErrPassword
	}

	return string(hash), nil
}
