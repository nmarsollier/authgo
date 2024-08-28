mockgen -source=./rabbit/rabbit_channel.go -destination=./rabbit/rabbit_channel_mocks.go -package=rabbit
mockgen -source=./tools/db/mongo_collection.go -destination=./tools/db/mongo_collection_mocks.go -package=db
mockgen -source=./log/logrus_logger.go -destination=./log/logrus_logger_mocks.go -package=log