package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type Claims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

func GenerateJWT(userId string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	MY_SECRET_KEY, ok := os.LookupEnv("MY_SECRET_KEY")
	if !ok {
		panic("secret key dont found")
	}

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte(MY_SECRET_KEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		MY_SECRET_KEY, ok := os.LookupEnv("MY_SECRET_KEY")
		if !ok {
			return nil, fmt.Errorf("secret key not found")
		}
		return []byte(MY_SECRET_KEY), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	return claims, nil
}
