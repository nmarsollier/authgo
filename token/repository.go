package token

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/tools/log"
)

var tableName = "tokens"

// insert crea un nuevo token y lo almacena en la db
func insert(userID string, deps ...interface{}) (*Token, error) {
	token := newToken(userID)

	tokenToInsert, err := attributevalue.MarshalMap(token)
	if err != nil {
		log.Get(deps...).Error(err)

		return nil, err
	}

	_, err = db.Get(deps...).PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			TableName: &tableName,
			Item:      tokenToInsert,
		},
	)
	if err != nil {
		log.Get(deps...).Error(err)

		return nil, err
	}

	return token, nil
}

// findByID busca un token en la db
func findByID(tokenID string, deps ...interface{}) (token *Token, err error) {
	response, err := db.Get().GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: tokenID,
			}},
		TableName: &tableName,
	})

	if err != nil || response == nil || response.Item == nil {
		log.Get(deps...).Error(err)

		return nil, errs.NotFound
	}

	err = attributevalue.UnmarshalMap(response.Item, &token)
	return
}

// delete como deshabilitado un token
func delete(tokenID string, deps ...interface{}) (err error) {
	_, err = db.Get(deps...).DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{
			Value: tokenID,
		}}, TableName: &tableName,
	})

	if err != nil {
		log.Get(deps...).Error(err)
	}
	return
}
