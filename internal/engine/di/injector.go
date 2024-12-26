package di

import (
	"github.com/nmarsollier/authgo/internal/engine/db"
	"github.com/nmarsollier/authgo/internal/engine/env"
	"github.com/nmarsollier/authgo/internal/engine/log"
	"github.com/nmarsollier/authgo/internal/engine/rbt"
	"github.com/nmarsollier/authgo/internal/rabbit"
	"github.com/nmarsollier/authgo/internal/token"
	"github.com/nmarsollier/authgo/internal/usecases"
	"github.com/nmarsollier/authgo/internal/user"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

// Singletons
var tokenCache token.TokenCache
var tokenCollection db.Collection
var userCollection db.Collection
var database *mongo.Database

type Injector interface {
	Database() *mongo.Database
	InvalidateTokenUseCase() usecases.InvalidateTokenUseCase
	Logger() log.LogRusEntry
	RabbitChannel() rbt.RabbitChannel
	SendLogoutService() rabbit.SendLogoutService
	SignInUseCase() usecases.SignInUseCase
	SignUpUseCase() usecases.SignUpUseCase
	TokenCache() token.TokenCache
	TokenCollection() db.Collection
	TokenRepository() token.TokenRepository
	TokenService() token.TokenService
	UserCollection() db.Collection
	UserRepository() user.UserRepository
	UserService() user.UserService
}

type Deps struct {
	CurrInvalidateTokenUseCase usecases.InvalidateTokenUseCase
	CurrLog                    log.LogRusEntry
	CurrSendLogoutService      rabbit.SendLogoutService
	CurrSignInUseCase          usecases.SignInUseCase
	CurrSignUpUseCase          usecases.SignUpUseCase
	CurrTokenRepository        token.TokenRepository
	CurrTokenService           token.TokenService
	CurrUserRepository         user.UserRepository
	CurrUserService            user.UserService
	CurrRabbitChannel          rbt.RabbitChannel
	CurrTokenCache             token.TokenCache
	CurrTokenCollection        db.Collection
	CurrUserCollection         db.Collection
	CurrDatabase               *mongo.Database
}

func NewInjector(log log.LogRusEntry) Injector {
	return &Deps{
		CurrLog: log,
	}
}

func (i *Deps) Database() *mongo.Database {
	if i.CurrDatabase != nil {
		return i.CurrDatabase
	}

	if database != nil {
		return database
	}

	database, err := db.NewDatabase(env.Get().MongoURL, "authgo")
	if err != nil {
		i.CurrLog.Fatal(err)
		return nil
	}

	return database
}

func (i *Deps) Logger() log.LogRusEntry {
	return i.CurrLog
}

func (i *Deps) TokenService() token.TokenService {
	if i.CurrTokenService != nil {
		return i.CurrTokenService
	}

	i.CurrTokenService = token.NewTokenService(i.CurrLog, i.TokenCache(), i.TokenRepository())
	return i.CurrTokenService
}

func (i *Deps) TokenRepository() token.TokenRepository {
	if i.CurrTokenRepository != nil {
		return i.CurrTokenRepository
	}

	repository, err := token.NewTokenRepository(i.CurrLog, i.TokenCollection())
	if err != nil {
		i.CurrLog.Fatal(err)
		return nil
	}

	i.CurrTokenRepository = repository
	return i.CurrTokenRepository
}

func (i *Deps) TokenCache() token.TokenCache {
	if i.CurrTokenCache != nil {
		return i.CurrTokenCache
	}

	if tokenCache != nil {
		return tokenCache
	}

	tokenCache = token.NewTokenCache()

	return tokenCache
}

func (i *Deps) UserService() user.UserService {
	if i.CurrUserService != nil {
		return i.CurrUserService
	}

	i.CurrUserService = user.NewUserService(i.UserRepository())

	return i.CurrUserService
}

func (i *Deps) UserRepository() user.UserRepository {
	if i.CurrUserRepository != nil {
		return i.CurrUserRepository
	}

	repository, err := user.NewUserRepository(i.CurrLog, i.UserCollection())
	if err != nil {
		i.CurrLog.Fatal(err)
		return nil
	}

	i.CurrUserRepository = repository
	return i.CurrUserRepository
}

func (i *Deps) SignInUseCase() usecases.SignInUseCase {
	if i.CurrSignInUseCase != nil {
		return i.CurrSignInUseCase
	}

	i.CurrSignInUseCase = usecases.NewSignInUseCase(i.UserService(), i.TokenService())

	return i.CurrSignInUseCase
}

func (i *Deps) SignUpUseCase() usecases.SignUpUseCase {
	if i.CurrSignUpUseCase != nil {
		return i.CurrSignUpUseCase
	}

	i.CurrSignUpUseCase = usecases.NewSignUpUseCase(i.UserService(), i.TokenService())

	return i.CurrSignUpUseCase
}

func (i *Deps) InvalidateTokenUseCase() usecases.InvalidateTokenUseCase {
	if i.CurrInvalidateTokenUseCase != nil {
		return i.CurrInvalidateTokenUseCase
	}

	i.CurrInvalidateTokenUseCase = usecases.NewInvalidateTokenUseCase(i.CurrLog, i.TokenService(), i.SendLogoutService())

	return i.CurrInvalidateTokenUseCase
}

func (i *Deps) SendLogoutService() rabbit.SendLogoutService {
	if i.CurrSendLogoutService != nil {
		return i.CurrSendLogoutService
	}

	i.CurrSendLogoutService, _ = rabbit.NewSendLogoutService(i.CurrLog, i.RabbitChannel())

	return i.CurrSendLogoutService
}

func (i *Deps) TokenCollection() db.Collection {
	if i.CurrTokenCollection != nil {
		return i.CurrTokenCollection
	}

	if tokenCollection != nil {
		return tokenCollection
	}

	tokenCollection, err := db.NewCollection(i.CurrLog, i.Database(), "tokens", "userId")
	if err != nil {
		i.CurrLog.Fatal(err)
		return nil
	}

	return tokenCollection
}

func (i *Deps) UserCollection() db.Collection {
	if i.CurrUserCollection != nil {
		return i.CurrUserCollection
	}

	if userCollection != nil {
		return userCollection
	}

	userCollection, err := db.NewCollection(i.CurrLog, i.Database(), "users")
	if err != nil {
		i.CurrLog.Fatal(err)
		return nil
	}

	return userCollection
}

func (i *Deps) RabbitChannel() rbt.RabbitChannel {
	if i.CurrRabbitChannel != nil {
		return i.CurrRabbitChannel
	}

	i.CurrRabbitChannel, _ = rbt.NewRabbitChannel(env.Get().RabbitURL, i.CurrLog)

	return i.CurrRabbitChannel
}

// IsDbTimeoutError función a llamar cuando se produce un error de db
func IsDbTimeoutError(err interface{}) {
	if err == topology.ErrServerSelectionTimeout {
		database = nil
	}
}
