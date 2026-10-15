package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("my super secret key")

func GenerateToken(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user-id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(jwtSecret)
}

func GetJWTSecret() []byte {
	return jwtSecret
}
