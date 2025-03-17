package db

import (
	"github.com/nmarsollier/authgo/internal/common/db"
	"github.com/nmarsollier/authgo/internal/common/log"
	"github.com/nmarsollier/authgo/internal/env"
)

var tokenCollection db.Collection

func TokenCollection() db.Collection {
	if tokenCollection != nil {
		return tokenCollection
	}

	log := log.Get(env.Get().FluentURL, env.Get().ServerName)

	tokenCollection, err := db.NewCollection(log, database(), "tokens", isDbTimeoutError, "userId")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return tokenCollection
}
