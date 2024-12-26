package user

import (
	"github.com/nmarsollier/authgo/internal/engine/errs"
)

type UserService interface {
	ChangePassword(userID string, current string, newPassword string) error
	Disable(userID string) error
	Enable(userID string) error
	FindById(userID string) (*UserData, error)
	FindAllUsers() ([]*UserData, error)
	Grant(userID string, permissions []string) error
	Granted(userID string, permission string) bool
	Revoke(userID string, permissions []string) error
	New(login string, name string, password string) (result *UserData, err error)
	SignIn(login string, password string) (user *UserData, err error)
}

func NewUserService(userRepository UserRepository) UserService {
	return &userService{
		repository: userRepository,
	}
}

type userService struct {
	repository UserRepository
}

func (s *userService) ChangePassword(userID string, current string, newPassword string) error {
	user, err := s.repository.FindByID(userID)
	if err != nil {
		return err
	}

	if err = user.validatePassword(current); err != nil {
		return err
	}

	if err = user.setPasswordText(newPassword); err != nil {
		return err
	}

	_, err = s.repository.Update(user)

	return err
}

func (s *userService) Disable(userID string) error {
	usr, err := s.repository.FindByID(userID)
	if err != nil {
		return err
	}

	usr.Enabled = false

	_, err = s.repository.Update(usr)

	return err
}

func (s *userService) FindAllUsers() (users []*UserData, err error) {
	user, err := s.repository.FindAll()

	if err != nil {
		return
	}

	for i := 0; i < len(user); i = i + 1 {
		users = append(users, NewUserData(user[i]))
	}

	return
}

func (s *userService) New(login string, name string, password string) (*UserData, error) {
	newUser := NewUser()
	newUser.Login = login
	newUser.Name = name
	newUser.setPasswordText(password)

	result, err := s.repository.Insert(newUser)
	return NewUserData(result), err
}

func (s *userService) Get(userID string) (*UserData, error) {
	user, err := s.repository.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if !user.Enabled {
		return nil, errs.NotFound
	}

	return NewUserData(user), err
}

func (s *userService) Grant(userID string, permissions []string) error {
	user, err := s.repository.FindByID(userID)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.grant(value)
	}
	_, err = s.repository.Update(user)

	return err
}

func (s *userService) Granted(userID string, permission string) bool {
	usr, err := s.repository.FindByID(userID)
	if err != nil {
		return false
	}

	return usr.granted(permission)
}

func (s *userService) Revoke(userID string, permissions []string) error {
	user, err := s.repository.FindByID(userID)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.revoke(value)
	}
	_, err = s.repository.Update(user)

	return err
}

func (s *userService) SignIn(login string, password string) (*UserData, error) {
	user, err := s.repository.FindByLogin(login)
	if err != nil {
		return nil, err
	}

	if !user.Enabled {
		return nil, errs.Unauthorized
	}

	if err = user.validatePassword(password); err != nil {
		return nil, err
	}

	return NewUserData(user), nil
}

func (s *userService) Enable(userID string) error {
	usr, err := s.repository.FindByID(userID)
	if err != nil {
		return err
	}

	usr.Enabled = true
	_, err = s.repository.Update(usr)

	return err
}

func (s *userService) FindById(userID string) (*UserData, error) {
	user, err := s.repository.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return NewUserData(user), err
}
