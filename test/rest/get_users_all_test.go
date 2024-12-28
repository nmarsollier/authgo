package rest

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/internal/rest"
	"github.com/nmarsollier/authgo/internal/user"
	"github.com/nmarsollier/authgo/test/engine/di"
	"github.com/nmarsollier/authgo/test/mock"
	"github.com/nmarsollier/commongo/db"
	"github.com/nmarsollier/commongo/test/mktools"
	"github.com/nmarsollier/commongo/test/mockgen"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestGetUsersHappyPath(t *testing.T) {
	userData, _ := mock.TestAdminUser()
	tokenData, tokenString := mock.TestToken()

	// Db Mocks
	ctrl := gomock.NewController(t)
	mongo := mockgen.NewMockCollection(ctrl)
	mock.ExpectTokenAuthFindOne(t, mongo, tokenData)

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
			data := mockgen.NewMockCursor(ctrl)
			data.EXPECT().Next(gomock.Any()).Return(true).Times(2)
			data.EXPECT().Next(gomock.Any()).Return(false).Times(1)

			data.EXPECT().Decode(gomock.Any()).DoAndReturn(
				func(updated *user.User) error {
					testUser, _ := mock.TestUser()

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
	deps := di.NewTestInjector(ctrl, 2, 0, 1, 0, 0, 0)
	deps.SetUserCollection(mongo)
	deps.SetTokenCollection(mongo)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestGetRequest("/users/all", tokenString)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	result := w.Body.String()
	assert.NotEmpty(t, result)
}

func TestGetUsersFindError(t *testing.T) {
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

	mongodb.EXPECT().Find(gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, filter interface{}) (db.Cursor, error) {
			// Asign return values
			return nil, mongo.ErrNoDocuments
		},
	).Times(1)

	// REQUEST
	deps := di.NewTestInjector(ctrl, 2, 1, 1, 0, 0, 0)
	deps.SetUserCollection(mongodb)
	deps.SetTokenCollection(mongodb)

	r := TestRouter(ctrl, deps)
	rest.InitRoutes(r)

	req, w := mktools.TestGetRequest("/users/all", tokenString)
	r.ServeHTTP(w, req)

	mktools.AssertDocumentNotFound(t, w)
}
