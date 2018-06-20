package db

import (
	"context"

	"github.com/mongodb/mongo-go-driver/core/option"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// WrapCollection Instancia db.Collection haciendo un wrap a un mongo.Collection real.
func WrapCollection(coll *mongo.Collection) Collection {
	return &fakeCollection{
		c: coll,
	}
}

// Collection es un wrapper de mongo.Collection que nos permite mockear mongo.Collection
type Collection interface {
	Name() string
	InsertOne(ctx context.Context, document interface{},
		opts ...option.InsertOneOptioner) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, documents []interface{},
		opts ...option.InsertManyOptioner) (*mongo.InsertManyResult, error)
	DeleteOne(ctx context.Context, filter interface{},
		opts ...option.DeleteOptioner) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, filter interface{},
		opts ...option.DeleteOptioner) (*mongo.DeleteResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{},
		options ...option.UpdateOptioner) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{},
		opts ...option.UpdateOptioner) (*mongo.UpdateResult, error)
	ReplaceOne(ctx context.Context, filter interface{},
		replacement interface{}, opts ...option.ReplaceOptioner) (*mongo.UpdateResult, error)
	Aggregate(ctx context.Context, pipeline interface{},
		opts ...option.AggregateOptioner) (mongo.Cursor, error)
	Count(ctx context.Context, filter interface{},
		opts ...option.CountOptioner) (int64, error)
	Distinct(ctx context.Context, fieldName string, filter interface{},
		opts ...option.DistinctOptioner) ([]interface{}, error)
	Find(ctx context.Context, filter interface{},
		opts ...option.FindOptioner) (mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{},
		opts ...option.FindOneOptioner) Decoder
	FindOneAndDelete(ctx context.Context, filter interface{},
		opts ...option.FindOneAndDeleteOptioner) Decoder
	FindOneAndReplace(ctx context.Context, filter interface{},
		replacement interface{}, opts ...option.FindOneAndReplaceOptioner) Decoder
	FindOneAndUpdate(ctx context.Context, filter interface{},
		update interface{}, opts ...option.FindOneAndUpdateOptioner) Decoder
	Watch(ctx context.Context, pipeline interface{},
		opts ...option.ChangeStreamOptioner) (mongo.Cursor, error)
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
	opts ...option.InsertOneOptioner) (*mongo.InsertOneResult, error) {
	return c.c.InsertOne(ctx, document, opts...)
}

func (c *fakeCollection) InsertMany(
	ctx context.Context,
	documents []interface{},
	opts ...option.InsertManyOptioner) (*mongo.InsertManyResult, error) {
	return c.c.InsertMany(ctx, documents, opts...)
}

func (c *fakeCollection) DeleteOne(
	ctx context.Context,
	filter interface{},
	opts ...option.DeleteOptioner) (*mongo.DeleteResult, error) {
	return c.c.DeleteOne(ctx, filter, opts...)
}

func (c *fakeCollection) DeleteMany(
	ctx context.Context,
	filter interface{},
	opts ...option.DeleteOptioner) (*mongo.DeleteResult, error) {
	return c.c.DeleteMany(ctx, filter, opts...)
}

func (c *fakeCollection) UpdateOne(
	ctx context.Context,
	filter interface{},
	update interface{},
	options ...option.UpdateOptioner) (*mongo.UpdateResult, error) {
	return c.c.UpdateOne(ctx, filter, update, options...)
}

func (c *fakeCollection) UpdateMany(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...option.UpdateOptioner) (*mongo.UpdateResult, error) {
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
	opts ...option.AggregateOptioner) (mongo.Cursor, error) {
	return c.c.Aggregate(ctx, pipeline, opts...)
}

func (c *fakeCollection) Count(
	ctx context.Context,
	filter interface{},
	opts ...option.CountOptioner) (int64, error) {
	return c.c.Count(ctx, filter, opts...)
}

func (c *fakeCollection) Distinct(
	ctx context.Context,
	fieldName string,
	filter interface{},
	opts ...option.DistinctOptioner) ([]interface{}, error) {
	return c.c.Distinct(ctx, fieldName, filter, opts...)
}

func (c *fakeCollection) Find(
	ctx context.Context,
	filter interface{},
	opts ...option.FindOptioner) (mongo.Cursor, error) {
	return c.c.Find(ctx, filter, opts...)
}

func (c *fakeCollection) FindOne(
	ctx context.Context,
	filter interface{},
	opts ...option.FindOneOptioner) Decoder {
	return c.c.FindOne(ctx, filter, opts...)
}

func (c *fakeCollection) FindOneAndDelete(
	ctx context.Context,
	filter interface{},
	opts ...option.FindOneAndDeleteOptioner) Decoder {
	return c.c.FindOneAndDelete(ctx, filter, opts...)
}

func (c *fakeCollection) FindOneAndReplace(
	ctx context.Context,
	filter interface{},
	replacement interface{},
	opts ...option.FindOneAndReplaceOptioner) Decoder {
	return c.c.FindOneAndReplace(ctx, filter, replacement, opts...)
}

func (c *fakeCollection) FindOneAndUpdate(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...option.FindOneAndUpdateOptioner) Decoder {
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
