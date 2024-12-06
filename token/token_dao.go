package token

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/nmarsollier/authgo/tools/env"
)

var tableName = "tokens"

var (
	once     sync.Once
	instance TokenDao
)

type TokenDao interface {
	FindById(id string) (token *Token, err error)

	Insert(user *Token) (err error)

	Delete(id string) (err error)
}

func GetTokenDao(deps ...interface{}) (instance TokenDao, err error) {
	for _, o := range deps {
		if client, ok := o.(TokenDao); ok {
			return client, nil
		}
	}

	once.Do(func() {
		customCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			env.Get().AwsAccessKeyId,
			env.Get().AwsSecret,
			"",
		))

		cfg, e := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(env.Get().AwsRegion),
			config.WithCredentialsProvider(customCreds),
		)
		if e != nil {
			err = e
		}

		client := dynamodb.NewFromConfig(cfg)

		instance = &tokenDao{
			client: client,
		}
	})

	return
}

type tokenDao struct {
	client *dynamodb.Client
}

func (r *tokenDao) FindById(id string) (*Token, error) {
	token := Token{ID: id}
	tokenId, err := attributevalue.Marshal(token.ID)
	if err != nil {
		return nil, err
	}

	response, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{"id": tokenId}, TableName: &tableName,
	})

	if err != nil || response == nil || response.Item == nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(response.Item, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *tokenDao) Insert(user *Token) error {
	tokenToInsert, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			TableName: &tableName,
			Item:      tokenToInsert,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *tokenDao) Delete(id string) error {
	token := Token{ID: id}
	tokenId, err := attributevalue.Marshal(token.ID)
	if err != nil {
		return err
	}

	_, err = r.client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{"id": tokenId}, TableName: &tableName,
	})

	if err != nil {
		return err
	}

	return nil
}
