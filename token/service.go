package token

import (
	"fmt"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"

	"github.com/nmarsollier/authgo/tools/env"
	"github.com/nmarsollier/authgo/tools/errors"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/nmarsollier/authgo/rabbit"
	gocache "github.com/patrickmn/go-cache"
)

var cache = gocache.New(60*time.Minute, 10*time.Minute)

// Payload es la información cifrada que se guarda en el token
type Payload struct {
	TokenID string
	UserID  string
}

type serviceImpl struct {
	tokenDao dao
}

// Service es la interfaz ue define el servicio
type Service interface {
	Create(userID objectid.ObjectID) (string, error)
	Validate(token string) (*Payload, error)
	Invalidate(token string) error
}

// NewService retorna una nueva instancia del servicio
func NewService() Service {
	return serviceImpl{
		tokenDao: newDao(),
	}
}

// NewTestingService retorna un servicio con fines de test
func NewTestingService(fakeDao dao) Service {
	return serviceImpl{
		tokenDao: fakeDao,
	}
}

// Create crea un token
/**
 * @apiDefine TokenResponse
 *
 * @apiSuccessExample {json} Respuesta
 *     HTTP/1.1 200 OK
 *     {
 *       "token": "{Token de autorización}"
 *     }
 */
func (s serviceImpl) Create(userID objectid.ObjectID) (string, error) {
	token := newToken()
	token.UserID = userID

	token, err := s.tokenDao.insert(token)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"tokenID": token.ID.Hex(),
		"userID":  token.UserID.Hex(),
	})

	// Sign and get the complete encoded token as a string using the secret
	return jwtToken.SignedString([]byte(env.Get().JWTSecret))
}

// Validate valida un token
/**
 * @apiDefine AuthHeader
 *
 * @apiExample {String} Header Autorización
 *    Authorization=bearer {token}
 *
 * @apiErrorExample 401 Unauthorized
 *     HTTP/1.1 401 Unauthorized
 *     {
 *        "error" : "Unauthorized"
 *     }
 */
func (s serviceImpl) Validate(token string) (*Payload, error) {
	// Si esta en cache, retornamos el cache
	if found, ok := cache.Get(token); ok {
		if payload, ok := found.(Payload); ok {
			return &payload, nil
		}
	}

	// Sino validamos el token y lo agregamos al cache
	payload, err := extractPayload(token)
	if err != nil {
		return nil, err
	}

	// Buscamos el token en la db para validarlo
	dbToken, err := s.tokenDao.findByID(payload.TokenID)
	if err != nil || !dbToken.Enabled {
		return nil, errors.Unauthorized
	}

	// Todo bien, se agrega al cache y se retorna
	cache.Set(token, payload, gocache.DefaultExpiration)

	return payload, nil
}

// Invalidate invalida un token
func (s serviceImpl) Invalidate(token string) error {
	payload, err := s.Validate(token)
	if err != nil {
		return errors.Unauthorized
	}

	if err = s.tokenDao.delete(payload.TokenID); err != nil {
		return err
	}

	go func() {
		if err = rabbit.SendLogout("bearer " + token); err != nil {
			log.Output(1, "Rabbit logout no se pudo enviar")
		}

		cache.Delete(token)
	}()

	return nil
}

// extract payload from token string
func extractPayload(tokenString string) (*Payload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(env.Get().JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.Unauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.Unauthorized
	}

	payload := &Payload{
		UserID:  claims["userID"].(string),
		TokenID: claims["tokenID"].(string),
	}

	return payload, nil
}
