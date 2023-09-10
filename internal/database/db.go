package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	path   string
	mux    *sync.RWMutex
	id     int
	userId int
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := DB{
		path:   path,
		mux:    &sync.RWMutex{},
		id:     0,
		userId: 0,
	}
	err := db.ensureDB()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	fmt.Println("Creating chirp method")
	chirp := Chirp{}
	dat, err := db.loadDB()
	if err != nil {
		return chirp, err
	}
	chirp.Chirp = body
	db.id++
	chirp.ID = db.id

	dat.Chirps[db.id] = chirp
	fmt.Println("Before actual db write")
	db.writeDB(dat)
	fmt.Printf("AFter actual db write %v", chirp)

	return chirp, nil
}

// CreateChirp creates a new user and saves it to disk
func (db *DB) CreateUser(body string, password string) (User, error) {
	fmt.Println("Creating chirp method")
	user := User{}
	dat, err := db.loadDB()
	if err != nil {
		return user, err
	}
	user.Email = body
	db.id++
	user.ID = db.id

	dat.Users[db.id] = user
	fmt.Println("Before actual db write")
	db.writeDB(dat)
	fmt.Printf("AFter actual db write %v", user)

	return user, nil
}

// CreateChirp creates a new user with password and saves it to disk
func (db *DB) CreateUserWithPassword(email string, password string) (User, error) {
	fmt.Println("Creating chirp method")

	user := User{}
	dat, err := db.loadDB()

	if err != nil {
		return user, err
	}

	// existingUser, ok := dat.Users[db.id]
	// if ok {
	// 	fmt.Println("User already exists")
	// 	return existingUser, errors.New("User already exists")
	// }

	hashPass, err := hashPassword(password)
	if err != nil {
		return user, err
	}

	user.Email = email
	user.Password = hashPass
	db.id++
	user.ID = db.id

	dat.Users[db.id] = user
	fmt.Println("Before actual db write")
	db.writeDB(dat)
	fmt.Printf("AFter actual db write %v", user)

	return user, nil
}

func hashPassword(password string) (string, error) {
	// Generate a bcrypt hash of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error generating bcrypt hash:", err)
		return "", err
	}

	// Store the hashed password in your database or wherever you need to
	return string(hashedPassword), nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	chirps := []Chirp{}
	dat, err := db.loadDB()
	if err != nil {
		return chirps, err
	}

	for _, v := range dat.Chirps {
		chirps = append(chirps, v)
	}

	return chirps, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetUsers() ([]User, error) {
	users := []User{}
	dat, err := db.loadDB()
	if err != nil {
		return users, err
	}

	for _, v := range dat.Users {
		users = append(users, v)
	}

	return users, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	user := User{}

	users, err := db.GetUsers()

	if err != nil {
		return user, err
	}

	for _, usr := range users {
		if usr.Email == email {
			user = usr
			break
		}
	}

	return user, nil
}

func (db *DB) UpdateUser(id int, newEmail string, newPassword string) (bool, User, error) {
	data, err := db.loadDB()
	if err != nil {
		return false, User{}, err
	}

	user, ok := data.Users[id]
	if !ok {
		return false, User{}, fmt.Errorf("user not found")
	}

	user.Email = newEmail

	hashPass, err := hashPassword(newPassword)
	if err != nil {
		return false, user, err
	}

	fmt.Printf("old password %v\n", user.Password)

	user.Password = hashPass

	fmt.Printf("new password %v\n", user.Password)
	data.Users[id] = user

	fmt.Printf("updated db %v\n", data.Users)
	fmt.Printf("updated db %v\n", data)

	err = db.writeDB(data)

	if err != nil {
		return false, user, err
	}

	return true, user, nil
}

func (db *DB) ValidatePasswordForUser(user User, password string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return false, err
	}

	return true, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	f, err := os.OpenFile(db.path, os.O_RDONLY|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	defer f.Close()

	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	fmt.Printf("Calling load db")
	err := db.ensureDB()
	dbStruct := DBStructure{
		Chirps: make(map[int]Chirp),
		Users:  make(map[int]User),
	}
	fmt.Println("dbStructure made")
	if err != nil {
		return dbStruct, err
	}
	fmt.Println("Trying to read file")
	dat, err := os.ReadFile(db.path)
	if err != nil {
		return dbStruct, err
	}

	if len(dat) > 0 {
		err = json.Unmarshal(dat, &dbStruct)
		if err != nil {
			return dbStruct, err
		}
	}

	fmt.Println(dbStruct)

	return dbStruct, nil

}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	fmt.Printf("trying to write data %v\n", dbStructure)
	err := db.ensureDB()
	if err != nil {
		return err
	}

	dat, err := json.Marshal(dbStructure)

	if err != nil {
		return err
	}

	err = os.WriteFile("database.json", dat, os.ModeAppend)
	if err != nil {
		return err
	}

	return nil
}
