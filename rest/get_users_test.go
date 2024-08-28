package rest

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/log"
	"github.com/nmarsollier/authgo/rest/server"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/user"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestGetUsersHappyPath(t *testing.T) {
	userData, _ := user.TestAdminUser()
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenAuthFindOne(t, mongo, tokenData)

	mongo.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserIdFilter, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, filter.ID)

			// Assign return values
			*updated = *userData
			return nil
		},
	).Times(1)

	mongo.EXPECT().Find(gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter interface{}) (db.Cursor, error) {
			data := db.NewMockCursor(ctrl)
			data.EXPECT().Next(gomock.Any()).Return(true).Times(2)
			data.EXPECT().Next(gomock.Any()).Return(false).Times(1)

			data.EXPECT().Decode(gomock.Any()).DoAndReturn(
				func(updated *user.User) error {
					testUser, _ := user.TestUser()

					*updated = *testUser

					return nil
				},
			).Times(2)

			data.EXPECT().Close(gomock.Any()).Return(nil).Times(1)

			// Asign return values
			return data, nil
		},
	).Times(1)

	// REQUEST
	r := server.TestRouter(mongo, log.NewTestLogger(ctrl, 6, 0, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestGetRequest("/v1/users", tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	result := w.Body.String()
	assert.NotEmpty(t, result)
}

func TestGetUsersFindError(t *testing.T) {
	userData, _ := user.TestAdminUser()
	tokenData, tokenString := token.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	token.ExpectTokenAuthFindOne(t, mongodb, tokenData)

	mongodb.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserIdFilter, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, filter.ID)

			// Asign return values
			*updated = *userData
			return nil
		},
	).Times(1)

	mongodb.EXPECT().Find(gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter interface{}) (db.Cursor, error) {
			// Asign return values
			return nil, mongo.ErrNoDocuments
		},
	).Times(1)

	// REQUEST
	r := server.TestRouter(mongodb, log.NewTestLogger(ctrl, 6, 1, 1, 0, 0, 0))
	InitRoutes()

	req, w := server.TestGetRequest("/v1/users", tokenString)
	r.ServeHTTP(w, req)

	server.AssertDocumentNotFound(t, w)
}
