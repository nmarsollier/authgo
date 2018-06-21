package mocks

import (
	"context"
	"errors"

	"github.com/nmarsollier/authgo/tools/db"

	"github.com/mongodb/mongo-go-driver/core/option"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/stretchr/testify/mock"
)

// FakeCollection es una implamentacion falsa de db.Collection que nos permite testear el sistema.
type FakeCollection struct {
	mock.Mock
}

// Name mocked
func (mc *FakeCollection) Name() string {
	return "MockedConnection"
}

// InsertOne mocked
func (mc *FakeCollection) InsertOne(
	ctx context.Context,
	document interface{},
	opts ...option.InsertOneOptioner) (*mongo.InsertOneResult, error) {

	res := mc.Called(ctx, document, opts)

	err, _ := res.Get(1).(error)

	return &mongo.InsertOneResult{
			InsertedID: res.Get(0),
		},
		err
}

// InsertMany mocked
func (mc *FakeCollection) InsertMany(
	ctx context.Context,
	documents []interface{},
	opts ...option.InsertManyOptioner) (*mongo.InsertManyResult, error) {

	res := mc.Called(ctx, documents, opts)

	err, _ := res.Get(1).(error)
	data, _ := res.Get(0).([]interface{})

	return &mongo.InsertManyResult{
			InsertedIDs: data,
		},
		err
}

// DeleteOne mocked
func (mc *FakeCollection) DeleteOne(ctx context.Context,
	filter interface{},
	opts ...option.DeleteOptioner) (*mongo.DeleteResult, error) {

	res := mc.Called(ctx, filter, opts)

	err, _ := res.Get(1).(error)
	data, _ := res.Get(0).(int64)

	return &mongo.DeleteResult{
			DeletedCount: data,
		},
		err
}

// DeleteMany mocked
func (mc *FakeCollection) DeleteMany(ctx context.Context,
	filter interface{},
	opts ...option.DeleteOptioner) (*mongo.DeleteResult, error) {

	res := mc.Called(ctx, filter, opts)

	err, _ := res.Get(1).(error)
	data, _ := res.Get(0).(int64)

	return &mongo.DeleteResult{
			DeletedCount: data,
		},
		err
}

// UpdateOne mocked
func (mc *FakeCollection) UpdateOne(ctx context.Context,
	filter interface{},
	update interface{},
	options ...option.UpdateOptioner) (*mongo.UpdateResult, error) {
	res := mc.Called(ctx, filter, options)

	matchedCount, _ := res.Get(0).(int64)
	modifiedCount, _ := res.Get(1).(int64)
	err, _ := res.Get(3).(error)

	return &mongo.UpdateResult{
			MatchedCount:  matchedCount,
			ModifiedCount: modifiedCount,
			UpsertedID:    res.Get(2),
		},
		err
}

// UpdateMany mocked
func (mc *FakeCollection) UpdateMany(ctx context.Context,
	filter interface{},
	update interface{},
	opts ...option.UpdateOptioner) (*mongo.UpdateResult, error) {

	res := mc.Called(ctx, filter, update, opts)

	matchedCount, _ := res.Get(0).(int64)
	modifiedCount, _ := res.Get(1).(int64)
	err, _ := res.Get(3).(error)

	return &mongo.UpdateResult{
			MatchedCount:  matchedCount,
			ModifiedCount: modifiedCount,
			UpsertedID:    res.Get(2),
		},
		err
}

// ReplaceOne mocked
func (mc *FakeCollection) ReplaceOne(ctx context.Context,
	filter interface{},
	replacement interface{},
	opts ...option.ReplaceOptioner) (*mongo.UpdateResult, error) {

	res := mc.Called(ctx, filter, replacement, opts)

	matchedCount, _ := res.Get(0).(int64)
	modifiedCount, _ := res.Get(1).(int64)
	err, _ := res.Get(3).(error)

	return &mongo.UpdateResult{
			MatchedCount:  matchedCount,
			ModifiedCount: modifiedCount,
			UpsertedID:    res.Get(2),
		},
		err
}

// Aggregate mocked
func (mc *FakeCollection) Aggregate(ctx context.Context,
	pipeline interface{},
	opts ...option.AggregateOptioner) (mongo.Cursor, error) {

	res := mc.Called(ctx, pipeline, opts)

	if err, ok := res.Get(1).(error); ok && err != nil {
		return nil, err
	}

	if mocked, ok := res.Get(0).(mongo.Cursor); ok {
		return mocked, nil
	}

	return nil, errors.New("MockedConnection.Aggregate - Not Implemented")
}

// Count mocked
func (mc *FakeCollection) Count(ctx context.Context,
	filter interface{},
	opts ...option.CountOptioner) (int64, error) {

	res := mc.Called(ctx, filter, opts)

	count, _ := res.Get(0).(int64)
	err, _ := res.Get(1).(error)

	return count, err
}

// Distinct mocked
func (mc *FakeCollection) Distinct(ctx context.Context,
	fieldName string,
	filter interface{},
	opts ...option.DistinctOptioner) ([]interface{}, error) {
	return nil, errors.New("MockedConnection.Distinct - Not Implemented")
}

// Find mocked
func (mc *FakeCollection) Find(ctx context.Context,
	filter interface{},
	opts ...option.FindOptioner) (mongo.Cursor, error) {

	res := mc.Called(ctx, filter, opts)

	if err, ok := res.Get(1).(error); ok && err != nil {
		return nil, err
	}

	if mocked, ok := res.Get(0).(mongo.Cursor); ok {
		return mocked, nil
	}

	return nil, errors.New("MockedConnection.Find - Not Implemented")
}

// FindOne mocked
func (mc *FakeCollection) FindOne(
	ctx context.Context,
	filter interface{},
	opts ...option.FindOneOptioner) db.Decoder {
	res := mc.Called(ctx, filter, opts)

	if mocked, ok := res.Get(0).(db.Decoder); ok {
		return mocked
	}
	return FakeDecoder(nil)
}

// FindOneAndDelete mocked
func (mc *FakeCollection) FindOneAndDelete(
	ctx context.Context, filter interface{},
	opts ...option.FindOneAndDeleteOptioner) db.Decoder {

	res := mc.Called(ctx, filter, opts)
	if mocked, ok := res.Get(0).(db.Decoder); ok {
		return mocked
	}
	return FakeDecoder(nil)
}

// FindOneAndReplace mocked
func (mc *FakeCollection) FindOneAndReplace(
	ctx context.Context,
	filter interface{},
	replacement interface{},
	opts ...option.FindOneAndReplaceOptioner) db.Decoder {

	res := mc.Called(ctx, filter, opts)
	if mocked, ok := res.Get(0).(db.Decoder); ok {
		return mocked
	}
	return FakeDecoder(nil)
}

// FindOneAndUpdate mocked
func (mc *FakeCollection) FindOneAndUpdate(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...option.FindOneAndUpdateOptioner) db.Decoder {

	res := mc.Called(ctx, filter, opts)
	if mocked, ok := res.Get(0).(db.Decoder); ok {
		return mocked
	}
	return FakeDecoder(nil)
}

// Watch mocked
func (mc *FakeCollection) Watch(
	ctx context.Context,
	pipeline interface{},
	opts ...option.ChangeStreamOptioner) (mongo.Cursor, error) {

	res := mc.Called(ctx, pipeline, opts)

	if err, ok := res.Get(1).(error); ok && err != nil {
		return nil, err
	}

	if mocked, ok := res.Get(0).(mongo.Cursor); ok {
		return mocked, nil
	}

	return nil, errors.New("MockedConnection.Watch - Not Implemented")
}

// Indexes mocked
func (mc *FakeCollection) Indexes() mongo.IndexView {
	return mongo.IndexView{}
}

// Drop mocked
func (mc *FakeCollection) Drop(ctx context.Context) error {
	return errors.New("MockedConnection.Drop - Not Implemented")
}
