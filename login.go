package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/scottEAdams1/Chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	//Receive the body from the json
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	//Get all users from the database
	users, err := cfg.DB.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//Get the user with the correct email
	var user database.User
	for _, user1 := range users {
		if user1.Email == params.Email {
			user = user1
		}
	}
	if user.Email == "" {
		respondWithError(w, 400, "Email doesn't exist")
		return
	}

	//Check if the password is correct
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(params.Password))
	if err != nil {
		respondWithError(w, 401, "Unauthorised")
		return
	}

	var expirationTime time.Time
	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > 60*60*24 {
		expirationTime = time.Now().UTC().Add(24 * time.Hour)
	} else {
		expirationTime = time.Now().UTC().Add(time.Duration(params.ExpiresInSeconds) * time.Second)
	}

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Subject:   strconv.Itoa(user.ID),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	type User struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}

	userResponse := User{
		ID:    user.ID,
		Email: user.Email,
		Token: signedToken,
	}
	respondWithJSON(w, 200, userResponse)
}
