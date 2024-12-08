package user

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/tools/log"
)

var tableName = "users"

type DbUserUpdateDocumentBody struct {
	Name        string `validate:"required,min=1,max=100"`
	Password    string `validate:"required"`
	Permissions []string
	Enabled     bool
	Updated     time.Time
}

func insert(user *User, deps ...interface{}) (_ *User, err error) {
	if err = user.validateSchema(); err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	userToInsert, err := attributevalue.MarshalMap(user)
	if err != nil {
		log.Get(deps...).Error(err)
		return
	}

	_, err = db.Get(deps...).PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			TableName: &tableName,
			Item:      userToInsert,
		},
	)

	if err != nil {
		log.Get(deps...).Error(err)
		return
	}

	return user, nil
}

func update(user *User, deps ...interface{}) (err error) {
	if err = user.validateSchema(); err != nil {
		log.Get(deps...).Error(err)
		return
	}

	userKey, err := attributevalue.MarshalMap(map[string]interface{}{
		"id": user.ID,
	})
	if err != nil {
		log.Get(deps...).Error(err)
		return
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

	_, err = db.Get(deps...).UpdateItem(
		context.TODO(),
		&dynamodb.UpdateItemInput{
			TableName:                 &tableName,
			Key:                       userKey,
			UpdateExpression:          aws.String("SET name = :name, password = :password, permissions = :permissions, enabled = :enabled, updated = :updated"),
			ExpressionAttributeValues: update,
		},
	)

	if err != nil {
		log.Get(deps...).Error(err)
	}

	return nil
}

// FindAll devuelve todos los usuarios
func findAll(deps ...interface{}) (users []*User, err error) {
	result, err := db.Get(deps...).Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: &tableName,
	})
	if err != nil {
		log.Get(deps...).Error(err)
		return
	}

	err = attributevalue.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		log.Get(deps...).Error(err)
	}

	return
}

// FindByID lee un usuario desde la db
func findByID(userID string, deps ...interface{}) (user *User, err error) {
	response, err := db.Get(deps...).GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: userID,
			}},
		TableName: &tableName,
	})

	if err != nil || response == nil || response.Item == nil {
		log.Get(deps...).Error(err)

		return nil, errs.NotFound
	}

	err = attributevalue.UnmarshalMap(response.Item, &user)
	if err != nil {
		log.Get(deps...).Error(err)
	}

	return
}

// FindByLogin lee un usuario desde la db
func findByLogin(login string, deps ...interface{}) (user *User, err error) {
	expr, err := expression.NewBuilder().WithKeyCondition(
		expression.Key("login").Equal(expression.Value(login)),
	).Build()
	if err != nil {
		return nil, err
	}

	response, err := db.Get(deps...).Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 &tableName,
		IndexName:                 aws.String("login-index"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	if err != nil || len(response.Items) == 0 {
		log.Get(deps...).Error(err)
		return
	}

	err = attributevalue.UnmarshalMap(response.Items[0], &user)
	if err != nil {
		log.Get(deps...).Error(err)
	}

	return
}
