package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// IsDbUniqueKeyError retorna true si el error es de indice único
func IsDbUniqueKeyError(err error) bool {
	if wErr, ok := err.(mongo.WriteException); ok {
		for i := 0; i < len(wErr.WriteErrors); i++ {
			if wErr.WriteErrors[i].Code == 11000 {
				return true
			}
		}
	}
	return false
}
