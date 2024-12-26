package di

import (
	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/internal/engine/db"
	"github.com/nmarsollier/authgo/internal/engine/di"
	"github.com/nmarsollier/authgo/internal/engine/log"
	"github.com/nmarsollier/authgo/internal/engine/rbt"
	"github.com/nmarsollier/authgo/internal/rabbit"
	"github.com/nmarsollier/authgo/internal/token"
	"github.com/nmarsollier/authgo/internal/usecases"
	"github.com/nmarsollier/authgo/internal/user"
	tlog "github.com/nmarsollier/authgo/test/engine/log"
	"github.com/nmarsollier/authgo/test/mock"
	trabbit "github.com/nmarsollier/authgo/test/rabbit"
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
			CurrLog: tlog.NewTestLogger(ctrl, withFieldCount, errorCount, infoCount, dataCount, warnCount, fatalCount),
		},
	}

	mongo := mock.NewMockCollection(ctrl)
	result.SetTokenCollection(mongo)
	result.SetUserCollection(mongo)
	result.SetRabbitChannel(trabbit.DefaultMockRabbitChannel(ctrl, 0))
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

func (t *TestInjector) SetSendLogoutService(service rabbit.SendLogoutService) {
	t.CurrSendLogoutService = service
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

func (t *TestInjector) SetRabbitChannel(channel rbt.RabbitChannel) {
	t.CurrRabbitChannel = channel
}

func (t *TestInjector) SetTokenCache(cache token.TokenCache) {
	t.CurrTokenCache = cache
}

func (t *TestInjector) SetTokenCollection(collection db.Collection) {
	t.CurrTokenCollection = collection
}

func (t *TestInjector) SetUserCollection(collection db.Collection) {
	t.CurrUserCollection = collection
}
