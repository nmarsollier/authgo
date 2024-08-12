package rest

import (
	"encoding/json"
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

func TestPostUserPasswordHappyPath(t *testing.T) {
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

	userCollection.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter primitive.M, update primitive.M) (int64, error) {
			// Check parameters
			assert.Equal(t, userData.ID, filter["_id"])

			userP := update["$set"].(primitive.M)
			assert.Equal(t, true, userP["enabled"])
			assert.Equal(t, "Name", userP["name"])
			assert.NotEmpty(t, "Password")
			assert.Equal(t, []string{"user", "other"}, userP["permissions"])

			// Asign return values
			return 1, nil
		},
	).Times(1)

	// REQUEST
	r := engine.TestRouter(token.NewTokenOption(tokenCollection), user.NewOptions(userCollection))
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user/password", changePasswordBody{Current: "123", New: "456"}, tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostUserPasswordMissingCurrent(t *testing.T) {
	tokenData, tokenString := tests.TestToken()

	// Token Dao Mocks
	ctrl := gomock.NewController(t)
	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectFindOneForToken(t, tokenCollection, tokenData)

	// REQUEST
	r := engine.TestRouter(token.NewTokenOption(tokenCollection))
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user/password", changePasswordBody{New: "456"}, tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Contains(t, result["error"], "Current")
	assert.Contains(t, result["error"], "required")
}

func TestPostUserPasswordMissingNew(t *testing.T) {
	tokenData, tokenString := tests.TestToken()

	// Token Dao Mocks
	ctrl := gomock.NewController(t)
	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectFindOneForToken(t, tokenCollection, tokenData)

	// REQUEST
	r := engine.TestRouter(token.NewTokenOption(tokenCollection))
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user/password", changePasswordBody{Current: "123"}, tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Contains(t, result["error"], "New")
	assert.Contains(t, result["error"], "required")
}

func TestPostUserPasswordWrongCurrent(t *testing.T) {
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
	r := engine.TestRouter(token.NewTokenOption(tokenCollection), user.NewOptions(userCollection))
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user/password", changePasswordBody{Current: "456", New: "456"}, tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	assert.NotEmpty(t, result["error"])
	assert.Contains(t, result["error"], "password")
	assert.Contains(t, result["error"], "invalid")
}

func TestPostUserPasswordUserNotFound(t *testing.T) {
	tokenData, tokenString := tests.TestToken()

	// Token Dao Mocks
	ctrl := gomock.NewController(t)
	tokenCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectFindOneForToken(t, tokenCollection, tokenData)

	// User Dao Mocks
	userCollection := db.NewMockMongoCollection(ctrl)
	tests.ExpectFindOneError(userCollection, app_errors.NotFound, 1)

	// REQUEST
	r := engine.TestRouter(token.NewTokenOption(tokenCollection), user.NewOptions(userCollection))
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user/password", changePasswordBody{Current: "123", New: "456"}, tokenString)
	r.ServeHTTP(w, req)

	tests.AssertDocumentNotFound(t, w)
}

func TestPostUserPasswordUpdateFails(t *testing.T) {
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

	tests.ExpectUpdateOneError(userCollection, app_errors.ErrID, 1)

	// REQUEST
	r := engine.TestRouter(token.NewTokenOption(tokenCollection), user.NewOptions(userCollection))
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user/password", changePasswordBody{Current: "123", New: "456"}, tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	assert.NotEmpty(t, result["error"])
	assert.Contains(t, result["error"], "id")
	assert.Contains(t, result["error"], "Invalid")
}
