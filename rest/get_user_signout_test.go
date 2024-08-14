package rest

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/tools/apperr"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/tests"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetUserSignOutHappyPath(t *testing.T) {
	tokenData, tokenString := tests.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)

	tests.ExpectFindOneForToken(t, mongo, tokenData)

	mongo.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter primitive.M, update primitive.M) (int64, error) {
			assert.Equal(t, tokenData.ID, filter["_id"].(primitive.ObjectID))

			assert.Equal(t, false, update["$set"].(primitive.M)["enabled"])

			return 1, nil
		},
	).Times(1)

	rabbitMock := tests.MockRabbitChannel(ctrl, 1)
	rabbitMock.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(exchange string, routingKey string, body []byte) error {
			assert.Equal(t, "auth", exchange)
			assert.Equal(t, "", routingKey)
			bodyStr := string(body)
			assert.Contains(t, bodyStr, "logout")
			assert.Contains(t, bodyStr, "bearer")

			return nil
		},
	).Times(1)

	// REQUEST
	r := engine.TestRouter(mongo, rabbitMock)
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/user/signout", tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	time.Sleep(50 * time.Millisecond)
}

func TestGetUserSignOutDbUpdateError(t *testing.T) {
	tokenData, tokenString := tests.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)
	tests.ExpectFindOneForToken(t, mongo, tokenData)
	tests.ExpectUpdateOneError(mongo, apperr.NotFound, 1)

	// REQUEST
	r := engine.TestRouter(mongo)
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/user/signout", tokenString)
	r.ServeHTTP(w, req)

	tests.AssertDocumentNotFound(t, w)
}

func TestGetUserSignOutInvalidToken(t *testing.T) {
	// REQUEST
	r := engine.TestRouter()
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/user/signout", "123")
	r.ServeHTTP(w, req)

	tests.AssertUnauthorized(t, w)
}

func TestGetUserSignOutDbFindError(t *testing.T) {
	_, tokenString := tests.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)

	tests.ExpectFindOneError(mongo, apperr.NotFound, 1)

	// REQUEST
	r := engine.TestRouter(mongo)
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/user/signout", tokenString)
	r.ServeHTTP(w, req)

	tests.AssertUnauthorized(t, w)
}
