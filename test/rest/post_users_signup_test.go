package rest

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/internal/rest"
	"github.com/nmarsollier/authgo/test/engine/di"
	"github.com/nmarsollier/authgo/test/mock"
	"github.com/nmarsollier/commongo/errs"
	"github.com/nmarsollier/commongo/test/mktools"
	"github.com/nmarsollier/commongo/test/mockgen"
	"github.com/stretchr/testify/assert"
)

func TestPostUserInHappyPath(t *testing.T) {
	userData, password := mock.TestUser()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := mockgen.NewMockCollection(ctrl)

	mktools.ExpectInsertOne(mongodb, 1)

	mktools.ExpectInsertOne(mongodb, 1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 1, 0, 1, 0, 0, 0)
	deps.SetUserCollection(mongodb)
	deps.SetTokenCollection(mongodb)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestPostRequest("/users/signup", SignUpRequest{Login: userData.Login, Password: password, Name: userData.Name}, "")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.NotEmpty(t, result["token"], "")
}

func TestPostUserMissingLogin(t *testing.T) {
	_, password := mock.TestUser()

	// REQUEST
	ctrl := gomock.NewController(t)
	deps := di.NewTestInjector(ctrl, 1, 0, 1, 0, 0, 0)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestPostRequest("/users/signup", SignUpRequest{Password: password}, "")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	result := w.Body.String()
	assert.Contains(t, result, "login")
	assert.Contains(t, result, "name")
	assert.Contains(t, result, "required")
}

func TestPostUserInvalidLoginMinRule(t *testing.T) {
	// REQUEST
	ctrl := gomock.NewController(t)
	deps := di.NewTestInjector(ctrl, 1, 1, 1, 0, 0, 0)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestPostRequest("/users/signup", SignUpRequest{Login: "a", Name: "b", Password: "c"}, "")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	result := w.Body.String()
	assert.Contains(t, result, "login")
	assert.Contains(t, result, "min")
}

func TestPostUserIvalidPassword(t *testing.T) {
	userData, _ := mock.TestUser()

	// REQUEST
	ctrl := gomock.NewController(t)
	deps := di.NewTestInjector(ctrl, 1, 0, 1, 0, 0, 0)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestPostRequest("/users/signup", SignUpRequest{Name: userData.Name, Login: userData.Login}, "")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	result := w.Body.String()
	assert.Contains(t, result, "password")
	assert.Contains(t, result, "required")
}

func TestPostUserDatabaseError(t *testing.T) {
	userData, password := mock.TestUser()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := mockgen.NewMockCollection(ctrl)

	mktools.ExpectInsertOneError(mongodb, mktools.TestOtherDbError, 1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 1, 1, 1, 0, 0, 0)
	deps.SetUserCollection(mongodb)
	deps.SetTokenCollection(mongodb)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestPostRequest("/users/signup", SignUpRequest{Login: userData.Login, Password: password, Name: userData.Name}, "")
	r.ServeHTTP(w, req)

	mktools.AssertInternalServerError(t, w)
}

func TestPostUserAlreayExist(t *testing.T) {
	userData, password := mock.TestUser()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := mockgen.NewMockCollection(ctrl)

	mktools.ExpectInsertOneError(mongodb, mktools.TestIsUniqueError, 1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 1, 1, 1, 0, 0, 0)
	deps.SetUserCollection(mongodb)
	deps.SetTokenCollection(mongodb)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestPostRequest("/users/signup", SignUpRequest{Login: userData.Login, Password: password, Name: userData.Name}, "")
	r.ServeHTTP(w, req)

	mktools.AssertBadRequestError(t, w)
}

func TestPostTokenDatabaseError(t *testing.T) {
	userData, password := mock.TestUser()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := mockgen.NewMockCollection(ctrl)

	mktools.ExpectInsertOne(mongodb, 1)

	mktools.ExpectInsertOneError(mongodb, errs.Internal, 1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 1, 1, 1, 0, 0, 0)
	deps.SetUserCollection(mongodb)
	deps.SetTokenCollection(mongodb)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestPostRequest("/users/signup", SignUpRequest{Login: userData.Login, Password: password, Name: userData.Name}, "")
	r.ServeHTTP(w, req)

	mktools.AssertInternalServerError(t, w)
}

type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Login    string `json:"login" binding:"required"`
}
