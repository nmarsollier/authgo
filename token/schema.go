package token

import (
	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/nmarsollier/authgo/tools/env"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/**
 * @apiDefine TokenResponse
 *
 * @apiSuccessExample {json} Respuesta
 *     HTTP/1.1 200 OK
 *     {
 *       "token": "{Token de autorización}"
 *     }
 */

// Token es un objeto valor que representa un token.
type Token struct {
	ID      primitive.ObjectID `bson:"_id"`
	UserID  primitive.ObjectID `bson:"userId"`
	Enabled bool               `bson:"enabled"`
}

// Encode codifica un Token obteniendo el tokenString
func Encode(t *Token) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"tokenID": t.ID.Hex(),
		"userID":  t.UserID.Hex(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := jwtToken.SignedString([]byte(env.Get().JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// NewToken crea un nuevo Token con la información minima necesaria
func newToken(userID primitive.ObjectID) *Token {
	return &Token{
		ID:      primitive.NewObjectID(),
		UserID:  userID,
		Enabled: true,
	}
}
