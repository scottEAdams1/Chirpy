package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/scottEAdams1/Chirpy/internal/auth"
)

// Create a chirp and add it to the database
func (cfg *apiConfig) createChirps(w http.ResponseWriter, r *http.Request) {
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
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	//If the message if more than 140 characters, error
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	//Clean the text of the words "kerfuffle", "sharbert", "fornax"
	cleaned := getCleanedBody(params.Body)

	//Create a chirp from the cleaned text
	newChirp, err := cfg.DB.CreateChirp(cleaned, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//Send response
	respondWithJSON(w, http.StatusCreated, newChirp)
}

// Clean the text of the words "kerfuffle", "sharbert", "fornax"
func getCleanedBody(body string) string {
	//Split text into individual words
	words := strings.Split(body, " ")

	cleaned_words := []string{}

	//For each word, check if it is one to clean, then add it to the cleaned_words slice
	for _, word := range words {
		lowerword := strings.ToLower(word)
		if lowerword == "kerfuffle" || lowerword == "sharbert" || lowerword == "fornax" {
			word = "****"
		}
		cleaned_words = append(cleaned_words, word)
	}

	//Convert the slice to a string
	cleaned := strings.Join(cleaned_words, " ")
	return cleaned
}

// Respond to request with an error
func respondWithError(w http.ResponseWriter, code int, msg string) {
	//Server error
	if code > 499 {
		log.Printf("Error decoding parameters: %s", msg)
	}

	//Return an error struct with response
	type errorStruct struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errorStruct{
		Error: msg,
	})
}

// Respond in JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	//Set response to JSON
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	//Convert struct into JSON
	jsonResponse, jsonErr := json.Marshal(payload)
	if jsonErr != nil {
		log.Printf("Error marshalling JSON: %s", jsonErr)
		w.WriteHeader(500)
		return
	}

	//Add code(e.g. 200, 500) to header
	w.WriteHeader(code)

	//Add JSON to the body of the response
	w.Write(jsonResponse)
}
