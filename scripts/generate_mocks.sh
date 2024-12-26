rm -rf ./test/mock

mockgen -source=./internal/engine/db/collection.go -destination=./test/mock/mongo_collection_mocks.go -package=mock
mockgen -source=./internal/engine/di/injector.go -destination=./test/mock/injector_mocks.go -package=mock
mockgen -source=./internal/engine/log/logger.go -destination=./test/mock/logger_mocks.go -package=mock

mockgen -source=./internal/engine/rbt/rabbit_channel.go -destination=./test/mock/rabbit_channel_mocks.go -package=mock
mockgen -source=./internal/rabbit/send_logout.go -destination=./test/mock/send_logout_mocks.go -package=mock

mockgen -source=./internal/token/cache.go -destination=./test/mock/cache_mocks.go -package=mock
mockgen -source=./internal/token/repository.go -destination=./test/mock/token_repository_mocks.go -package=mock
mockgen -source=./internal/token/service.go -destination=./test/mock/token_service_mocks.go -package=mock

mockgen -source=./internal/usecases/invalidate_token.go -destination=./test/mock/invalidate_token_mocks.go -package=mock
mockgen -source=./internal/usecases/sign_in.go -destination=./test/mock/sign_in_mocks.go -package=mock
mockgen -source=./internal/usecases/sign_up.go -destination=./test/mock/sign_up_mocks.go -package=mock

mockgen -source=./internal/user/repository.go -destination=./test/mock/user_repository_mocks.go -package=mock
mockgen -source=./internal/user/service.go -destination=./test/mock/user_service_mocks.go -package=mock
