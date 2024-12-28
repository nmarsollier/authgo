package resolvers

import (
	"context"

	"github.com/nmarsollier/authgo/internal/graph/model"
	"github.com/nmarsollier/authgo/internal/graph/tools"
	"github.com/nmarsollier/authgo/internal/usecases"
	"github.com/nmarsollier/authgo/internal/user"
)

func ChangePassword(ctx context.Context, oldPassword string, newPassword string) (bool, error) {
	token, err := tools.HeaderToken(ctx)
	if err != nil {
		return false, err
	}

	di := tools.GqlDi(ctx)
	if err := di.UserService().ChangePassword(token.UserID.Hex(), oldPassword, newPassword); err != nil {
		return false, err
	}

	return true, nil
}

func CurrentUser(ctx context.Context) (*model.User, error) {
	token, err := tools.HeaderToken(ctx)
	if err != nil {
		return nil, err
	}

	di := tools.GqlDi(ctx)
	user, err := di.UserService().FindById(token.UserID.Hex())
	if err != nil {
		return nil, err
	}

	return toUser(user), nil
}

func Disable(ctx context.Context, userID string) (bool, error) {
	if err := tools.ValidateAdmin(ctx); err != nil {
		return false, err
	}

	di := tools.GqlDi(ctx)

	if err := di.UserService().Disable(userID); err != nil {
		return false, err
	}

	return true, nil
}

func Enable(ctx context.Context, userID string) (bool, error) {
	if err := tools.ValidateAdmin(ctx); err != nil {
		return false, err
	}

	di := tools.GqlDi(ctx)

	if err := di.UserService().Enable(userID); err != nil {
		return false, err
	}

	return true, nil
}

func FindAllUsers(ctx context.Context) ([]*model.User, error) {
	if err := tools.ValidateAdmin(ctx); err != nil {
		return nil, err
	}

	env := tools.GqlDi(ctx)
	users, err := env.UserService().FindAllUsers()

	if err != nil {
		return nil, err
	}

	result := make([]*model.User, len(users))
	for i := range users {
		result[i] = toUser(users[i])
	}

	return result, nil
}

func FindUserByID(ctx context.Context, id string) (*model.User, error) {
	di := tools.GqlDi(ctx)
	user, err := di.UserService().FindById(id)
	if err != nil {
		return nil, err
	}

	return toUser(user), nil
}

func Grant(ctx context.Context, userID string, permissions []string) (bool, error) {
	if err := tools.ValidateAdmin(ctx); err != nil {
		return false, err
	}

	di := tools.GqlDi(ctx)

	if err := di.UserService().Grant(userID, permissions); err != nil {
		return false, err
	}

	return true, nil
}
func Revoke(ctx context.Context, userID string, permissions []string) (bool, error) {
	if err := tools.ValidateAdmin(ctx); err != nil {
		return false, err
	}

	di := tools.GqlDi(ctx)

	if err := di.UserService().Revoke(userID, permissions); err != nil {
		return false, err
	}

	return true, nil
}

// SignIn is the resolver for the signIn field.
func SignIn(ctx context.Context, login string, password string) (*model.Token, error) {
	env := tools.GqlDi(ctx)

	tokenString, err := env.SignInUseCase().SignIn(&usecases.SignInRequest{Login: login, Password: password})
	if err != nil {
		return nil, err
	}

	return &model.Token{
		Token: tokenString.Token,
	}, nil
}

func SignOut(ctx context.Context) (bool, error) {
	tokenString, err := tools.TokenString(ctx)
	if err != nil {
		return false, err
	}

	env := tools.GqlDi(ctx)

	if err := env.InvalidateTokenUseCase().InvalidateToken(tokenString); err != nil {
		return false, err
	}

	return true, nil
}

func SignUp(ctx context.Context, name string, login string, password string) (*model.Token, error) {
	env := tools.GqlDi(ctx)
	token, err := env.SignUpUseCase().SignUp(&usecases.SignUpRequest{
		Name:     name,
		Login:    login,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return &model.Token{
		Token: token.Token,
	}, nil
}

func toUser(user *user.UserData) (result *model.User) {
	return &model.User{
		ID:          user.Id,
		Name:        user.Name,
		Permissions: user.Permissions,
		Login:       user.Login,
		Enabled:     user.Enabled,
	}
}
