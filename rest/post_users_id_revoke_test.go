package rest

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/app_errors"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/tests"
	"github.com/nmarsollier/authgo/user"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPostUserRevokeHappyPath(t *testing.T) {
	adminUserData, _ := tests.TestAdminUser()
	userData, _ := tests.TestUser()
	tokenData, tokenString := tests.TestToken()

	// Token Dao Mocks
	ctrl := gomock.NewController(t)
	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectFindOneForToken(t, tokenCollection, tokenData)

	// User Dao Mocks
	userCollection := db.NewMockMongoCollection(ctrl)
	userCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, params primitive.M, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, params["_id"])

			// Asign return values
			*updated = *adminUserData
			return nil
		},
	).Times(1)

	userCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, params primitive.M, updated *user.User) error {
			// Check parameters
			assert.Equal(t, userData.ID, params["_id"])

			// Asign return values
			*updated = *userData
			return nil
		},
	).Times(1)

	userCollection.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter primitive.M, update primitive.M) (int64, error) {
			// Check parameters
			assert.Equal(t, userData.ID, filter["_id"])

			userP := update["$set"].(primitive.M)
			assert.NotContains(t, userP["permissions"], "user")

			// Asign return values
			return 1, nil
		},
	).Times(1)

	// REQUEST
	r := engine.TestRouter(token.NewProps(tokenCollection), user.NewProps(userCollection))
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/users/"+tokenData.UserID.Hex()+"/revoke", revokePermissionBody{Permissions: []string{"user"}, UserId: userData.ID.Hex()}, tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostUserRevokeFindUserError_1(t *testing.T) {
	userData, _ := tests.TestUser()
	tokenData, tokenString := tests.TestToken()

	// Token Dao Mocks
	ctrl := gomock.NewController(t)
	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectFindOneForToken(t, tokenCollection, tokenData)

	// User Dao Mocks
	userCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectFindOneError(userCollection, app_errors.NotFound, 1)

	// REQUEST
	r := engine.TestRouter(token.NewProps(tokenCollection), user.NewProps(userCollection))
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/users/"+tokenData.UserID.Hex()+"/revoke", revokePermissionBody{Permissions: []string{"people"}, UserId: userData.ID.Hex()}, tokenString)
	r.ServeHTTP(w, req)

	tests.AssertUnauthorized(t, w)

}

func TestPostUserRevokeFindUserError_2(t *testing.T) {
	adminUserData, _ := tests.TestAdminUser()
	userData, _ := tests.TestUser()
	tokenData, tokenString := tests.TestToken()

	// Token Dao Mocks
	ctrl := gomock.NewController(t)
	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectFindOneForToken(t, tokenCollection, tokenData)

	// User Dao Mocks
	userCollection := db.NewMockMongoCollection(ctrl)
	userCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, params primitive.M, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, params["_id"])

			// Asign return values
			*updated = *adminUserData
			return nil
		},
	).Times(1)

	tests.ExpectFindOneError(userCollection, app_errors.NotFound, 1)

	// REQUEST
	r := engine.TestRouter(token.NewProps(tokenCollection), user.NewProps(userCollection))
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/users/"+tokenData.UserID.Hex()+"/revoke", revokePermissionBody{Permissions: []string{"people"}, UserId: userData.ID.Hex()}, tokenString)
	r.ServeHTTP(w, req)

	tests.AssertDocumentNotFound(t, w)
}

func TestPostUserRevokeNotAdmin(t *testing.T) {
	userData, _ := tests.TestUser()
	tokenData, tokenString := tests.TestToken()

	// Token Dao Mocks
	ctrl := gomock.NewController(t)
	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectFindOneForToken(t, tokenCollection, tokenData)

	// User Dao Mocks
	userCollection := db.NewMockMongoCollection(ctrl)
	userCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, params primitive.M, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, params["_id"])

			// Asign return values
			*updated = *userData
			return nil
		},
	).Times(1)

	// REQUEST
	r := engine.TestRouter(token.NewProps(tokenCollection), user.NewProps(userCollection))
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/users/"+tokenData.UserID.Hex()+"/revoke", revokePermissionBody{Permissions: []string{"people"}, UserId: userData.ID.Hex()}, tokenString)
	r.ServeHTTP(w, req)

	tests.AssertUnauthorized(t, w)
}

func TestGetUserSignOutMissingTokenHeader(t *testing.T) {

	// REQUEST
	r := engine.TestRouter()
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/users/123/revoke", revokePermissionBody{Permissions: []string{"people"}, UserId: "123"}, "")
	r.ServeHTTP(w, req)

	tests.AssertUnauthorized(t, w)
}
