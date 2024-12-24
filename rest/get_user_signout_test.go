package rest

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/engine/di"
	"github.com/nmarsollier/authgo/engine/errs"
	"github.com/nmarsollier/authgo/tests/engine/db"
	"github.com/nmarsollier/authgo/tests/rabbit"
	"github.com/nmarsollier/authgo/tests/router/server"
	ttoken "github.com/nmarsollier/authgo/tests/token"
	"github.com/nmarsollier/authgo/token"
	"github.com/stretchr/testify/assert"
)

func TestGetUserSignOutHappyPath(t *testing.T) {
	tokenData, tokenString := ttoken.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)

	mongo := db.NewMockMongoCollection(ctrl)
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

	r := server.TestRouter(ctrl, deps)
	InitRoutes(r)

	req, w := server.TestGetRequest("/users/signout", tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	time.Sleep(50 * time.Millisecond)
}

func TestGetUserSignOutDbUpdateError(t *testing.T) {
	tokenData, tokenString := ttoken.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)

	mongo := db.NewMockMongoCollection(ctrl)
	ttoken.ExpectTokenAuthFindOne(t, mongo, tokenData)
	db.ExpectUpdateOneError(mongo, errs.NotFound, 1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 2, 0, 1, 0, 0, 0)
	deps.SetUserCollection(mongo)
	deps.SetTokenCollection(mongo)

	r := server.TestRouter(ctrl, deps)
	InitRoutes(r)

	req, w := server.TestGetRequest("/users/signout", tokenString)
	r.ServeHTTP(w, req)

	server.AssertDocumentNotFound(t, w)
}

func TestGetUserSignOutInvalidToken(t *testing.T) {
	// Db Mocks
	ctrl := gomock.NewController(t)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 1, 1, 1, 0, 0, 0)

	r := server.TestRouter(ctrl, deps)
	InitRoutes(r)

	req, w := server.TestGetRequest("/users/signout", "123")
	r.ServeHTTP(w, req)

	server.AssertUnauthorized(t, w)
}

func TestGetUserSignOutDbFindError(t *testing.T) {
	_, tokenString := ttoken.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)
	db.ExpectFindOneError(mongo, errs.NotFound, 1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 1, 2, 1, 0, 0, 0)
	deps.SetUserCollection(mongo)
	deps.SetTokenCollection(mongo)

	r := server.TestRouter(ctrl, deps)
	InitRoutes(r)

	req, w := server.TestGetRequest("/users/signout", tokenString)
	r.ServeHTTP(w, req)

	server.AssertUnauthorized(t, w)
}
