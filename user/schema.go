package user

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"golang.org/x/crypto/bcrypt"
)

// User data structure
type User struct {
	ID       objectid.ObjectID `bson:"_id"`
	Name     string            `bson:"name"`
	Login    string            `bson:"login"`
	Password string            `bson:"password"`
	Roles    []string          `bson:"roles"`
	Enabled  bool              `bson:"enabled"`
	Created  time.Time         `bson:"created"`
	Updated  time.Time         `bson:"updated"`
}

func newUser() *User {
	return &User{
		ID:      objectid.New(),
		Enabled: true,
		Created: time.Now(),
		Updated: time.Now(),
	}
}

// StringID obtiene el id de usuario en string
func (e *User) StringID() string {
	return e.ID.Hex()
}

func (e *User) setPasswordText(pwd string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return ErrPassword
	}

	e.Password = string(hash)
	return nil
}

func (e *User) validatePassword(plainPwd string) error {
	err := bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(plainPwd))

	if err != nil {
		return ErrPassword
	}
	return nil
}
