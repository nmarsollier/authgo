package mock

import (
	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/internal/di"
	"github.com/nmarsollier/authgo/internal/token"
	"github.com/nmarsollier/authgo/internal/usecases"
	"github.com/nmarsollier/authgo/internal/user"
	"github.com/nmarsollier/commongo/cache"
	"github.com/nmarsollier/commongo/db"
	"github.com/nmarsollier/commongo/log"
	"github.com/nmarsollier/commongo/rbt"
	"github.com/nmarsollier/commongo/test/mktools"
	"github.com/nmarsollier/commongo/test/mockgen"
	"go.mongodb.org/mongo-driver/mongo"
)

type TestInjector struct {
	di.Deps
}

func NewTestInjector(
	ctrl *gomock.Controller, withFieldCount int, errorCount int, infoCount int, dataCount int, warnCount int, fatalCount int,
) *TestInjector {

	result := &TestInjector{
		Deps: di.Deps{
			CurrLog: mktools.NewTestLogger(ctrl, withFieldCount, errorCount, infoCount, dataCount, warnCount, fatalCount),
		},
	}

	mongo := mockgen.NewMockCollection(ctrl)
	result.SetTokenCollection(mongo)
	result.SetUserCollection(mongo)
	result.SetSendLogoutPublisher(mktools.NewMockRabbitPublisher[string](ctrl))
	return result
}

func (t *TestInjector) SetDatabase(database *mongo.Database) {
	t.CurrDatabase = database
}

func (t *TestInjector) SetLogger(log log.LogRusEntry) {
	t.CurrLog = log
}

func (t *TestInjector) SetInvalidateTokenUseCase(useCase usecases.InvalidateTokenUseCase) {
	t.CurrInvalidateTokenUseCase = useCase
}

func (t *TestInjector) SetSendLogoutPublisher(service rbt.RabbitPublisher[string]) {
	t.CurrSendLogout = service
}

func (t *TestInjector) SetSignInUseCase(useCase usecases.SignInUseCase) {
	t.CurrSignInUseCase = useCase
}

func (t *TestInjector) SetSignUpUseCase(useCase usecases.SignUpUseCase) {
	t.CurrSignUpUseCase = useCase
}

func (t *TestInjector) SetTokenRepository(repository token.TokenRepository) {
	t.CurrTokenRepository = repository
}

func (t *TestInjector) SetTokenService(service token.TokenService) {
	t.CurrTokenService = service
}

func (t *TestInjector) SetUserRepository(repository user.UserRepository) {
	t.CurrUserRepository = repository
}

func (t *TestInjector) SetUserService(service user.UserService) {
	t.CurrUserService = service
}

func (t *TestInjector) SetTokenCache(cache cache.Cache[token.Token]) {
	t.CurrTokenCache = cache
}

func (t *TestInjector) SetTokenCollection(collection db.Collection) {
	t.CurrTokenCollection = collection
}

func (t *TestInjector) SetUserCollection(collection db.Collection) {
	t.CurrUserCollection = collection
}
