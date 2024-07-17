package database

import (
	"errors"
	"sort"
)

// Create a user
func (db *DB) CreateUser(email string, password []byte) (User, error) {
	//Lock database
	db.mux.Lock()
	defer db.mux.Unlock()

	//Load database into struct
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	//Get the id for the user
	var nextID int
	if len(dbStruct.Users) > 0 {
		for id := range dbStruct.Users {
			if id > nextID {
				nextID = id
			}
		}
		nextID++
	} else {
		nextID = 1
	}

	//Create User in form of a struct
	newUser := User{
		ID:       nextID,
		Email:    email,
		Password: password,
	}

	//Check user doesn't already exist
	if dbStruct.Users[newUser.ID].Email != "" {
		return User{}, errors.New("already exists")
	}

	//Add the user to the database struct
	dbStruct.Users[newUser.ID] = newUser

	//Write database struct to the database file as JSON
	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}

// Get users from the database in order of email
func (db *DB) GetUsers() ([]User, error) {
	//Read lock database
	db.mux.RLock()
	defer db.mux.RUnlock()

	//Load database into struct
	dbStruct, err := db.loadDB()
	if err != nil {
		return []User{}, err
	}

	//Make a slice to hold all the users
	users := make([]User, 0, len(dbStruct.Users))
	for _, user := range dbStruct.Users {
		users = append(users, user)
	}

	//Order the users in order of email
	sort.Slice(users, func(i, j int) bool { return users[i].Email < users[j].Email })
	return users, nil
}
