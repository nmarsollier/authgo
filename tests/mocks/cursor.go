package mocks

import (
	"context"
	"encoding/json"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// Cursor un fake cursor para la data indicada
func Cursor(data []interface{}) mongo.Cursor {
	cursor := fakeCursor{
		data: data,
		idx:  -1,
	}
	return &cursor
}

type fakeCursor struct {
	data []interface{}
	idx  int
}

func (f *fakeCursor) ID() int64 {
	return 0
}

func (f *fakeCursor) Next(ctx context.Context) bool {
	f.idx = f.idx + 1
	return f.idx < len(f.data)
}

func (f *fakeCursor) Decode(v interface{}) error {
	jsonData, err := json.Marshal(f.data[f.idx])
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, &v)
	if err != nil {
		return err
	}

	return nil
}

func (f *fakeCursor) DecodeBytes() (reader bson.Reader, err error) {
	return nil, nil
}

func (f *fakeCursor) Err() error {
	return nil
}

func (f *fakeCursor) Close(ctx context.Context) error {
	return nil
}
