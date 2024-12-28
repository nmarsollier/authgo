package rest

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/internal/engine/errs"
	"github.com/nmarsollier/authgo/internal/rest"
	"github.com/nmarsollier/authgo/internal/user"
	"github.com/nmarsollier/authgo/test/engine/di"
	"github.com/nmarsollier/authgo/test/mock"
	"github.com/nmarsollier/authgo/test/mockgen"
	"github.com/stretchr/testify/assert"
)

func TestPostUserGrantHappyPath(t *testing.T) {
	adminUserData, _ := mock.TestAdminUser()
	userData, _ := mock.TestUser()
	tokenData, tokenString := mock.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := mockgen.NewMockCollection(ctrl)

	mock.ExpectTokenAuthFindOne(t, mongodb, tokenData)

	mongodb.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserIdFilter, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, filter.ID)

			// Asign return values
			*updated = *adminUserData
			return nil
		},
	).Times(1)

	mongodb.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserIdFilter, updated *user.User) error {
			// Check parameters
			assert.Equal(t, userData.ID, filter.ID)

			// Asign return values
			*updated = *userData
			return nil
		},
	).Times(1)

	mongodb.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserIdFilter, update user.DbUserUpdateDocument) (int64, error) {
			// Check parameters
			assert.Equal(t, userData.ID, filter.ID)

			assert.Contains(t, update.Set.Permissions, "people")

			// Asign return values
			return 1, nil
		},
	).Times(1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 2, 0, 1, 0, 0, 0)
	deps.SetUserCollection(mongodb)
	deps.SetTokenCollection(mongodb)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := TestPostRequest("/users/"+userData.ID.Hex()+"/grant", grantPermissionBody{Permissions: []string{"people"}}, tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostUserGrantFindUserError_1(t *testing.T) {
	userData, _ := mock.TestUser()
	tokenData, tokenString := mock.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := mockgen.NewMockCollection(ctrl)

	mock.ExpectTokenAuthFindOne(t, mongodb, tokenData)

	mock.ExpectFindOneError(mongodb, errs.NotFound, 1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 2, 1, 1, 0, 1, 0)
	deps.SetUserCollection(mongodb)
	deps.SetTokenCollection(mongodb)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := TestPostRequest("/users/"+userData.ID.Hex()+"/grant", grantPermissionBody{Permissions: []string{"people"}}, tokenString)
	r.ServeHTTP(w, req)

	AssertUnauthorized(t, w)

}

func TestPostUserGrantFindUserError_2(t *testing.T) {
	adminUserData, _ := mock.TestAdminUser()
	userData, _ := mock.TestUser()
	tokenData, tokenString := mock.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := mockgen.NewMockCollection(ctrl)

	mock.ExpectTokenAuthFindOne(t, mongodb, tokenData)

	mongodb.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserIdFilter, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, filter.ID)

			// Asign return values
			*updated = *adminUserData
			return nil
		},
	).Times(1)

	mock.ExpectFindOneError(mongodb, errs.NotFound, 1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 2, 1, 1, 0, 0, 0)
	deps.SetUserCollection(mongodb)
	deps.SetTokenCollection(mongodb)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := TestPostRequest("/users/"+userData.ID.Hex()+"/grant", grantPermissionBody{Permissions: []string{"people"}}, tokenString)
	r.ServeHTTP(w, req)

	AssertDocumentNotFound(t, w)
}

func TestPostUserGrantNotAdmin(t *testing.T) {
	userData, _ := mock.TestUser()
	tokenData, tokenString := mock.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := mockgen.NewMockCollection(ctrl)

	mock.ExpectTokenAuthFindOne(t, mongodb, tokenData)

	mongodb.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserIdFilter, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, filter.ID)

			// Asign return values
			*updated = *userData
			return nil
		},
	).Times(1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 2, 0, 1, 0, 1, 0)
	deps.SetUserCollection(mongodb)
	deps.SetTokenCollection(mongodb)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := TestPostRequest("/users/"+userData.ID.Hex()+"/grant", grantPermissionBody{Permissions: []string{"people"}}, tokenString)
	r.ServeHTTP(w, req)

	AssertUnauthorized(t, w)
}

type grantPermissionBody struct {
	Permissions []string `json:"permissions" binding:"required"`
}
