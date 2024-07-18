package database

import (
	"errors"
	"sort"
	"time"
)

// Create a token
func (db *DB) CreateToken(token string, user_id int) (Token, error) {
	//Lock database
	db.mux.Lock()
	defer db.mux.Unlock()

	//Load database into struct
	dbStruct, err := db.loadDB()
	if err != nil {
		return Token{}, err
	}

	//Get the expiration time
	expiration := time.Now().UTC().AddDate(0, 0, 60)

	//Create Token in form of a struct
	newToken := Token{
		Token:      token,
		UserID:     user_id,
		Expiration: expiration,
	}

	//Add the token to the database struct
	dbStruct.Tokens[newToken.Token] = newToken

	//Write database struct to the database file as JSON
	err = db.writeDB(dbStruct)
	if err != nil {
		return Token{}, err
	}
	return newToken, nil
}

// Get tokens from the database in order of id
func (db *DB) GetTokens() ([]Token, error) {
	//Read lock database
	db.mux.RLock()
	defer db.mux.RUnlock()

	//Load database into struct
	dbStruct, err := db.loadDB()
	if err != nil {
		return []Token{}, err
	}

	//Make a slice to hold all the tokens
	tokens := make([]Token, 0, len(dbStruct.Tokens))
	for _, token := range dbStruct.Tokens {
		tokens = append(tokens, token)
	}

	//Order the users in order of email
	sort.Slice(tokens, func(i, j int) bool { return tokens[i].UserID < tokens[j].UserID })
	return tokens, nil
}

// Get a single token based on the token string
func (db *DB) GetToken(refreshToken string) (Token, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	tokens, err := db.GetTokens()
	if err != nil {
		return Token{}, err
	}

	var returnToken Token
	for _, token := range tokens {
		if token.Token == refreshToken {
			returnToken = token
		}
	}

	if returnToken.Token == "" {
		return Token{}, errors.New("Token doesn't exist")
	}
	return returnToken, nil
}

// Remove a token from the database
func (db *DB) RemoveToken(tokenString string) error {
	//Lock database
	db.mux.Lock()
	defer db.mux.Unlock()

	//Load database into struct
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}

	//Remove token from database struct
	delete(dbStruct.Tokens, tokenString)

	//Write database struct to the database file as JSON
	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}
	return nil
}
