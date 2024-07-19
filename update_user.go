package main

import (
	"encoding/json"
	"net/http"

	"github.com/scottEAdams1/Chirpy/internal/auth"
)

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	//Get the token string from the header
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	//Get the user ID from the token string
	userID, err := auth.GetUserID(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	//Receive the body from the json
	type parameters struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	var isChirpyRed bool
	if params.IsChirpyRed {
		isChirpyRed = true
	} else {
		isChirpyRed = false
	}

	user, err := cfg.DB.UpdateUser(userID, params.Email, []byte(params.Password), isChirpyRed)
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
