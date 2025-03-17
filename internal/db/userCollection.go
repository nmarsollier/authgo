package db

import (
	"github.com/nmarsollier/authgo/internal/common/db"
	"github.com/nmarsollier/authgo/internal/common/log"
	"github.com/nmarsollier/authgo/internal/env"
)

var userCollection db.Collection

func UserCollection() db.Collection {
	if userCollection != nil {
		return userCollection
	}

	log := log.Get(env.Get().FluentURL, env.Get().ServerName)
	userCollection, err := db.NewCollection(log, database(), "users", isDbTimeoutError)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return userCollection
}
