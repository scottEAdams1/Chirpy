package database

import (
	"errors"
	"sort"

	"golang.org/x/crypto/bcrypt"
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
		ID:          nextID,
		Email:       email,
		Password:    password,
		IsChirpyRed: false,
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

// Get a single user based on the id
func (db *DB) GetUser(id int) (User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	users, err := db.GetUsers()
	if err != nil {
		return User{}, err
	}

	var returnUser User
	for _, user := range users {
		if user.ID == id {
			returnUser = user
		}
	}

	if returnUser.Email == "" {
		return User{}, errors.New("User doesn't exist")
	}
	return returnUser, nil
}

// Updates a user with new information
func (db *DB) UpdateUser(id int, email string, password []byte, isChirpyRed bool) (User, error) {
	//Lock database
	db.mux.Lock()
	defer db.mux.Unlock()

	//Load database into struct
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(password, 1)
	if err != nil {
		return User{}, err
	}

	//Create User in form of a struct
	newUser := User{
		ID:          id,
		Email:       email,
		Password:    hashedPassword,
		IsChirpyRed: isChirpyRed,
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

func (db *DB) UpdateRed(user User) (User, error) {
	//Lock database
	db.mux.Lock()
	defer db.mux.Unlock()

	//Load database into struct
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	//Create User in form of a struct
	newUser := User{
		ID:          user.ID,
		Email:       user.Email,
		Password:    user.Password,
		IsChirpyRed: true,
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
