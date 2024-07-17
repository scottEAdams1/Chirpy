package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	//Receive the body from the json
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	//Hash password
	bytePassword := []byte(params.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, 1)
	if err != nil {
		respondWithError(w, 400, err.Error())
	}

	//Create a user from the email
	newUser, err := cfg.DB.CreateUser(params.Email, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//Create a user without password
	type User struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	userWOPassword := User{
		ID:    newUser.ID,
		Email: newUser.Email,
	}
	//Send response
	respondWithJSON(w, http.StatusCreated, userWOPassword)
}
