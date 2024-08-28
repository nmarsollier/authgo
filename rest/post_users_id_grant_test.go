package rest

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/log"
	"github.com/nmarsollier/authgo/rest/server"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/user"
	"github.com/stretchr/testify/assert"
)

func TestPostUserGrantHappyPath(t *testing.T) {
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

			assert.Contains(t, update.Set.Permissions, "people")

			// Asign return values
			return 1, nil
		},
	).Times(1)

	// REQUEST
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 6, 0, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/users/"+userData.ID.Hex()+"/grant", grantPermissionBody{Permissions: []string{"people"}}, tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostUserGrantFindUserError_1(t *testing.T) {
	userData, _ := user.TestUser()
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenAuthFindOne(t, mongodb, tokenData)

	db.ExpectFindOneError(mongodb, errs.NotFound, 1)

	// REQUEST
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 6, 1, 1, 0, 1, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/users/"+userData.ID.Hex()+"/grant", grantPermissionBody{Permissions: []string{"people"}}, tokenString)
	r.ServeHTTP(w, req)

	server.AssertUnauthorized(t, w)

}

func TestPostUserGrantFindUserError_2(t *testing.T) {
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
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 6, 1, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/users/"+userData.ID.Hex()+"/grant", grantPermissionBody{Permissions: []string{"people"}}, tokenString)
	r.ServeHTTP(w, req)

	server.AssertDocumentNotFound(t, w)
}

func TestPostUserGrantNotAdmin(t *testing.T) {
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
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 6, 0, 1, 0, 1, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/users/"+userData.ID.Hex()+"/grant", grantPermissionBody{Permissions: []string{"people"}}, tokenString)
	r.ServeHTTP(w, req)

	server.AssertUnauthorized(t, w)
}
