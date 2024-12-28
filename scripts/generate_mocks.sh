rm -rf ./test/mockgen

mockgen -source=./internal/token/repository.go -destination=./test/mockgen/token_repository_mocks.go -package=mockgen
mockgen -source=./internal/token/service.go -destination=./test/mockgen/token_service_mocks.go -package=mockgen

mockgen -source=./internal/usecases/invalidate_token.go -destination=./test/mockgen/invalidate_token_mocks.go -package=mockgen
mockgen -source=./internal/usecases/sign_in.go -destination=./test/mockgen/sign_in_mocks.go -package=mockgen
mockgen -source=./internal/usecases/sign_up.go -destination=./test/mockgen/sign_up_mocks.go -package=mockgen

mockgen -source=./internal/user/repository.go -destination=./test/mockgen/user_repository_mocks.go -package=mockgen
mockgen -source=./internal/user/service.go -destination=./test/mockgen/user_service_mocks.go -package=mockgen
