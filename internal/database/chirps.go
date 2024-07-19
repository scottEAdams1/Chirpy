package database

import "sort"

// Create a chirp
func (db *DB) CreateChirp(body string, author_id int) (Chirp, error) {
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
		ID:       nextID,
		Body:     body,
		AuthorID: author_id,
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
