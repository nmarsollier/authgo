package db

// Decoder interfacea mongo.DocumentResult, para poder testear la app
type Decoder interface {
	Decode(v interface{}) error
}
