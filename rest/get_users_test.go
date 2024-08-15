package rest

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/tests"
	"github.com/nmarsollier/authgo/user"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestGetUsersHappyPath(t *testing.T) {
	userData, _ := tests.TestAdminUser()
	tokenData, tokenString := tests.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := db.NewMockMongoCollection(ctrl)

	tests.ExpectFindOneForToken(t, mongo, tokenData)

	mongo.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter user.DbUserIdFilter, updated *user.User) error {
			// Check parameters
			assert.Equal(t, tokenData.UserID, filter.ID)

			// Asign return values
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
					testUser, _ := tests.TestUser()

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
	r := engine.TestRouter(mongo)
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/users", tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	result := w.Body.String()
	assert.NotEmpty(t, result)
}

func TestGetUsersFindError(t *testing.T) {
	userData, _ := tests.TestAdminUser()
	tokenData, tokenString := tests.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongodb := db.NewMockMongoCollection(ctrl)

	tests.ExpectFindOneForToken(t, mongodb, tokenData)

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
	r := engine.TestRouter(mongodb)
	InitRoutes()

	req, w := tests.TestGetRequest("/v1/users", tokenString)
	r.ServeHTTP(w, req)

	tests.AssertDocumentNotFound(t, w)
}
