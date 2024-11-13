package rest

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/rabbit"
	"github.com/nmarsollier/authgo/rest/server"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/tools/log"
	"github.com/stretchr/testify/assert"
)

func TestGetUserSignOutHappyPath(t *testing.T) {
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenAuthFindOne(t, mongo, tokenData)

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
	r := server.TestRouter(mongo, rabbitMock, log.NewTestLogger(ctrl, 10, 0, 2, 1, 0, 0))
	InitRoutes()

	req, w := server.TestGetRequest("/v1/user/signout", tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	time.Sleep(50 * time.Millisecond)
}

func TestGetUserSignOutDbUpdateError(t *testing.T) {
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)
	token.ExpectTokenAuthFindOne(t, mongo, tokenData)
	db.ExpectUpdateOneError(mongo, errs.NotFound, 1)

	// REQUEST
	r := server.TestRouter(mongo, log.NewTestLogger(ctrl, 6, 0, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestGetRequest("/v1/user/signout", tokenString)
	r.ServeHTTP(w, req)

	server.AssertDocumentNotFound(t, w)
}

func TestGetUserSignOutInvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)

	// REQUEST
	r := server.TestRouter(log.NewTestLogger(ctrl, 5, 1, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestGetRequest("/v1/user/signout", "123")
	r.ServeHTTP(w, req)

	server.AssertUnauthorized(t, w)
}

func TestGetUserSignOutDbFindError(t *testing.T) {
	_, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)

	db.ExpectFindOneError(mongo, errs.NotFound, 1)

	// REQUEST
	r := server.TestRouter(mongo, log.NewTestLogger(ctrl, 5, 2, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestGetRequest("/v1/user/signout", tokenString)
	r.ServeHTTP(w, req)

	server.AssertUnauthorized(t, w)
}
