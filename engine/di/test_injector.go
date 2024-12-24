package di

import (
	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/engine/db"
	"github.com/nmarsollier/authgo/engine/log"
	"github.com/nmarsollier/authgo/rabbit"
	tdb "github.com/nmarsollier/authgo/tests/engine/db"
	tlog "github.com/nmarsollier/authgo/tests/engine/log"
	trabbit "github.com/nmarsollier/authgo/tests/rabbit"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/usecases"
	"github.com/nmarsollier/authgo/user"
)

type TestInjector struct {
	injector
}

func NewTestInjector(
	ctrl *gomock.Controller, withFieldCount int, errorCount int, infoCount int, dataCount int, warnCount int, fatalCount int,
) *TestInjector {

	result := &TestInjector{
		injector: injector{
			log: tlog.NewTestLogger(ctrl, withFieldCount, errorCount, infoCount, dataCount, warnCount, fatalCount),
		},
	}

	mongo := tdb.NewMockMongoCollection(ctrl)
	result.SetTokenCollection(mongo)
	result.SetUserCollection(mongo)
	result.SetRabbitChannel(trabbit.DefaultMockRabbitChannel(ctrl, 0))
	return result
}

func (t *TestInjector) SetLogger(log log.LogRusEntry) {
	t.log = log
}

func (t *TestInjector) SetInvalidateTokenUseCase(useCase usecases.InvalidateTokenUseCase) {
	t.invalidateTokenUseCase = useCase
}

func (t *TestInjector) SetSendLogoutService(service rabbit.SendLogoutService) {
	t.sendLogoutService = service
}

func (t *TestInjector) SetSignInUseCase(useCase usecases.SignInUseCase) {
	t.signInUseCase = useCase
}

func (t *TestInjector) SetSignUpUseCase(useCase usecases.SignUpUseCase) {
	t.signUpUseCase = useCase
}

func (t *TestInjector) SetTokenRepository(repository token.TokenRepository) {
	t.tokenRepository = repository
}

func (t *TestInjector) SetTokenService(service token.TokenService) {
	t.tokenService = service
}

func (t *TestInjector) SetUserRepository(repository user.UserRepository) {
	t.userRepository = repository
}

func (t *TestInjector) SetUserService(service user.UserService) {
	t.userService = service
}

func (t *TestInjector) SetRabbitChannel(channel rabbit.RabbitChannel) {
	t.rabbitChannel = channel
}

func (t *TestInjector) SetTokenCache(cache token.TokenCache) {
	t.tokenCache = cache
}

func (t *TestInjector) SetTokenCollection(collection db.MongoCollection) {
	t.tokenCollection = collection
}

func (t *TestInjector) SetUserCollection(collection db.MongoCollection) {
	t.userCollection = collection
}
