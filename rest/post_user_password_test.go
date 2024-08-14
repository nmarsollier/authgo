package rest

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/tools/tests"
	"github.com/nmarsollier/authgo/user"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPostUserPasswordHappyPath(t *testing.T) {
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

	mongodb.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
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
	r := engine.TestRouter(mongodb)
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user/password", changePasswordBody{Current: "123", New: "456"}, tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostUserPasswordMissingCurrent(t *testing.T) {
	tokenData, tokenString := tests.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	tests.ExpectFindOneForToken(t, mongodb, tokenData)

	// REQUEST
	r := engine.TestRouter(mongodb)
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user/password", changePasswordBody{New: "456"}, tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	result := w.Body.String()
	assert.Contains(t, result, "current")
	assert.Contains(t, result, "required")
}

func TestPostUserPasswordMissingNew(t *testing.T) {
	tokenData, tokenString := tests.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	tests.ExpectFindOneForToken(t, mongodb, tokenData)

	// REQUEST
	r := engine.TestRouter(mongodb)
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user/password", changePasswordBody{Current: "123"}, tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	result := w.Body.String()

	assert.Contains(t, result, "new")
	assert.Contains(t, result, "required")
}

func TestPostUserPasswordWrongCurrent(t *testing.T) {
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

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	tests.ExpectFindOneForToken(t, mongodb, tokenData)

	tests.ExpectFindOneError(mongodb, errs.NotFound, 1)

	// REQUEST
	r := engine.TestRouter(mongodb)
	InitRoutes()

	req, w := tests.TestPostRequest("/v1/user/password", changePasswordBody{Current: "123", New: "456"}, tokenString)
	r.ServeHTTP(w, req)

	tests.AssertDocumentNotFound(t, w)
}

func TestPostUserPasswordUpdateFails(t *testing.T) {
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

	tests.ExpectUpdateOneError(mongodb, user.ErrID, 1)

	// REQUEST
	r := engine.TestRouter(mongodb)
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
