package rest

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/internal/rest"
	"github.com/nmarsollier/authgo/internal/user"
	"github.com/nmarsollier/authgo/test/mock"
	"github.com/nmarsollier/commongo/test/mktools"
	"github.com/nmarsollier/commongo/test/mockgen"
	"github.com/stretchr/testify/assert"
)

func TestPostUserDisableHappyPath(t *testing.T) {
	userData, _ := mock.TestAdminUser()
	tokenData, tokenString := mock.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := mockgen.NewMockCollection(ctrl)

	mock.ExpectTokenAuthFindOne(t, mongodb, tokenData)

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
	deps := mock.NewTestInjector(ctrl, 2, 0, 1, 0, 0, 0)
	deps.SetUserCollection(mongodb)
	deps.SetTokenCollection(mongodb)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestPostRequest("/users/"+tokenData.UserID.Hex()+"/disable", "", tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostUserDisableFindUserError_1(t *testing.T) {
	tokenData, tokenString := mock.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := mockgen.NewMockCollection(ctrl)

	mock.ExpectTokenAuthFindOne(t, mongodb, tokenData)

	mktools.ExpectFindOneError(mongodb, user.ErrID, 1)

	// REQUEST
	deps := mock.NewTestInjector(ctrl, 2, 1, 1, 0, 1, 0)
	deps.SetUserCollection(mongodb)
	deps.SetTokenCollection(mongodb)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestPostRequest("/users/"+tokenData.UserID.Hex()+"/disable", "", tokenString)
	r.ServeHTTP(w, req)

	mktools.AssertUnauthorized(t, w)
}

func TestPostUserDisableFindUserError_2(t *testing.T) {
	userData, _ := mock.TestAdminUser()
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
	mktools.ExpectFindOneError(mongodb, user.ErrID, 1)

	// REQUEST
	deps := mock.NewTestInjector(ctrl, 2, 1, 1, 0, 0, 0)
	deps.SetUserCollection(mongodb)
	deps.SetTokenCollection(mongodb)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestPostRequest("/users/"+tokenData.UserID.Hex()+"/disable", "", tokenString)
	r.ServeHTTP(w, req)

	mktools.AssertBadRequestError(t, w)
}

func TestPostUserDisableNotAdmin(t *testing.T) {
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
	deps := mock.NewTestInjector(ctrl, 2, 0, 1, 0, 1, 0)
	deps.SetUserCollection(mongodb)
	deps.SetTokenCollection(mongodb)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestPostRequest("/users/"+tokenData.UserID.Hex()+"/disable", "", tokenString)
	r.ServeHTTP(w, req)

	mktools.AssertUnauthorized(t, w)
}
