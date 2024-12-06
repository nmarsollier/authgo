package user

import (
	"time"

	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// User data structure
type User struct {
	ID          string    `dynamodbav:"id"`
	Name        string    `dynamodbav:"name" validate:"required,min=1,max=100"`
	Login       string    `dynamodbav:"login" validate:"required,min=5,max=100"`
	Password    string    `dynamodbav:"password" validate:"required"`
	Permissions []string  `dynamodbav:"permissions"`
	Enabled     bool      `dynamodbav:"enabled"`
	Created     time.Time `dynamodbav:"created"`
	Updated     time.Time `dynamodbav:"updated"`
}

// NewUser Nueva instancia de usuario
func NewUser() *User {
	return &User{
		ID:          uuid.NewV4().String(),
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
