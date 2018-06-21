package user

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"golang.org/x/crypto/bcrypt"
	validator "gopkg.in/go-playground/validator.v9"
)

// User data structure
type User struct {
	ID          objectid.ObjectID `bson:"_id"`
	Name        string            `bson:"name" validate:"required,min=1,max=100"`
	Login       string            `bson:"login" validate:"required,min=5,max=100"`
	Password    string            `bson:"password" validate:"required"`
	Permissions []string          `bson:"permissions"`
	Enabled     bool              `bson:"enabled"`
	Created     time.Time         `bson:"created"`
	Updated     time.Time         `bson:"updated"`
}

func NewUser() *User {
	return &User{
		ID:          objectid.New(),
		Enabled:     true,
		Created:     time.Now(),
		Updated:     time.Now(),
		Permissions: []string{"user"},
	}
}

func (e *User) SetPasswordText(pwd string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return ErrPassword
	}

	e.Password = string(hash)
	return nil
}

func (e *User) ValidatePassword(plainPwd string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(plainPwd)); err != nil {
		return ErrPassword
	}
	return nil
}

// Granted verifica si el usuario tiene el permiso indicado
func (e *User) Granted(permission string) bool {
	for _, p := range e.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// Grant le otorga el permiso indicado al usuario
func (e *User) Grant(permission string) {
	if !e.Granted(permission) {
		e.Permissions = append(e.Permissions, permission)
	}
}

// Revoke le revoca el permiso indicado al usuario
func (e *User) Revoke(permission string) {
	if e.Granted(permission) {
		var newPermissions []string
		for _, p := range e.Permissions {
			if p != permission {
				newPermissions = append(newPermissions, p)
			}
		}
		e.Permissions = newPermissions
	}
}

func (e *User) ValidateSchema() error {
	validate := validator.New()
	return validate.Struct(e)
}
