package token

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nmarsollier/authgo/tools/env"
	"github.com/nmarsollier/authgo/tools/errors"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rabbit"
	gocache "github.com/patrickmn/go-cache"
)

var cache = gocache.New(60*time.Minute, 10*time.Minute)

var jwtSecret = []byte(env.Get().JWTSecret)

// Payload es la información cifrada que se guarda en el token
type Payload struct {
	TokenID string
	UserID  string
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
func Create(userID string) (string, error) {
	token := newToken()
	token.UserID = userID

	token, err := save(token)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"tokenID": token.ID(),
		"userID":  token.UserID,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := jwtToken.SignedString(jwtSecret)

	return tokenString, err
}

// Validate valida un token
/**
 * @apiDefine AuthHeader
 *
 * @apiParamExample {String} Header Autorización
 *    Authorization=bearer {token}
 *
 * @apiSuccessExample 401 Unauthorized
 *    HTTP/1.1 401 Unauthorized
 */
func Validate(c *gin.Context) (*Payload, error) {
	tokenString, err := getTokenHeader(c)
	if err != nil {
		return nil, err
	}

	// Si esta en cache, retornamos el cache
	if found, ok := cache.Get(tokenString); ok {
		if payload, ok := found.(Payload); ok {
			return &payload, nil
		}
	}

	// Sino validamos el token y lo agregamos al cache
	payload, err := extractPayload(tokenString)
	if err != nil {
		return nil, err
	}

	// Buscamos el token en la db para validarlo
	dbToken, err := findByID(payload.TokenID)
	if err != nil || !dbToken.Enabled {
		return nil, errors.Unauthorized
	}

	// Todo bien, se agrega al cache y se retorna
	cache.Set(tokenString, payload, gocache.DefaultExpiration)

	return payload, nil
}

// Invalidate invalida un token
func Invalidate(c *gin.Context) error {
	payload, err := Validate(c)
	if err != nil {
		return errors.Unauthorized
	}

	if err = delete(payload.TokenID); err != nil {
		return err
	}

	go func() {
		tokenString := c.GetHeader("Authorization")

		if err = rabbit.SendLogout(tokenString); err != nil {
			log.Output(1, "Rabbit logout no se pudo enviar")
		}

		cache.Delete(tokenString[7:])
	}()

	return nil
}

// get token from Authorization header
func getTokenHeader(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")
	if strings.Index(tokenString, "bearer ") != 0 {
		return "", errors.Unauthorized
	}
	tokenString = tokenString[7:]
	return tokenString, nil
}

// extract payload from token string
func extractPayload(tokenString string) (*Payload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecret, nil
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
