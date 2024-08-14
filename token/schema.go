package token

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/nmarsollier/authgo/tools/env"
	"github.com/nmarsollier/authgo/tools/errs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

// descifra el token string y devuelve los datos del payload
func ExtractPayload(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(env.Get().JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return "", "", errs.Unauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", "", errs.Unauthorized
	}

	tokenID := claims["tokenID"].(string)
	userID := claims["userID"].(string)

	return tokenID, userID, nil
}
