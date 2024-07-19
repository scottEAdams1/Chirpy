package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) chirpyRed(w http.ResponseWriter, r *http.Request) {
	//Receive the body from the json
	type data struct {
		UserID int `json:"user_id"`
	}
	type parameters struct {
		Data  data   `json:"data"`
		Event string `json:"event"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	user, err := cfg.DB.GetUser(params.Data.UserID)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	_, err = cfg.DB.UpdateRed(user)
	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}
	w.WriteHeader(204)
}
