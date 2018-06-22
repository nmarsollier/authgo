package security

import (
	"fmt"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/nmarsollier/authgo/rabbit"
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
	dao, err := newDao()
	if err != nil {
		return nil, err
	}

	return serviceStruct{
		dao: dao,
	}, nil
}

// MockedService con fines de testing para mockear db.collection
func MockedService(coll Dao) Service {
	return serviceStruct{
		dao: coll,
	}
}

// El repositorio
type serviceStruct struct {
	dao Dao
}

// Create crea un nuevo token y lo almacena en la db
func (d serviceStruct) Create(userID objectid.ObjectID) (*Token, error) {
	token, err := d.dao.Create(userID)
	if err != nil {
		return nil, err
	}

	cacheAdd(token)

	return token, nil
}

// Find busca un token en la db
func (d serviceStruct) Find(tokenID string) (*Token, error) {
	return d.dao.FindByID(tokenID)
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
	token, err := d.dao.FindByID(tokenID)
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

	if err = d.dao.Delete(token.ID); err != nil {
		return err
	}

	go func() {
		cacheRemove(token)

		if err = rabbit.SendLogout("bearer " + tokenString); err != nil {
			log.Output(1, "Rabbit logout no se pudo enviar")
		}
	}()

	return nil
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
