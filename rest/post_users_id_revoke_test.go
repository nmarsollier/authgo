package rest

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/rest/server"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/user"
	"github.com/stretchr/testify/assert"
)

func TestPostUserRevokeHappyPath(t *testing.T) {
	adminUserData, _ := user.TestAdminUser()
	userData, _ := user.TestUser()
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenAuthFindOne(t, mongodb, tokenData)

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
	r := server.TestRouter(mongodb)
	InitRoutes()

	req, w := server.TestPostRequest("/v1/users/"+userData.ID.Hex()+"/revoke", revokePermissionBody{Permissions: []string{"user"}}, tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostUserRevokeFindUserError_1(t *testing.T) {
	userData, _ := user.TestUser()
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenAuthFindOne(t, mongodb, tokenData)

	db.ExpectFindOneError(mongodb, errs.NotFound, 1)

	// REQUEST
	r := server.TestRouter(mongodb)
	InitRoutes()

	req, w := server.TestPostRequest("/v1/users/"+userData.ID.Hex()+"/revoke", revokePermissionBody{Permissions: []string{"people"}}, tokenString)
	r.ServeHTTP(w, req)

	server.AssertUnauthorized(t, w)

}

func TestPostUserRevokeFindUserError_2(t *testing.T) {
	adminUserData, _ := user.TestAdminUser()
	userData, _ := user.TestUser()
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenAuthFindOne(t, mongodb, tokenData)

	mongodb.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserIdFilter, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, filter.ID)

			// Asign return values
			*updated = *adminUserData
			return nil
		},
	).Times(1)

	db.ExpectFindOneError(mongodb, errs.NotFound, 1)

	// REQUEST
	r := server.TestRouter(mongodb)
	InitRoutes()

	req, w := server.TestPostRequest("/v1/users/"+userData.ID.Hex()+"/revoke", revokePermissionBody{Permissions: []string{"people"}}, tokenString)
	r.ServeHTTP(w, req)

	server.AssertDocumentNotFound(t, w)
}

func TestPostUserRevokeNotAdmin(t *testing.T) {
	userData, _ := user.TestUser()
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenAuthFindOne(t, mongodb, tokenData)

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
	r := server.TestRouter(mongodb)
	InitRoutes()

	req, w := server.TestPostRequest("/v1/users/"+userData.ID.Hex()+"/revoke", revokePermissionBody{Permissions: []string{"people"}}, tokenString)
	r.ServeHTTP(w, req)

	server.AssertUnauthorized(t, w)
}

func TestGetUserSignOutMissingTokenHeader(t *testing.T) {

	// REQUEST
	r := server.TestRouter()
	InitRoutes()

	req, w := server.TestPostRequest("/v1/users/123/revoke", revokePermissionBody{Permissions: []string{"people"}}, "")
	r.ServeHTTP(w, req)

	server.AssertUnauthorized(t, w)
}
