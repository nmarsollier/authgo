package db

import (
	"context"

	"github.com/mongodb/mongo-go-driver/core/option"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/aggregateopt"
	"github.com/mongodb/mongo-go-driver/mongo/countopt"
	"github.com/mongodb/mongo-go-driver/mongo/deleteopt"
	"github.com/mongodb/mongo-go-driver/mongo/distinctopt"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/insertopt"
	"github.com/mongodb/mongo-go-driver/mongo/updateopt"
)

// Decoder permite mockear mongo.DocumentResult, para poder testear la app
type Decoder interface {
	Decode(v interface{}) error
}

// WrapCollection Instancia db.Collection haciendo un wrap a un mongo.Collection real.
func WrapCollection(coll *mongo.Collection) Collection {
	return &fakeCollection{
		c: coll,
	}
}

// Collection es un wrapper de mongo.Collection que nos permite mockear mongo.Collection
type Collection interface {
	Name() string
	InsertOne(ctx context.Context, document interface{}, opts ...insertopt.One) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, documents []interface{}, opts ...insertopt.Many) (*mongo.InsertManyResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...deleteopt.Delete) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...deleteopt.Delete) (*mongo.DeleteResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, options ...updateopt.Update) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...updateopt.Update) (*mongo.UpdateResult, error)
	ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...option.ReplaceOptioner) (*mongo.UpdateResult, error)
	Aggregate(ctx context.Context, pipeline interface{}, opts ...aggregateopt.Aggregate) (mongo.Cursor, error)
	Count(ctx context.Context, filter interface{}, opts ...countopt.Count) (int64, error)
	Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...distinctopt.Distinct) ([]interface{}, error)
	Find(ctx context.Context, filter interface{}, opts ...findopt.Find) (mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...findopt.One) Decoder
	FindOneAndDelete(ctx context.Context, filter interface{}, opts ...findopt.DeleteOne) Decoder
	FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...findopt.ReplaceOne) Decoder
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...findopt.UpdateOne) Decoder
	Watch(ctx context.Context, pipeline interface{}, opts ...option.ChangeStreamOptioner) (mongo.Cursor, error)
	Indexes() mongo.IndexView
	Drop(ctx context.Context) error
}

type fakeCollection struct {
	c *mongo.Collection
}

func (c *fakeCollection) Name() string {
	return c.c.Name()
}

func (c *fakeCollection) InsertOne(
	ctx context.Context,
	document interface{},
	opts ...insertopt.One) (*mongo.InsertOneResult, error) {
	return c.c.InsertOne(ctx, document, opts...)
}

func (c *fakeCollection) InsertMany(
	ctx context.Context,
	documents []interface{},
	opts ...insertopt.Many) (*mongo.InsertManyResult, error) {
	return c.c.InsertMany(ctx, documents, opts...)
}

func (c *fakeCollection) DeleteOne(
	ctx context.Context,
	filter interface{},
	opts ...deleteopt.Delete) (*mongo.DeleteResult, error) {
	return c.c.DeleteOne(ctx, filter, opts...)
}

func (c *fakeCollection) DeleteMany(
	ctx context.Context,
	filter interface{},
	opts ...deleteopt.Delete) (*mongo.DeleteResult, error) {
	return c.c.DeleteMany(ctx, filter, opts...)
}

func (c *fakeCollection) UpdateOne(
	ctx context.Context,
	filter interface{},
	update interface{},
	options ...updateopt.Update) (*mongo.UpdateResult, error) {
	return c.c.UpdateOne(ctx, filter, update, options...)
}

func (c *fakeCollection) UpdateMany(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...updateopt.Update) (*mongo.UpdateResult, error) {
	return c.c.UpdateMany(ctx, filter, update, opts...)
}

func (c *fakeCollection) ReplaceOne(
	ctx context.Context,
	filter interface{},
	replacement interface{},
	opts ...option.ReplaceOptioner) (*mongo.UpdateResult, error) {
	return c.c.ReplaceOne(ctx, replacement, opts)
}

func (c *fakeCollection) Aggregate(
	ctx context.Context,
	pipeline interface{},
	opts ...aggregateopt.Aggregate) (mongo.Cursor, error) {
	return c.c.Aggregate(ctx, pipeline, opts...)
}

func (c *fakeCollection) Count(
	ctx context.Context,
	filter interface{},
	opts ...countopt.Count) (int64, error) {
	return c.c.Count(ctx, filter, opts...)
}

func (c *fakeCollection) Distinct(
	ctx context.Context,
	fieldName string,
	filter interface{},
	opts ...distinctopt.Distinct) ([]interface{}, error) {
	return c.c.Distinct(ctx, fieldName, filter, opts...)
}

func (c *fakeCollection) Find(
	ctx context.Context,
	filter interface{},
	opts ...findopt.Find) (mongo.Cursor, error) {
	return c.c.Find(ctx, filter, opts...)
}

func (c *fakeCollection) FindOne(
	ctx context.Context,
	filter interface{},
	opts ...findopt.One) Decoder {
	return c.c.FindOne(ctx, filter, opts...)
}

func (c *fakeCollection) FindOneAndDelete(
	ctx context.Context,
	filter interface{},
	opts ...findopt.DeleteOne) Decoder {
	return c.c.FindOneAndDelete(ctx, filter, opts...)
}

func (c *fakeCollection) FindOneAndReplace(
	ctx context.Context,
	filter interface{},
	replacement interface{},
	opts ...findopt.ReplaceOne) Decoder {
	return c.c.FindOneAndReplace(ctx, filter, replacement, opts...)
}

func (c *fakeCollection) FindOneAndUpdate(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...findopt.UpdateOne) Decoder {
	return c.c.FindOneAndUpdate(ctx, filter, update, opts...)
}

func (c *fakeCollection) Watch(
	ctx context.Context,
	pipeline interface{},
	opts ...option.ChangeStreamOptioner) (mongo.Cursor, error) {
	return c.c.Watch(ctx, pipeline)
}

func (c *fakeCollection) Indexes() mongo.IndexView {
	return c.c.Indexes()
}

func (c *fakeCollection) Drop(ctx context.Context) error {
	return c.c.Drop(ctx)
}
