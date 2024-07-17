package main

import (
	"encoding/json"
	"net/http"

	"github.com/scottEAdams1/Chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
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
