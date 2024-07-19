package auth

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}

	tokenString := authHeader[len("Bearer "):]
	if len(tokenString) == 0 {
		return "", errors.New("no token")
	}
	return tokenString, nil

}

func GetUserID(tokenString, jwtSecret string) (int, error) {
	type MyCustomClaims struct {
		jwt.RegisteredClaims
	}
	claims := &MyCustomClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(*jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return 0, errors.New("unauthorised")
	}
	userID := claims.Subject
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		return 0, err
	}
	return userIDint, nil
}
