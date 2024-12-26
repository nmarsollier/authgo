rm -rf ./test/mockgen

mockgen -source=./internal/engine/db/collection.go -destination=./test/mockgen/mongo_collection_mocks.go -package=mockgen
mockgen -source=./internal/engine/di/injector.go -destination=./test/mockgen/injector_mocks.go -package=mockgen
mockgen -source=./internal/engine/log/logger.go -destination=./test/mockgen/logger_mocks.go -package=mockgen

mockgen -source=./internal/engine/rbt/rabbit_channel.go -destination=./test/mockgen/rabbit_channel_mocks.go -package=mockgen
mockgen -source=./internal/rabbit/send_logout.go -destination=./test/mockgen/send_logout_mocks.go -package=mockgen

mockgen -source=./internal/token/cache.go -destination=./test/mockgen/cache_mocks.go -package=mockgen
mockgen -source=./internal/token/repository.go -destination=./test/mockgen/token_repository_mocks.go -package=mockgen
mockgen -source=./internal/token/service.go -destination=./test/mockgen/token_service_mocks.go -package=mockgen

mockgen -source=./internal/usecases/invalidate_token.go -destination=./test/mockgen/invalidate_token_mocks.go -package=mockgen
mockgen -source=./internal/usecases/sign_in.go -destination=./test/mockgen/sign_in_mocks.go -package=mockgen
mockgen -source=./internal/usecases/sign_up.go -destination=./test/mockgen/sign_up_mocks.go -package=mockgen

mockgen -source=./internal/user/repository.go -destination=./test/mockgen/user_repository_mocks.go -package=mockgen
mockgen -source=./internal/user/service.go -destination=./test/mockgen/user_service_mocks.go -package=mockgen
