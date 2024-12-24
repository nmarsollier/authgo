package user

import (
	"context"

	"github.com/nmarsollier/authgo/engine/db"
	"github.com/nmarsollier/authgo/engine/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewCollection(log log.LogRusEntry) (db.MongoCollection, error) {
	database, err := db.CurrentDatabase()
	if err != nil {
		return nil, err
	}

	col := database.Collection("users")

	_, err = col.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    DbUserLoginFilter{Login: ""},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Error(err)
	}

	return db.NewMongoCollection(col), nil
}
