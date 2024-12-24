package di

import (
	"github.com/nmarsollier/authgo/engine/db"
	"github.com/nmarsollier/authgo/engine/log"
	"github.com/nmarsollier/authgo/rabbit"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/usecases"
	"github.com/nmarsollier/authgo/user"
)

// Singletons
var tokenCache token.TokenCache
var tokenCollection db.MongoCollection
var userCollection db.MongoCollection

type Injector interface {
	InvalidateTokenUseCase() usecases.InvalidateTokenUseCase
	Logger() log.LogRusEntry
	RabbitChannel() rabbit.RabbitChannel
	SendLogoutService() rabbit.SendLogoutService
	SignInUseCase() usecases.SignInUseCase
	SignUpUseCase() usecases.SignUpUseCase
	TokenCache() token.TokenCache
	TokenCollection() db.MongoCollection
	TokenRepository() token.TokenRepository
	TokenService() token.TokenService
	UserCollection() db.MongoCollection
	UserRepository() user.UserRepository
	UserService() user.UserService
}

type injector struct {
	invalidateTokenUseCase usecases.InvalidateTokenUseCase
	log                    log.LogRusEntry
	sendLogoutService      rabbit.SendLogoutService
	signInUseCase          usecases.SignInUseCase
	signUpUseCase          usecases.SignUpUseCase
	tokenRepository        token.TokenRepository
	tokenService           token.TokenService
	userRepository         user.UserRepository
	userService            user.UserService
	rabbitChannel          rabbit.RabbitChannel
	tokenCache             token.TokenCache
	tokenCollection        db.MongoCollection
	userCollection         db.MongoCollection
}

func NewInjector(log log.LogRusEntry) Injector {
	return &injector{
		log: log,
	}
}

func (i *injector) Logger() log.LogRusEntry {
	return i.log
}

func (i *injector) TokenService() token.TokenService {
	if i.tokenService != nil {
		return i.tokenService
	}

	i.tokenService = token.NewTokenService(i.log, i.TokenCache(), i.TokenRepository())
	return i.tokenService
}

func (i *injector) TokenRepository() token.TokenRepository {
	if i.tokenRepository != nil {
		return i.tokenRepository
	}

	repository, err := token.NewTokenRepository(i.log, i.TokenCollection())
	if err != nil {
		i.log.Fatal(err)
		return nil
	}

	i.tokenRepository = repository
	return i.tokenRepository
}

func (i *injector) TokenCache() token.TokenCache {
	if i.tokenCache != nil {
		return i.tokenCache
	}

	if tokenCache != nil {
		return tokenCache
	}

	tokenCache = token.NewTokenCache()

	return tokenCache
}

func (i *injector) UserService() user.UserService {
	if i.userService != nil {
		return i.userService
	}

	i.userService = user.NewUserService(i.log, i.UserRepository())

	return i.userService
}

func (i *injector) UserRepository() user.UserRepository {
	if i.userRepository != nil {
		return i.userRepository
	}

	repository, err := user.NewUserRepository(i.log, i.UserCollection())
	if err != nil {
		i.log.Fatal(err)
		return nil
	}

	i.userRepository = repository
	return i.userRepository
}

func (i *injector) SignInUseCase() usecases.SignInUseCase {
	if i.signInUseCase != nil {
		return i.signInUseCase
	}

	i.signInUseCase = usecases.NewSignInUseCase(i.UserService(), i.TokenService())

	return i.signInUseCase
}

func (i *injector) SignUpUseCase() usecases.SignUpUseCase {
	if i.signUpUseCase != nil {
		return i.signUpUseCase
	}

	i.signUpUseCase = usecases.NewSignUpUseCase(i.UserService(), i.TokenService())

	return i.signUpUseCase
}

func (i *injector) InvalidateTokenUseCase() usecases.InvalidateTokenUseCase {
	if i.invalidateTokenUseCase != nil {
		return i.invalidateTokenUseCase
	}

	i.invalidateTokenUseCase = usecases.NewInvalidateTokenUseCase(i.log, i.TokenService(), i.SendLogoutService())

	return i.invalidateTokenUseCase
}

func (i *injector) SendLogoutService() rabbit.SendLogoutService {
	if i.sendLogoutService != nil {
		return i.sendLogoutService
	}

	i.sendLogoutService, _ = rabbit.NewSendLogoutService(i.log, i.RabbitChannel())

	return i.sendLogoutService
}

func (i *injector) TokenCollection() db.MongoCollection {
	if i.tokenCollection != nil {
		return i.tokenCollection
	}

	if tokenCollection != nil {
		return tokenCollection
	}

	tokenCollection, err := token.NewCollection(i.log)
	if err != nil {
		i.log.Fatal(err)
		return nil
	}

	return tokenCollection
}

func (i *injector) UserCollection() db.MongoCollection {
	if i.userCollection != nil {
		return i.userCollection
	}

	if userCollection != nil {
		return userCollection
	}

	userCollection, err := user.NewCollection(i.log)
	if err != nil {
		i.log.Fatal(err)
		return nil
	}

	return userCollection
}

func (i *injector) RabbitChannel() rabbit.RabbitChannel {
	if i.rabbitChannel != nil {
		return i.rabbitChannel
	}

	i.rabbitChannel, _ = rabbit.NewChannel(i.log)

	return i.rabbitChannel
}
