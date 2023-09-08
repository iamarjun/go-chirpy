package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
	id   int
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := DB{
		path: path,
		mux:  &sync.RWMutex{},
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
		fmt.Printf("Trying to unmarshal existing dbstructure %v\n", dat)
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

	fmt.Printf("Trying to marshal dbstructure %v\n", dbStructure)

	dat, err := json.Marshal(dbStructure)

	if err != nil {
		return err
	}

	fmt.Printf("Marshaled dbstructure %v\n", dat)

	err = os.WriteFile("database.json", dat, os.ModeAppend)
	if err != nil {
		return err
	}

	return nil
}
