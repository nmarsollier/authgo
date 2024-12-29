package rest

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/internal/rest"
	"github.com/nmarsollier/authgo/internal/token"
	"github.com/nmarsollier/authgo/test/mock"
	"github.com/nmarsollier/commongo/errs"
	"github.com/nmarsollier/commongo/test/mktools"
	"github.com/nmarsollier/commongo/test/mockgen"
	"github.com/stretchr/testify/assert"
)

func TestGetUserSignOutHappyPath(t *testing.T) {
	tokenData, tokenString := mock.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)

	mongo := mockgen.NewMockCollection(ctrl)
	mock.ExpectTokenAuthFindOne(t, mongo, tokenData)

	mongo.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter token.DbTokenIdFilter, update token.DbDeleteTokenDocument) (int64, error) {
			assert.Equal(t, tokenData.ID, filter.ID)

			assert.Equal(t, false, update.Set.Enabled)

			return 1, nil
		},
	).Times(1)

	rabbitMock := mktools.NewMockRabbitPublisher[string](ctrl)
	rabbitMock.EXPECT().Publish(
		gomock.Any(),
		gomock.Any()).DoAndReturn(
		func(token string, bodyStr string) error {
			assert.Contains(t, bodyStr, "Bearer")

			return nil
		},
	).Times(1)

	// REQUEST
	deps := mock.NewTestInjector(ctrl, 2, 0, 1, 0, 0, 0)
	deps.SetUserCollection(mongo)
	deps.SetTokenCollection(mongo)
	deps.SetSendLogoutPublisher(rabbitMock)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestGetRequest("/users/signout", tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	time.Sleep(50 * time.Millisecond)
}

func TestGetUserSignOutDbUpdateError(t *testing.T) {
	tokenData, tokenString := mock.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)

	mongo := mockgen.NewMockCollection(ctrl)
	mock.ExpectTokenAuthFindOne(t, mongo, tokenData)
	mktools.ExpectUpdateOneError(mongo, errs.NotFound, 1)

	// REQUEST
	deps := mock.NewTestInjector(ctrl, 2, 0, 1, 0, 0, 0)
	deps.SetUserCollection(mongo)
	deps.SetTokenCollection(mongo)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestGetRequest("/users/signout", tokenString)
	r.ServeHTTP(w, req)

	mktools.AssertDocumentNotFound(t, w)
}

func TestGetUserSignOutInvalidToken(t *testing.T) {
	// Db Mocks
	ctrl := gomock.NewController(t)

	// REQUEST
	deps := mock.NewTestInjector(ctrl, 1, 1, 1, 0, 0, 0)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestGetRequest("/users/signout", "123")
	r.ServeHTTP(w, req)

	mktools.AssertUnauthorized(t, w)
}

func TestGetUserSignOutDbFindError(t *testing.T) {
	_, tokenString := mock.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := mockgen.NewMockCollection(ctrl)
	mktools.ExpectFindOneError(mongo, errs.NotFound, 1)

	// REQUEST
	deps := mock.NewTestInjector(ctrl, 1, 2, 1, 0, 0, 0)
	deps.SetUserCollection(mongo)
	deps.SetTokenCollection(mongo)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestGetRequest("/users/signout", tokenString)
	r.ServeHTTP(w, req)

	mktools.AssertUnauthorized(t, w)
}
