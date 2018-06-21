package user

import (
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errors"
)

type serviceImpl struct {
	userDao      dao
	tokenService token.Service
}

// Service es la interfaz ue define el servicio
type Service interface {
	SignUp(user *SignUpRequest) (string, error)
	SignIn(login string, password string) (string, error)
	GetUser(userID string) (*User, error)
	ChangePassword(userID string, current string, newPassword string) error
	Grant(userID string, permissions []string) error
	Revoke(userID string, permissions []string) error
	Granted(userID string, permission string) bool
	Disable(userID string) error
	Enable(userID string) error
	Users() ([]*User, error)
}

// NewService retorna una nueva instancia del servicio
func NewService() Service {
	return serviceImpl{
		userDao:      newDao(),
		tokenService: token.NewService(),
	}
}

// NewTestingService retorna un servicio con fines de test
func NewTestingService(fakeDao dao, fakedTokenService token.Service) Service {
	return serviceImpl{
		userDao:      fakeDao,
		tokenService: fakedTokenService,
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
	newUser := newUser()
	newUser.Login = user.Login
	newUser.Name = user.Name
	newUser.setPasswordText(user.Password)

	newUser, err := s.userDao.insert(newUser)
	if err != nil {
		if db.IsUniqueKeyError(err) {
			return "", ErrLoginExist
		}
		return "", err
	}

	return s.tokenService.Create(newUser.ID)
}

// SignIn is the controller to sign in users
func (s serviceImpl) SignIn(login string, password string) (string, error) {
	user, err := s.userDao.findByLogin(login)
	if err != nil {
		return "", err
	}

	if !user.Enabled {
		return "", errors.Unauthorized
	}

	if err = user.validatePassword(password); err != nil {
		return "", err
	}

	return s.tokenService.Create(user.ID)
}

// GetUser wrapper para obtener un usuario
func (s serviceImpl) GetUser(userID string) (*User, error) {
	return s.userDao.findByID(userID)
}

// ChangePassword cambiar la contraseña del usuario indicado
func (s serviceImpl) ChangePassword(userID string, current string, newPassword string) error {
	user, err := s.userDao.findByID(userID)
	if err != nil {
		return err
	}

	if err = user.validatePassword(current); err != nil {
		return err
	}

	if err = user.setPasswordText(newPassword); err != nil {
		return err
	}

	_, err = s.userDao.update(user)

	return err
}

// Grant Le habilita los permisos enviados por parametros
func (s serviceImpl) Grant(userID string, permissions []string) error {
	user, err := s.userDao.findByID(userID)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.grant(value)
	}
	_, err = s.userDao.update(user)

	return err
}

// Revoke Le revoca los permisos enviados por parametros
func (s serviceImpl) Revoke(userID string, permissions []string) error {
	user, err := s.userDao.findByID(userID)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.revoke(value)
	}
	_, err = s.userDao.update(user)

	return err
}

//Granted verifica si el usuario tiene el permiso
func (s serviceImpl) Granted(userID string, permission string) bool {
	usr, err := s.userDao.findByID(userID)
	if err != nil {
		return false
	}

	return usr.granted(permission)
}

//Disable deshabilita un usuario
func (s serviceImpl) Disable(userID string) error {
	usr, err := s.userDao.findByID(userID)
	if err != nil {
		return err
	}

	usr.Enabled = false

	_, err = s.userDao.update(usr)

	return err
}

//Enable habilita un usuario
func (s serviceImpl) Enable(userID string) error {
	usr, err := s.userDao.findByID(userID)
	if err != nil {
		return err
	}

	usr.Enabled = true
	_, err = s.userDao.update(usr)

	return err
}

// Users wrapper para obtener todos los usuarios
func (s serviceImpl) Users() ([]*User, error) {
	return s.userDao.findAll()
}
