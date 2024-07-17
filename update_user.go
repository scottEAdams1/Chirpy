package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
)

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, 500, "Authorization header not found")
		return
	}

	tokenString := authHeader[len("Bearer "):]
	type MyCustomClaims struct {
		jwt.RegisteredClaims
	}
	claims := &MyCustomClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(*jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
		respondWithError(w, 401, "Unauthorised")
		return
	}
	userID := claims.Subject
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	//Receive the body from the json
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	user, err := cfg.DB.UpdateUser(userIDint, params.Email, []byte(params.Password))
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	type User struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	userResponse := User{
		ID:    user.ID,
		Email: user.Email,
	}
	respondWithJSON(w, 200, userResponse)
}
