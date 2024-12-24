mockgen -source=./engine/db/mongo_collection.go -destination=./tests/engine/db/mongo_collection_mocks.go -package=db
mockgen -source=./engine/di/injector.go -destination=./tests/engine/di/injector_mocks.go -package=di
mockgen -source=./engine/log/logger.go -destination=./tests/engine/log/logger_mocks.go -package=log

mockgen -source=./rabbit/rabbit_channel.go -destination=./tests/rabbit/rabbit_channel_mocks.go -package=rabbit
mockgen -source=./rabbit/send_logout.go -destination=./tests/rabbit/send_logout_mocks.go -package=rabbit

mockgen -source=./token/cache.go -destination=./tests/token/cache_mocks.go -package=token
mockgen -source=./token/repository.go -destination=./tests/token/repository_mocks.go -package=token
mockgen -source=./token/service.go -destination=./tests/token/service_mocks.go -package=token

mockgen -source=./usecases/invalidate_token.go -destination=./tests/usecases/invalidate_token_mocks.go -package=usecases
mockgen -source=./usecases/sign_in.go -destination=./tests/usecases/sign_in_mocks.go -package=usecases
mockgen -source=./usecases/sign_up.go -destination=./tests/usecases/sign_up_mocks.go -package=usecases

mockgen -source=./user/repository.go -destination=./tests/user/repository_mocks.go -package=user
mockgen -source=./user/service.go -destination=./tests/user/service_mocks.go -package=user
