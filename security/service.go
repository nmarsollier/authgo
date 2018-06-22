package security

import (
	"context"
	"fmt"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/nmarsollier/authgo/rabbit"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/env"
	"github.com/nmarsollier/authgo/tools/errors"
)

// Service es la interfaz con los método expuestos por este dao
type Service interface {
	Create(userID objectid.ObjectID) (*Token, error)
	Find(tokenID string) (*Token, error)
	Validate(tokenString string) (*Token, error)
	Invalidate(tokenString string) error
}

// NewService devuelve el servicio principal de seguridad
func NewService() (Service, error) {
	database, err := db.Get()
	if err != nil {
		return nil, err
	}

	collection := database.Collection("tokens")

	_, err = collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.NewDocument(
				bson.EC.String("userId", ""),
			),
			Options: bson.NewDocument(),
		},
	)
	if err != nil {
		return nil, err
	}

	return serviceStruct{
		collection: db.WrapCollection(collection),
	}, nil
}

// MockedService con fines de testing para mockear db.collection
func MockedService(coll db.Collection) Service {
	return serviceStruct{
		collection: coll,
	}
}

// El repositorio
type serviceStruct struct {
	collection db.Collection
}

// Create crea un nuevo token y lo almacena en la db
func (d serviceStruct) Create(userID objectid.ObjectID) (*Token, error) {
	token := newToken(userID)

	_, err := d.collection.InsertOne(context.Background(), token)
	if err != nil {
		return nil, err
	}

	cacheAdd(token)

	return token, nil
}

// Find busca un token en la db
func (d serviceStruct) Find(tokenID string) (*Token, error) {
	_id, err := objectid.FromHex(tokenID)
	if err != nil {
		return nil, errors.Unauthorized
	}

	token := &Token{}
	filter := bson.NewDocument(bson.EC.ObjectID("_id", _id))

	if err = d.collection.FindOne(context.Background(), filter).Decode(token); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Unauthorized
		}
		return nil, err
	}

	return token, nil
}

// Validate dado un tokenString devuelve el Token asociado
func (d serviceStruct) Validate(tokenString string) (*Token, error) {
	if token, err := cacheGet(tokenString); err == nil {
		return token, err
	}

	// Sino validamos el token y lo agregamos al cache
	tokenID, _, err := extractPayload(tokenString)
	if err != nil {
		return nil, err
	}

	// Buscamos el token en la db para validarlo
	token, err := d.Find(tokenID)
	if err != nil || !token.Enabled {
		return nil, errors.Unauthorized
	}

	// Todo bien, se agrega al cache y se retorna
	cacheAdd(token)

	return token, nil
}

// Invalidate invalida un token
func (d serviceStruct) Invalidate(tokenString string) error {
	token, err := d.Validate(tokenString)
	if err != nil {
		return errors.Unauthorized
	}

	if err = d.disable(token); err != nil {
		return err
	}

	go func() {
		if err = rabbit.SendLogout("bearer " + tokenString); err != nil {
			log.Output(1, "Rabbit logout no se pudo enviar")
		}

		cacheRemove(token)
	}()

	return nil
}

// marca como deshabilitado un token
func (d serviceStruct) disable(token *Token) error {
	token.Enabled = false
	doc, err := bson.NewDocumentEncoder().EncodeDocument(token)
	if err != nil {
		return err
	}

	_, err = d.collection.UpdateOne(context.Background(),
		bson.NewDocument(bson.EC.ObjectID("_id", token.ID)),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				doc.LookupElement("enabled"),
			),
		))

	return err
}

// descifra el token string y devuelve los datos del payload
func extractPayload(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(env.Get().JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.Unauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", "", errors.Unauthorized
	}

	tokenID := claims["tokenID"].(string)
	userID := claims["userID"].(string)

	return tokenID, userID, nil
}
