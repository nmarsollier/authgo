mockgen -source=./rabbit/rabbit.go -destination=./rabbit/mocks.go -package=rabbit
mockgen -source=./tools/db/mongo_collection.go -destination=./tools/db/mocks.go -package=db