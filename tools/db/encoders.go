package db

import "github.com/mongodb/mongo-go-driver/bson"

// EncodeDocument encode interface as bson document
func EncodeDocument(v interface{}) (doc *bson.Document, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}

	err = bson.Unmarshal(data, &doc)
	return doc, err
}
