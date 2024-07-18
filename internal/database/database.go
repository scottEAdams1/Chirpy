package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp    `json:"chrips"`
	Users  map[int]User     `json:"users"`
	Tokens map[string]Token `json:"tokens"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

type Token struct {
	Token      string    `json:"token"`
	UserID     int       `json:"user_id"`
	Expiration time.Time `json:"expiration"`
}

// Create a database
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	//Check if the database exists
	err := db.ensureDB()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Check if the database exists
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		dbStruct := DBStructure{
			Chirps: make(map[int]Chirp),
			Users:  make(map[int]User),
			Tokens: make(map[string]Token),
		}
		return db.writeDB(dbStruct)
	}
	return err
}

// Load the database into a struct
func (db *DB) loadDB() (DBStructure, error) {
	//Read the file holding the JSON database
	structure, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	//Make database struct to hold the chirps
	dbStructure := DBStructure{
		Chirps: make(map[int]Chirp),
		Users:  make(map[int]User),
		Tokens: make(map[string]Token),
	}

	//Convert the file data into the struct
	err = json.Unmarshal(structure, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}
	return dbStructure, nil
}

// Write the database struct to the file
func (db *DB) writeDB(dbStructure DBStructure) error {
	//Convert database struct to JSON
	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	//Write the JSON to the database JSON file
	err = os.WriteFile(db.path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
