package rest

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/internal/engine/errs"
	"github.com/nmarsollier/authgo/internal/rest"
	"github.com/nmarsollier/authgo/internal/token"
	"github.com/nmarsollier/authgo/test/engine/db"
	"github.com/nmarsollier/authgo/test/engine/di"
	"github.com/nmarsollier/authgo/test/mock"
	"github.com/nmarsollier/authgo/test/rabbit"
	ttoken "github.com/nmarsollier/authgo/test/token"
	"github.com/stretchr/testify/assert"
)

func TestGetUserSignOutHappyPath(t *testing.T) {
	tokenData, tokenString := ttoken.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)

	mongo := mock.NewMockCollection(ctrl)
	ttoken.ExpectTokenAuthFindOne(t, mongo, tokenData)

	mongo.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter token.DbTokenIdFilter, update token.DbDeleteTokenDocument) (int64, error) {
			assert.Equal(t, tokenData.ID, filter.ID)

			assert.Equal(t, false, update.Set.Enabled)

			return 1, nil
		},
	).Times(1)

	rabbitMock := rabbit.DefaultMockRabbitChannel(ctrl, 1)
	rabbitMock.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(exchange string, routingKey string, body []byte) error {
			assert.Equal(t, "auth", exchange)
			assert.Equal(t, "", routingKey)
			bodyStr := string(body)
			assert.Contains(t, bodyStr, "Bearer")

			return nil
		},
	).Times(1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 6, 0, 2, 1, 0, 0)
	deps.SetUserCollection(mongo)
	deps.SetTokenCollection(mongo)
	deps.SetRabbitChannel(rabbitMock)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := TestGetRequest("/users/signout", tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	time.Sleep(50 * time.Millisecond)
}

func TestGetUserSignOutDbUpdateError(t *testing.T) {
	tokenData, tokenString := ttoken.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)

	mongo := mock.NewMockCollection(ctrl)
	ttoken.ExpectTokenAuthFindOne(t, mongo, tokenData)
	db.ExpectUpdateOneError(mongo, errs.NotFound, 1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 2, 0, 1, 0, 0, 0)
	deps.SetUserCollection(mongo)
	deps.SetTokenCollection(mongo)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := TestGetRequest("/users/signout", tokenString)
	r.ServeHTTP(w, req)

	AssertDocumentNotFound(t, w)
}

func TestGetUserSignOutInvalidToken(t *testing.T) {
	// Db Mocks
	ctrl := gomock.NewController(t)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 1, 1, 1, 0, 0, 0)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := TestGetRequest("/users/signout", "123")
	r.ServeHTTP(w, req)

	AssertUnauthorized(t, w)
}

func TestGetUserSignOutDbFindError(t *testing.T) {
	_, tokenString := ttoken.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := mock.NewMockCollection(ctrl)
	db.ExpectFindOneError(mongo, errs.NotFound, 1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 1, 2, 1, 0, 0, 0)
	deps.SetUserCollection(mongo)
	deps.SetTokenCollection(mongo)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := TestGetRequest("/users/signout", tokenString)
	r.ServeHTTP(w, req)

	AssertUnauthorized(t, w)
}
