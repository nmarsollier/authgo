package rest

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/tests"
	"github.com/nmarsollier/authgo/user"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPostUserEnableHappyPath(t *testing.T) {
	userData, _ := tests.TestAdminUser()
	tokenData, tokenString := tests.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	tests.ExpectFindOneForToken(t, mongodb, tokenData)

	mongodb.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, params primitive.M, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, params["_id"])

			// Asign return values
			*updated = *userData
			return nil
		},
	).Times(2)

	mongodb.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter primitive.M, update primitive.M) (int64, error) {
			// Check parameters
			assert.Equal(t, userData.ID, filter["_id"])

			userP := update["$set"].(primitive.M)
			assert.Equal(t, true, userP["enabled"])

			// Asign return values
			return 1, nil
		},
	).Times(1)

	// REQUEST
	r := engine.TestRouter(mongodb)
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/users/"+tokenData.UserID.Hex()+"/enable", "", tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostUserEnableFindUserError_1(t *testing.T) {
	tokenData, tokenString := tests.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	tests.ExpectFindOneForToken(t, mongodb, tokenData)

	tests.ExpectFindOneError(mongodb, user.ErrID, 1)

	// REQUEST
	r := engine.TestRouter(mongodb)
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/users/"+tokenData.UserID.Hex()+"/enable", "", tokenString)
	r.ServeHTTP(w, req)

	tests.AssertUnauthorized(t, w)
}

func TestPostUserEnableFindUserError_2(t *testing.T) {
	userData, _ := tests.TestAdminUser()
	tokenData, tokenString := tests.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	tests.ExpectFindOneForToken(t, mongodb, tokenData)

	mongodb.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, params primitive.M, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, params["_id"])

			// Asign return values
			*updated = *userData
			return nil
		},
	).Times(1)
	tests.ExpectFindOneError(mongodb, user.ErrID, 1)

	// REQUEST
	r := engine.TestRouter(mongodb)
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/users/"+tokenData.UserID.Hex()+"/enable", "", tokenString)
	r.ServeHTTP(w, req)

	tests.AssertBadRequestError(t, w)
}

func TestPostUserEnableNotAdmin(t *testing.T) {
	userData, _ := tests.TestUser()
	tokenData, tokenString := tests.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	tests.ExpectFindOneForToken(t, mongodb, tokenData)

	mongodb.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, params primitive.M, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, params["_id"])

			// Asign return values
			*updated = *userData
			return nil
		},
	).Times(1)

	// REQUEST
	r := engine.TestRouter(mongodb)
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/users/"+tokenData.UserID.Hex()+"/enable", "", tokenString)
	r.ServeHTTP(w, req)

	tests.AssertUnauthorized(t, w)
}
