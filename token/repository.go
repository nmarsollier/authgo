package token

import (
	"context"

	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/tools/log"
)

// insert crea un nuevo token y lo almacena en la db
func insert(userID string, deps ...interface{}) (token *Token, err error) {
	token = newToken(userID)

	conn, err := db.GetPostgresClient(deps...)
	if err != nil {
		log.Get(deps...).Error(err)

		return nil, err
	}

	query := `INSERT INTO Tokens (ID, UserID, Enabled) VALUES ($1, $2, $3)`
	_, err = conn.Exec(context.TODO(), query, token.ID, token.UserID, token.Enabled)
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	return token, nil
}

// findByID busca un token en la db
func findByID(tokenID string, deps ...interface{}) (*Token, error) {
	conn, err := db.GetPostgresClient(deps...)
	if err != nil {
		log.Get(deps...).Error(err)

		return nil, err
	}

	var token Token
	err = conn.QueryRow(context.Background(), "SELECT id, userId, enabled FROM Tokens WHERE id=$1 and enabled=true", tokenID).Scan(&token.ID, &token.UserID, &token.Enabled)
	if err != nil {
		return nil, errs.NotFound
	}

	return &token, nil
}

// delete como deshabilitado un token
func delete(tokenID string, deps ...interface{}) (err error) {
	conn, err := db.GetPostgresClient(deps...)
	if err != nil {
		log.Get(deps...).Error(err)

		return err
	}

	_, err = conn.Exec(context.Background(), "UPDATE Tokens set enabled=FALSE WHERE id=$1", tokenID)
	if err != nil {
		return errs.NotFound
	}

	return
}
