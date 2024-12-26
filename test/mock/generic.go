package mock

import (
	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/test/mockgen"
)

func ExpectFindOneError(coll *mockgen.MockCollection, err error, times int) {
	coll.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(arg1 interface{}, params interface{}, update interface{}) error {
			return err
		},
	).Times(times)
}
