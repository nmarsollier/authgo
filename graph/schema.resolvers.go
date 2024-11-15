package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.56

import (
	"context"
	"fmt"

	"github.com/nmarsollier/authgo/graph/model"
	"github.com/nmarsollier/authgo/graph/tools"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/user"
)

// SignIn is the resolver for the signIn field.
func (r *mutationResolver) SignIn(ctx context.Context, login string, password string) (*user.TokenResponse, error) {
	env := tools.GqlCtx(ctx)

	tokenString, err := user.SignIn(user.SignInRequest{Login: login, Password: password}, env...)
	if err != nil {
		return nil, err
	}

	return tokenString, nil
}

// SignUp is the resolver for the signUp field.
func (r *mutationResolver) SignUp(ctx context.Context, name string, login string, password string) (*user.TokenResponse, error) {
	panic(fmt.Errorf("not implemented: SignUp - signUp"))
}

// SignOut is the resolver for the signOut field.
func (r *mutationResolver) SignOut(ctx context.Context) (bool, error) {
	tokenString, err := tools.TokenString(ctx)
	if err != nil {
		return false, err
	}

	env := tools.GqlCtx(ctx)

	if err := token.Invalidate(tokenString, env...); err != nil {
		return false, err
	}

	return true, nil
}

// ChangePassword is the resolver for the changePassword field.
func (r *mutationResolver) ChangePassword(ctx context.Context, oldPassword string, newPassword string) (bool, error) {
	token, err := tools.HeaderToken(ctx)
	if err != nil {
		return false, err
	}

	env := tools.GqlCtx(ctx)
	if err := user.ChangePassword(token.UserID.Hex(), oldPassword, newPassword, env...); err != nil {
		return false, err
	}

	return true, nil
}

// Enable is the resolver for the enable field.
func (r *mutationResolver) Enable(ctx context.Context, userID string) (bool, error) {
	if err := tools.ValidateAdmin(ctx); err != nil {
		return false, err
	}

	env := tools.GqlCtx(ctx)

	if err := user.Enable(userID, env...); err != nil {
		return false, err
	}

	return true, nil
}

// Disable is the resolver for the disable field.
func (r *mutationResolver) Disable(ctx context.Context, userID string) (bool, error) {
	if err := tools.ValidateAdmin(ctx); err != nil {
		return false, err
	}

	env := tools.GqlCtx(ctx)

	if err := user.Disable(userID, env...); err != nil {
		return false, err
	}

	return true, nil
}

// Grant is the resolver for the grant field.
func (r *mutationResolver) Grant(ctx context.Context, userID string, permissions []string) (bool, error) {
	if err := tools.ValidateAdmin(ctx); err != nil {
		return false, err
	}

	env := tools.GqlCtx(ctx)

	if err := user.Grant(userID, permissions, env...); err != nil {
		return false, err
	}

	return true, nil
}

// Revoke is the resolver for the revoke field.
func (r *mutationResolver) Revoke(ctx context.Context, userID string, permissions []string) (bool, error) {
	if err := tools.ValidateAdmin(ctx); err != nil {
		return false, err
	}

	env := tools.GqlCtx(ctx)

	if err := user.Revoke(userID, permissions, env...); err != nil {
		return false, err
	}

	return true, nil
}

// CurrentUser is the resolver for the currentUser field.
func (r *queryResolver) CurrentUser(ctx context.Context) (*user.UserResponse, error) {
	token, err := tools.HeaderToken(ctx)
	if err != nil {
		return nil, err
	}

	env := tools.GqlCtx(ctx)
	user, err := user.Get(token.UserID.Hex(), env...)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*user.UserResponse, error) {
	if err := tools.ValidateAdmin(ctx); err != nil {
		return nil, err
	}

	env := tools.GqlCtx(ctx)
	result, err := user.FindAllUsers(env...)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Mutation returns model.MutationResolver implementation.
func (r *Resolver) Mutation() model.MutationResolver { return &mutationResolver{r} }

// Query returns model.QueryResolver implementation.
func (r *Resolver) Query() model.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
