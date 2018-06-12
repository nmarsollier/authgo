package lookup

import (
	"github.com/mongodb/mongo-go-driver/bson"
)

// StringDefault will return key value or default value if missing.
func StringDefault(d bson.Document, key string, def string) string {
	elem, err := d.LookupErr(key)
	if err != nil {
		return def
	}

	return elem.StringValue()
}

// String will return key value or default value if missing.
func String(d bson.Document, key string) string {
	return StringDefault(d, key, "")
}

// BoolDefault will return key value or default value if missing.
func BoolDefault(d bson.Document, key string, def bool) bool {
	elem, err := d.LookupErr(key)
	if err != nil {
		return def
	}

	return elem.Boolean()
}

// Bool will return key value or default value if missing.
func Bool(d bson.Document, key string) bool {
	return BoolDefault(d, key, false)
}

// StringArrayDefault will return key value or default value if missing.
func StringArrayDefault(d bson.Document, key string, def []string) []string {
	elem, err := d.LookupErr(key)

	if err != nil {
		return def
	}

	array, ok := elem.MutableArrayOK()
	if !ok {
		return def
	}

	result := make([]string, array.Len())
	for i := 0; i < array.Len(); i++ {
		rol, _ := array.Lookup(uint(i))
		result[i] = rol.StringValue()
	}

	return result
}

// StringArray will return key value or default value if missing.
func StringArray(d bson.Document, key string) []string {
	return StringArrayDefault(d, key, []string{})
}

// ObjectIDDefault will return key value or default value if missing.
func ObjectIDDefault(d bson.Document, key string, def string) string {
	elem, err := d.LookupErr(key)
	if err != nil {
		return def
	}

	return elem.ObjectID().Hex()
}

// ObjectID will return key value or default value if missing.
func ObjectID(d bson.Document, key string) string {
	return ObjectIDDefault(d, key, "")
}
