package user

import (
	"github.com/nmarsollier/authgo/security"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errors"
)

type serviceImpl struct {
	dao        Dao
	secService security.Service
}

// Service es la interfaz ue define el servicio
type Service interface {
	SignUp(user *SignUpRequest) (string, error)
	SignIn(login string, password string) (string, error)
	Get(userID string) (*User, error)
	ChangePassword(userID string, current string, newPassword string) error
	Grant(userID string, permissions []string) error
	Revoke(userID string, permissions []string) error
	Granted(userID string, permission string) bool
	Disable(userID string) error
	Enable(userID string) error
	Users() ([]*User, error)
}

// NewService retorna una nueva instancia del servicio
func NewService() (Service, error) {
	secService, err := security.NewService()
	if err != nil {
		return nil, err
	}

	dao, err := newDao()
	if err != nil {
		return nil, err
	}

	return serviceImpl{
		dao:        dao,
		secService: secService,
	}, nil
}

// MockedService permite mockear el servicio
func MockedService(fakeDao Dao, fakeTRepo security.Service) Service {
	return serviceImpl{
		dao:        fakeDao,
		secService: fakeTRepo,
	}
}

// SignUpRequest es un nuevo usuario
type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Login    string `json:"login" binding:"required"`
}

// SignUp is the controller to signup new users
func (s serviceImpl) SignUp(user *SignUpRequest) (string, error) {
	newUser := NewUser()
	newUser.Login = user.Login
	newUser.Name = user.Name
	newUser.SetPasswordText(user.Password)

	newUser, err := s.dao.Insert(newUser)
	if err != nil {
		if db.IsUniqueKeyError(err) {
			return "", ErrLoginExist
		}
		return "", err
	}

	token, err := s.secService.Create(newUser.ID)
	if err != nil {
		return "", nil
	}

	return token.Encode()
}

// SignIn is the controller to sign in users
func (s serviceImpl) SignIn(login string, password string) (string, error) {
	user, err := s.dao.FindByLogin(login)
	if err != nil {
		return "", err
	}

	if !user.Enabled {
		return "", errors.Unauthorized
	}

	if err = user.ValidatePassword(password); err != nil {
		return "", err
	}

	token, err := s.secService.Create(user.ID)
	if err != nil {
		return "", nil
	}

	return token.Encode()
}

// Get wrapper para obtener un usuario
func (s serviceImpl) Get(userID string) (*User, error) {
	return s.dao.FindByID(userID)
}

// ChangePassword cambiar la contraseña del usuario indicado
func (s serviceImpl) ChangePassword(userID string, current string, newPassword string) error {
	user, err := s.dao.FindByID(userID)
	if err != nil {
		return err
	}

	if err = user.ValidatePassword(current); err != nil {
		return err
	}

	if err = user.SetPasswordText(newPassword); err != nil {
		return err
	}

	_, err = s.dao.Update(user)

	return err
}

// Grant Le habilita los permisos enviados por parametros
func (s serviceImpl) Grant(userID string, permissions []string) error {
	user, err := s.dao.FindByID(userID)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.Grant(value)
	}
	_, err = s.dao.Update(user)

	return err
}

// Revoke Le revoca los permisos enviados por parametros
func (s serviceImpl) Revoke(userID string, permissions []string) error {
	user, err := s.dao.FindByID(userID)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.Revoke(value)
	}
	_, err = s.dao.Update(user)

	return err
}

//Granted verifica si el usuario tiene el permiso
func (s serviceImpl) Granted(userID string, permission string) bool {
	usr, err := s.dao.FindByID(userID)
	if err != nil {
		return false
	}

	return usr.Granted(permission)
}

//Disable deshabilita un usuario
func (s serviceImpl) Disable(userID string) error {
	usr, err := s.dao.FindByID(userID)
	if err != nil {
		return err
	}

	usr.Enabled = false

	_, err = s.dao.Update(usr)

	return err
}

//Enable habilita un usuario
func (s serviceImpl) Enable(userID string) error {
	usr, err := s.dao.FindByID(userID)
	if err != nil {
		return err
	}

	usr.Enabled = true
	_, err = s.dao.Update(usr)

	return err
}

// Users wrapper para obtener todos los usuarios
func (s serviceImpl) Users() ([]*User, error) {
	return s.dao.FindAll()
}
