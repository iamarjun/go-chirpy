package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	path   string
	mux    *sync.RWMutex
	id     int
	userId int
}

type DBStructure struct {
	Chirps        map[int]Chirp    `json:"chirps"`
	Users         map[int]User     `json:"users"`
	RefreshTokens map[string]Token `json:"refresh_tokens"`
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
func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {
	chirp := Chirp{}
	dat, err := db.loadDB()
	if err != nil {
		return chirp, err
	}

	chirp.Chirp = body
	chirp.AuthorID = authorID
	db.id++
	chirp.ID = db.id

	dat.Chirps[db.id] = chirp
	db.writeDB(dat)

	return chirp, nil
}

// CreateChirp creates a new user and saves it to disk
func (db *DB) CreateUser(body string, password string) (User, error) {
	user := User{}
	dat, err := db.loadDB()
	if err != nil {
		return user, err
	}
	user.Email = body
	db.userId++
	user.ID = db.userId

	dat.Users[db.userId] = user
	db.writeDB(dat)

	return user, nil
}

// CreateChirp creates a new user with password and saves it to disk
func (db *DB) CreateUserWithPassword(email string, password string) (User, error) {
	user := User{}
	dat, err := db.loadDB()

	if err != nil {
		return user, err
	}

	hashPass, err := hashPassword(password)
	if err != nil {
		return user, err
	}

	user.Email = email
	user.Password = hashPass
	db.userId++
	user.ID = db.userId

	dat.Users[db.userId] = user
	db.writeDB(dat)

	return user, nil
}

func hashPassword(password string) (string, error) {
	// Generate a bcrypt hash of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Store the hashed password in your database or wherever you need to
	return string(hashedPassword), nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirpsByAuthorId(authorID int) ([]Chirp, error) {
	chirps := []Chirp{}
	dat, err := db.loadDB()
	if err != nil {
		return chirps, err
	}

	for _, chirp := range dat.Chirps {
		if chirp.AuthorID == authorID {
			chirps = append(chirps, chirp)
		}
	}

	return chirps, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	chirps := []Chirp{}
	dat, err := db.loadDB()
	if err != nil {
		return chirps, err
	}

	for _, chirp := range dat.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirpsById(chirpID int, authorID int) (Chirp, error) {
	chirps, err := db.GetChirpsByAuthorId(authorID)
	if err != nil {
		return Chirp{}, err
	}

	for _, chirp := range chirps {
		if chirp.ID == chirpID {
			return chirp, nil
		}
	}

	return Chirp{}, fmt.Errorf("chirp not found")
}

// GetChirps returns all chirps in the database
func (db *DB) DeleteChirpsById(chirpID int, authorID int) (bool, error) {
	data, err := db.loadDB()
	if err != nil {
		return false, err
	}

	chirp, ok := data.Chirps[chirpID]

	if !ok {
		return false, fmt.Errorf("chirp not found")
	}

	if chirp.AuthorID != authorID {
		return false, fmt.Errorf("not authorized to delete this chirp")
	}

	delete(data.Chirps, chirpID)

	err = db.writeDB(data)

	if err != nil {
		return false, err
	}

	return true, nil
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
			return user, nil
		}
	}

	return user, fmt.Errorf("user not found")
}

func (db *DB) GetUserById(userId int) (User, error) {
	user := User{}

	users, err := db.GetUsers()

	if err != nil {
		return user, err
	}

	for _, usr := range users {
		if usr.ID == userId {
			user = usr
			return user, nil
		}
	}

	return user, fmt.Errorf("user not found")
}

func (db *DB) MarkUserAsChirpRed(userId int) (bool, error) {
	data, err := db.loadDB()

	if err != nil {
		return false, err
	}

	user, err := db.GetUserById(userId)

	if err != nil {
		return false, err
	}

	user.IsChirpyRed = true

	data.Users[userId] = user

	err = db.writeDB(data)

	if err != nil {
		return false, err
	}

	return true, nil
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

	user.Password = hashPass

	data.Users[id] = user

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
	err := db.ensureDB()
	dbStruct := DBStructure{
		Chirps:        make(map[int]Chirp),
		Users:         make(map[int]User),
		RefreshTokens: make(map[string]Token),
	}
	if err != nil {
		return dbStruct, err
	}
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

	return dbStruct, nil

}

func (db *DB) RevokeRefreshToken(token string) (bool, error) {
	data, err := db.loadDB()
	if err != nil {
		return false, err
	}

	data.RefreshTokens[token] = Token{
		IsRevoked: true,
		Timestamp: time.Now().UTC().String(),
	}

	err = db.writeDB(data)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *DB) IsTokenRevoked(token string) (bool, error) {
	data, err := db.loadDB()
	if err != nil {
		return false, err
	}

	_, ok := data.RefreshTokens[token]

	if !ok {
		return false, nil
	}

	return true, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
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
