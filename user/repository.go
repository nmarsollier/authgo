package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/tools/log"
)

var tableName = "users"

type DbUserUpdateDocumentBody struct {
	Name        string `validate:"required,min=1,max=100"`
	Password    string `validate:"required"`
	Permissions []string
	Enabled     bool
	Updated     time.Time
}

func insert(user *User, deps ...interface{}) (_ *User, err error) {
	conn, err := db.GetPostgresClient(deps...)
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	query := `INSERT INTO Users (ID, Name, Login, Password, Permissions, Enabled, Created, Updated) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = conn.Exec(context.TODO(), query, user.ID, user.Name, user.Login, user.Password, user.Permissions, user.Enabled, user.Created, user.Updated)
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	return user, nil
}

func update(user *User, deps ...interface{}) (err error) {
	if err := user.validateSchema(); err != nil {
		log.Get(deps...).Error(err)
		return err
	}

	conn, err := db.GetPostgresClient(deps...)
	if err != nil {
		log.Get(deps...).Error(err)
		return err
	}

	query := `UPDATE Users SET Name=$1, Password=$2, Permissions=$3, Enabled=$4, Updated=$5 WHERE ID=$6`
	_, err = conn.Exec(context.TODO(), query, user.Name, user.Password, user.Permissions, user.Enabled, user.Updated, user.ID)
	if err != nil {
		log.Get(deps...).Error(err)
		return err
	}

	return nil
}

// FindAll devuelve todos los usuarios
func findAll(deps ...interface{}) (users []*User, err error) {
	conn, err := db.GetPostgresClient(deps...)
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	query := `SELECT ID, Name, Login, Password, Permissions, Enabled, Created, Updated FROM Users`
	rows, err := conn.Query(context.TODO(), query)
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.Permissions, &user.Enabled, &user.Created, &user.Updated)
		if err != nil {
			log.Get(deps...).Error(err)
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	return users, nil
}

// FindByID lee un usuario desde la db
func findByID(userID string, deps ...interface{}) (*User, error) {
	conn, err := db.GetPostgresClient(deps...)
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	var user User
	query := `SELECT ID, Name, Login, Password, Permissions, Enabled, Created, Updated FROM Users WHERE ID=$1`
	err = conn.QueryRow(context.TODO(), query, userID).Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.Permissions, &user.Enabled, &user.Created, &user.Updated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NotFound
		}
		log.Get(deps...).Error(err)
		return nil, err
	}

	return &user, err
}

// FindByLogin lee un usuario desde la db
func findByLogin(login string, deps ...interface{}) (*User, error) {
	conn, err := db.GetPostgresClient(deps...)
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	var user User
	query := `SELECT ID, Name, Login, Password, Permissions, Enabled, Created, Updated FROM Users WHERE Login=$1`
	err = conn.QueryRow(context.TODO(), query, login).Scan(&user.ID, &user.Name, &user.Login, &user.Password, &user.Permissions, &user.Enabled, &user.Created, &user.Updated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NotFound
		}
		log.Get(deps...).Error(err)
		return nil, err
	}

	return &user, err
}
