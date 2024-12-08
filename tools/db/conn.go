package db

import (
	"context"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/nmarsollier/authgo/tools/env"
)

var (
	once     sync.Once
	instance *dynamodb.Client
)

func Get(deps ...interface{}) *dynamodb.Client {
	once.Do(func() {
		customCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			env.Get().AwsAccessKeyId,
			env.Get().AwsSecret,
			"",
		))

		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(env.Get().AwsRegion),
			config.WithCredentialsProvider(customCreds),
		)
		if err != nil {
			log.Fatalf("Error cargando la configuración: %v", err)
		}

		client := dynamodb.NewFromConfig(cfg)

		instance = client
	})

	return instance
}
