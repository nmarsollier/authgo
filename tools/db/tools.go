package db

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

// IsDbTimeoutError función a llamar cuando se produce un error de db
func IsDbTimeoutError(err interface{}) {
	if err == topology.ErrServerSelectionTimeout {
		database = nil
	}
}

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
