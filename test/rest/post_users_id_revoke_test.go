package rest

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/internal/rest"
	"github.com/nmarsollier/authgo/internal/user"
	"github.com/nmarsollier/authgo/test/engine/di"
	"github.com/nmarsollier/authgo/test/mock"
	"github.com/nmarsollier/commongo/errs"
	"github.com/nmarsollier/commongo/test/mktools"
	"github.com/nmarsollier/commongo/test/mockgen"
	"github.com/stretchr/testify/assert"
)

func TestPostUserRevokeHappyPath(t *testing.T) {
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

			assert.NotContains(t, update.Set.Permissions, "user")

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

	req, w := mktools.TestPostRequest("/users/"+userData.ID.Hex()+"/revoke", revokePermissionBody{Permissions: []string{"user"}}, tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostUserRevokeFindUserError_1(t *testing.T) {
	userData, _ := mock.TestUser()
	tokenData, tokenString := mock.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := mockgen.NewMockCollection(ctrl)

	mock.ExpectTokenAuthFindOne(t, mongodb, tokenData)

	mktools.ExpectFindOneError(mongodb, errs.NotFound, 1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 2, 1, 1, 0, 1, 0)
	deps.SetUserCollection(mongodb)
	deps.SetTokenCollection(mongodb)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestPostRequest("/users/"+userData.ID.Hex()+"/revoke", revokePermissionBody{Permissions: []string{"people"}}, tokenString)
	r.ServeHTTP(w, req)

	mktools.AssertUnauthorized(t, w)

}

func TestPostUserRevokeFindUserError_2(t *testing.T) {
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

	mktools.ExpectFindOneError(mongodb, errs.NotFound, 1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 2, 1, 1, 0, 0, 0)
	deps.SetUserCollection(mongodb)
	deps.SetTokenCollection(mongodb)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestPostRequest("/users/"+userData.ID.Hex()+"/revoke", revokePermissionBody{Permissions: []string{"people"}}, tokenString)
	r.ServeHTTP(w, req)

	mktools.AssertDocumentNotFound(t, w)
}

func TestPostUserRevokeNotAdmin(t *testing.T) {
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

	req, w := mktools.TestPostRequest("/users/"+userData.ID.Hex()+"/revoke", revokePermissionBody{Permissions: []string{"people"}}, tokenString)
	r.ServeHTTP(w, req)

	mktools.AssertUnauthorized(t, w)
}

func TestGetUserSignOutMissingTokenHeader(t *testing.T) {

	// REQUEST
	ctrl := gomock.NewController(t)
	deps := di.NewTestInjector(ctrl, 1, 1, 1, 0, 0, 0)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestPostRequest("/users/123/revoke", revokePermissionBody{Permissions: []string{"people"}}, "")
	r.ServeHTTP(w, req)

	mktools.AssertUnauthorized(t, w)
}

type revokePermissionBody struct {
	Permissions []string `json:"permissions" binding:"required"`
}
