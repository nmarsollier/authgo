package user

import (
	"context"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errors"
)

type daoImpl struct {
	dbCollection db.Collection
}

type Dao interface {
	Collection() (db.Collection, error)
	Insert(user *User) (*User, error)
	Update(user *User) (*User, error)
	FindAll() ([]*User, error)
	FindByID(userID string) (*User, error)
	FindByLogin(login string) (*User, error)
	Delete(userID string) error
	GetID(ID string) (*objectid.ObjectID, error)
}

func newDao() Dao {
	return daoImpl{}
}

func NewTestingDao(coll db.Collection) Dao {
	return daoImpl{
		dbCollection: coll,
	}
}

// UsersCollection obtiene la colección de Usuarios
func (d daoImpl) Collection() (db.Collection, error) {
	if d.dbCollection != nil {
		return d.dbCollection, nil
	}

	database, err := db.Get()
	if err != nil {
		return nil, err
	}

	collection := database.Collection("users")

	_, err = collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.NewDocument(
				bson.EC.String("login", ""),
			),
			Options: bson.NewDocument(
				bson.EC.Boolean("unique", true),
			),
		},
	)
	if err != nil {
		log.Output(1, err.Error())
	}

	d.dbCollection = db.WrapCollection(collection)
	return d.dbCollection, nil
}

func (d daoImpl) Insert(user *User) (*User, error) {
	if err := user.ValidateSchema(); err != nil {
		return nil, err
	}

	collection, err := d.Collection()
	if err != nil {
		return nil, err
	}

	if _, err = collection.InsertOne(context.Background(), user); err != nil {
		return nil, err
	}

	return user, nil
}

func (d daoImpl) Update(user *User) (*User, error) {
	if err := user.ValidateSchema(); err != nil {
		return nil, err
	}

	collection, err := d.Collection()
	if err != nil {
		return nil, err
	}

	user.Updated = time.Now()

	doc, err := bson.NewDocumentEncoder().EncodeDocument(user)
	if err != nil {
		return nil, err
	}

	_, err = collection.UpdateOne(context.Background(),
		bson.NewDocument(doc.LookupElement("_id")),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				doc.LookupElement("password"),
				doc.LookupElement("name"),
				doc.LookupElement("enabled"),
				doc.LookupElement("updated"),
				doc.LookupElement("permissions"),
			),
		))

	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindAll devuelve todos los usuarios
func (d daoImpl) FindAll() ([]*User, error) {
	collection, err := d.Collection()
	if err != nil {
		return nil, err
	}

	users := []*User{}
	filter := bson.NewDocument()
	cur, err := collection.Find(context.Background(), filter, nil)
	defer cur.Close(context.Background())

	if err != nil {
		return nil, err
	}

	for cur.Next(context.Background()) {
		user := &User{}
		if err := cur.Decode(user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// FindByID lee un usuario desde la db
func (d daoImpl) FindByID(userID string) (*User, error) {
	_id, err := objectid.FromHex(userID)
	if err != nil {
		return nil, errors.ErrID
	}

	collection, err := d.Collection()
	if err != nil {
		return nil, err
	}

	user := &User{}
	filter := bson.NewDocument(bson.EC.ObjectID("_id", _id))
	if err = collection.FindOne(context.Background(), filter).Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

// FindByLogin lee un usuario desde la db
func (d daoImpl) FindByLogin(login string) (*User, error) {
	collection, collectionError := d.Collection()
	if collectionError != nil {
		return nil, collectionError
	}

	user := &User{}
	filter := bson.NewDocument(bson.EC.String("login", login))
	err := collection.FindOne(context.Background(), filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrLogin
		}
		return nil, err
	}

	return user, nil
}

// Delete marca un usuario como borrado en la base de datos
func (d daoImpl) Delete(userID string) error {
	_id, err := d.GetID(userID)
	if err != nil {
		return err
	}

	collection, err := d.Collection()
	if err != nil {
		return err
	}

	user := NewUser()
	user.ID = *_id
	user.Enabled = false
	user.Updated = time.Now()

	doc, err := bson.NewDocumentEncoder().EncodeDocument(user)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(context.Background(),
		bson.NewDocument(doc.LookupElement("_id")),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				doc.LookupElement("enabled"),
				doc.LookupElement("updated"),
			),
		))

	return err
}

func (d daoImpl) GetID(ID string) (*objectid.ObjectID, error) {
	_id, err := objectid.FromHex(ID)
	if err != nil {
		return nil, errors.ErrID
	}
	return &_id, nil
}
