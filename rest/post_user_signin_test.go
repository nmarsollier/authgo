package rest

import (
	"encoding/json"
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
	"go.mongodb.org/mongo-driver/mongo"
)

func TestPostSignInHappyPath(t *testing.T) {
	userData, password := user.TestUser()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	mongodb.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserLoginFilter, updated *user.User) error {
			// Check parameters
			assert.Equal(t, userData.Login, filter.Login)

			// Asign return values
			*updated = *userData
			return nil
		},
	).Times(1)

	mongodb.EXPECT().InsertOne(gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, token *token.Token) (string, error) {
			assert.Equal(t, true, token.Enabled)
			assert.Equal(t, userData.ID, token.UserID)
			return "", nil
		},
	).Times(1)

	// REQUEST
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 5, 0, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user/signin", user.SignInRequest{Login: userData.Login, Password: password}, "")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.NotEmpty(t, result["token"])

	tokenS := result["token"].(string)
	tokenID, userID, err := token.ExtractPayload(tokenS)
	assert.Equal(t, userID, userID)
	assert.NotEmpty(t, tokenID)
	assert.NoError(t, err)
}

func TestPostSignInWrongPassword(t *testing.T) {
	userData, _ := user.TestUser()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	mongodb.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserLoginFilter, updated *user.User) error {
			// Check parameters
			assert.Equal(t, userData.Login, filter.Login)

			// Asign return values
			*updated = *userData
			return nil
		},
	).Times(1)

	// REQUEST
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 5, 0, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user/signin", user.SignInRequest{Login: userData.Login, Password: "wrong"}, "")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	assert.NotEmpty(t, result["error"])
	assert.Contains(t, result["error"], "password")
	assert.Contains(t, result["error"], "invalid")
}

func TestPostSignInUserDisabled(t *testing.T) {
	userData, password := user.TestUser()
	userData.Enabled = false

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	mongodb.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserLoginFilter, updated *user.User) error {
			// Check parameters
			assert.Equal(t, userData.Login, filter.Login)

			// Asign return values
			*updated = *userData
			return nil
		},
	).Times(1)

	// REQUEST
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 5, 0, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user/signin", user.SignInRequest{Login: userData.Login, Password: password}, "")
	r.ServeHTTP(w, req)

	server.AssertUnauthorized(t, w)
}

func TestPostSignInMissingLogin(t *testing.T) {
	_, password := user.TestUser()

	// REQUEST
	ctrl := gomock.NewController(t)
	r := server.TestRouter(log.NewTestLogger(ctrl, 5, 0, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user/signin", user.SignInRequest{Password: password}, "")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	result := w.Body.String()
	assert.Contains(t, result, "login")
	assert.Contains(t, result, "required")
}

func TestPostSignInMissingPassword(t *testing.T) {
	userData, _ := user.TestUser()

	// REQUEST
	ctrl := gomock.NewController(t)
	r := server.TestRouter(log.NewTestLogger(ctrl, 5, 0, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user/signin", user.SignInRequest{Login: userData.Login}, "")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	result := w.Body.String()
	assert.Contains(t, result, "password")
	assert.Contains(t, result, "required")
}

func TestPostSignInUserDbError(t *testing.T) {
	userData, password := user.TestUser()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	db.ExpectFindOneError(mongodb, errs.Internal, 1)

	// REQUEST
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 5, 1, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user/signin", user.SignInRequest{Login: userData.Login, Password: password}, "")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, result["error"], "Internal server error")
}

func TestPostSignInUserNotFound(t *testing.T) {
	userData, password := user.TestUser()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	db.ExpectFindOneError(mongodb, mongo.ErrNoDocuments, 1)

	// REQUEST
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 5, 1, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user/signin", user.SignInRequest{Login: userData.Login, Password: password}, "")
	r.ServeHTTP(w, req)

	server.AssertBadRequestError(t, w)
}

func TestPostTokenDbError(t *testing.T) {
	userData, password := user.TestUser()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	db.ExpectInsertOneError(mongodb, user.ErrID, 1)

	mongodb.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserLoginFilter, updated *user.User) error {
			// Check parameters
			assert.Equal(t, userData.Login, filter.Login)

			// Asign return values
			*updated = *userData
			return nil
		},
	).Times(1)

	// REQUEST
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 5, 1, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestPostRequest("/v1/user/signin", user.SignInRequest{Login: userData.Login, Password: password}, "")
	r.ServeHTTP(w, req)

	server.AssertUnauthorized(t, w)
}
