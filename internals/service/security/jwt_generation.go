package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	userId string
	jwt.RegisteredClaims
}

func GenerateJWtToken(userId string, duration time.Duration, secretKey string) (string, error) {

	claims := CustomClaims{
		userId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pollyroute-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)

	if err != nil {
		return "", fmt.Errorf("Failed to signed the token : %w", err)
	}

	return signedToken, nil

}

func VerifyJwtToken(tokenString string, secretKey string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("Invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.userId, nil
	}

	return "", fmt.Errorf("token is invalid or claims could not be parsed")
}
