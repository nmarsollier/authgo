package user

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

var tableName = "users"

var (
	once     sync.Once
	instance UserDao
)

type UserDao interface {
	FindById(key string) (*User, error)

	FindByLogin(login string) (*User, error)

	Insert(user *User) error

	Update(user *User) error

	FindAll() ([]*User, error)
}

func GetUserDao(deps ...interface{}) (instance UserDao, err error) {
	for _, o := range deps {
		if client, ok := o.(UserDao); ok {
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

		instance = &authDao{
			client: client,
		}
	})

	return
}

type authDao struct {
	client *dynamodb.Client
}

func (r *authDao) FindById(key string) (*User, error) {
	user := User{ID: key}
	userId, err := attributevalue.Marshal(user.ID)
	if err != nil {
		return nil, err
	}

	response, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{"id": userId}, TableName: &tableName,
	})

	if err != nil || response == nil || response.Item == nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(response.Item, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authDao) FindByLogin(login string) (*User, error) {
	user := User{Login: login}
	loginAttr, err := attributevalue.Marshal(user.ID)
	if err != nil {
		return nil, err
	}

	response, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{"login": loginAttr}, TableName: &tableName,
	})

	if err != nil || response == nil || response.Item == nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(response.Item, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authDao) Insert(user *User) error {
	userToInsert, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			TableName: &tableName,
			Item:      userToInsert,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *authDao) Update(user *User) error {
	userKey, err := attributevalue.MarshalMap(map[string]interface{}{
		"ID": user.ID,
	})
	if err != nil {
		return err
	}

	update, err := attributevalue.MarshalMap(map[string]interface{}{
		":name":        user.Name,
		":password":    user.Password,
		":permissions": user.Permissions,
		":enabled":     user.Enabled,
		":updated":     user.Updated,
	})
	if err != nil {
		return err
	}

	_, err = r.client.UpdateItem(
		context.TODO(),
		&dynamodb.UpdateItemInput{
			TableName:                 &tableName,
			Key:                       userKey,
			UpdateExpression:          aws.String("SET Name = :name, Password = :password, Permissions = :permissions, Enabled = :enabled, Updated = :updated"),
			ExpressionAttributeValues: update,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *authDao) FindAll() ([]*User, error) {
	result, err := r.client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: &tableName,
	})
	if err != nil {
		return nil, err
	}

	var users []*User
	err = attributevalue.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
