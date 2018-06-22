package security

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/nmarsollier/authgo/tools/env"
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
	ID      objectid.ObjectID `bson:"_id"`
	UserID  objectid.ObjectID `bson:"userId"`
	Enabled bool              `bson:"enabled"`
}

// Encode codifica un Token obteniendo el tokenString
func (t Token) Encode() (string, error) {
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
func newToken(userID objectid.ObjectID) *Token {
	return &Token{
		ID:      objectid.New(),
		UserID:  userID,
		Enabled: true,
	}
}
