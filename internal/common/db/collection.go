package db

import (
	"context"

	"github.com/nmarsollier/authgo/internal/common/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection interface {
	FindOne(ctx context.Context, filter interface{}, v interface{}) error

	InsertOne(ctx context.Context, document interface{}) (id interface{}, error error)

	UpdateOne(ctx context.Context, filter interface{}, update interface{}, optn *options.UpdateOptions) (modified int64, error error)

	Find(ctx context.Context, filter interface{}) (cur Cursor, err error)

	ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}) (modified int64, error error)
}

func NewCollection(
	log log.LogRusEntry,
	database *mongo.Database,
	collectionName string,
	onError func(error),
	indexes ...string,
) (col Collection, err error) {
	collection := database.Collection(collectionName)

	for _, index := range indexes {
		_, err = collection.Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys:    bson.M{index: 1},
				Options: nil,
			},
		)
		if err != nil {
			onError(err)
			log.Info(err)
			return nil, err
		}
	}
	if err != nil {
		onError(err)
		log.Error(err)
		return nil, err
	}

	return &mongoCollection{
		collection: collection,
		onError:    onError,
	}, nil
}

type mongoCollection struct {
	collection *mongo.Collection
	onError    func(error)
}

func (m *mongoCollection) FindOne(ctx context.Context, filter interface{}, v interface{}) error {
	if err := m.collection.FindOne(context.Background(), filter).Decode(v); err != nil {
		m.onError(err)
		return err
	}
	return nil
}

func (m *mongoCollection) InsertOne(ctx context.Context, document interface{}) (id interface{}, error error) {
	insertedId, err := m.collection.InsertOne(context.Background(), document)
	if err != nil {
		m.onError(err)
		return nil, err
	}
	return insertedId.InsertedID, nil
}

func (m *mongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, optn *options.UpdateOptions) (modified int64, error error) {
	insertedId, err := m.collection.UpdateOne(context.Background(), filter, update, optn)
	if err != nil {
		m.onError(err)
		return 0, err
	}
	return insertedId.ModifiedCount, nil
}

func (m *mongoCollection) Find(ctx context.Context, filter interface{}) (cur Cursor, err error) {
	cursor, err := m.collection.Find(context.Background(), filter)
	if err != nil {
		m.onError(err)
		return nil, err
	}
	return NewCursor(cursor), nil
}

type Cursor interface {
	Close(ctx context.Context) error
	Next(ctx context.Context) bool
	Decode(val interface{}) error
}

func NewCursor(cursor *mongo.Cursor) Cursor {
	return &mongoCursor{
		cursor: cursor,
	}
}

type mongoCursor struct {
	cursor *mongo.Cursor
}

func (c *mongoCursor) Close(ctx context.Context) error {
	return c.cursor.Close(ctx)
}

func (c *mongoCursor) Next(ctx context.Context) bool {
	return c.cursor.Next(ctx)
}
func (c *mongoCursor) Decode(val interface{}) error {
	return c.cursor.Decode(val)
}

func (m *mongoCollection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}) (modified int64, error error) {
	insertedId, err := m.collection.ReplaceOne(context.Background(), filter, replacement)
	if err != nil {
		return 0, err
	}
	return insertedId.ModifiedCount, nil
}
