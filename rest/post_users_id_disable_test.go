package rest

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/rest/server"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/log"
	"github.com/nmarsollier/authgo/user"
	"github.com/stretchr/testify/assert"
)

func TestPostUserDisableHappyPath(t *testing.T) {
	userData, _ := user.TestAdminUser()
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenAuthFindOne(t, mongodb, tokenData)

	// User Dao Mocks
	mongodb.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserIdFilter, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, filter.ID)

			// Asign return values
			*updated = *userData
			return nil
		},
	).Times(2)

	mongodb.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserIdFilter, update user.DbUserUpdateDocument) (int64, error) {
			// Check parameters
			assert.Equal(t, userData.ID, filter.ID)

			assert.Equal(t, false, update.Set.Enabled)

			// Asign return values
			return 1, nil
		},
	).Times(1)

	// REQUEST
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 6, 0, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/users/"+tokenData.UserID.Hex()+"/disable", "", tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostUserDisableFindUserError_1(t *testing.T) {
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenAuthFindOne(t, mongodb, tokenData)

	db.ExpectFindOneError(mongodb, user.ErrID, 1)

	// REQUEST
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 6, 1, 1, 0, 1, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/users/"+tokenData.UserID.Hex()+"/disable", "", tokenString)
	r.ServeHTTP(w, req)

	server.AssertUnauthorized(t, w)
}

func TestPostUserDisableFindUserError_2(t *testing.T) {
	userData, _ := user.TestAdminUser()
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
	db.ExpectFindOneError(mongodb, user.ErrID, 1)

	// REQUEST
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 6, 1, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/users/"+tokenData.UserID.Hex()+"/disable", "", tokenString)
	r.ServeHTTP(w, req)

	server.AssertBadRequestError(t, w)
}

func TestPostUserDisableNotAdmin(t *testing.T) {
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

	req, w := server.TestPostRequest("/v1/users/"+tokenData.UserID.Hex()+"/disable", "", tokenString)
	r.ServeHTTP(w, req)

	server.AssertUnauthorized(t, w)
}
