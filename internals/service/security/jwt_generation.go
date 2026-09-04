package security

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserId string
	jwt.RegisteredClaims
}

func GenerateJWtToken(userId string, duration time.Duration, secretKey string) (string, error) {
	log.Println("GenerateJWtToken:", userId)
	claims := CustomClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pollyroute-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", fmt.Errorf("Failed to signed the token : %w", err)
	}

	return signedToken, nil

}

func VerifyJwtToken(tokenString string, secretKey string) (string, error) {
	log.Println("VerifyJwtToken:", tokenString)
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
		log.Println("Verified UserId: ", claims.UserId)
		return claims.UserId, nil
	}

	return "", fmt.Errorf("token is invalid or claims could not be parsed")
}
