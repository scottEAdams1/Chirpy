package database

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chrips"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
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

// Create a chirp
func (db *DB) CreateChirp(body string) (Chirp, error) {
	//Lock database
	db.mux.Lock()
	defer db.mux.Unlock()

	//Load database into struct
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	//Get the id for the chirp
	var nextID int
	if len(dbStruct.Chirps) > 0 {
		for id := range dbStruct.Chirps {
			if id > nextID {
				nextID = id
			}
		}
		nextID++
	} else {
		nextID = 1
	}

	//Create Chirp in form of a struct
	newChirp := Chirp{
		ID:   nextID,
		Body: body,
	}

	//Add the chirp to the database struct
	dbStruct.Chirps[newChirp.ID] = newChirp

	//Write database struct to the database file as JSON
	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, err
	}
	return newChirp, nil
}

// Get chirps from the database in order of id
func (db *DB) GetChirps() ([]Chirp, error) {
	//Read lock database
	db.mux.RLock()
	defer db.mux.RUnlock()

	//Load database into struct
	dbStruct, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	//Make a slice to hold all the chirps
	chirps := make([]Chirp, 0, len(dbStruct.Chirps))
	for _, chirp := range dbStruct.Chirps {
		chirps = append(chirps, chirp)
	}

	//Order the chirps in order of id
	sort.Slice(chirps, func(i, j int) bool { return chirps[i].ID < chirps[j].ID })
	return chirps, nil
}

// Check if the database exists
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		dbStruct := DBStructure{
			Chirps: make(map[int]Chirp),
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
