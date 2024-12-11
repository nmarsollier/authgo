package token

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/nmarsollier/authgo/tools/env"
	"github.com/nmarsollier/authgo/tools/errs"
	uuid "github.com/satori/go.uuid"
)

// Token es una estructura valor que representa un token.
type Token struct {
	ID      string
	UserID  string
	Enabled bool
}

// NewToken crea un nuevo Token con la información minima necesaria
func newToken(userID string) *Token {
	return &Token{
		ID:      uuid.NewV4().String(),
		UserID:  userID,
		Enabled: true,
	}
}

// Encode codifica un Token obteniendo el tokenString
func Encode(t *Token) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"tokenID": t.ID,
		"userID":  t.UserID,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := jwtToken.SignedString([]byte(env.Get().JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
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
