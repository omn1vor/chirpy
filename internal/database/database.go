package database

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp        `json:"chirps"`
	Users         map[int]User         `json:"users"`
	RevokedTokens map[string]time.Time `json:"revoked_tokens"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	if err := db.ensureDB(); err != nil {
		return nil, err
	}

	return &db, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	file, err := os.OpenFile(db.path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	stats, err := file.Stat()
	if err != nil {
		return err
	}

	if stats.Size() == 0 {
		entries := DBStructure{
			Chirps:        map[int]Chirp{},
			Users:         map[int]User{},
			RevokedTokens: map[string]time.Time{},
		}
		err = db.writeDB(&entries)
		if err != nil {
			return err
		}
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (*DBStructure, error) {
	bytes, err := os.ReadFile(db.path)
	if err != nil {
		return nil, err
	}

	entries := DBStructure{}
	if err := json.Unmarshal(bytes, &entries); err != nil {
		return nil, err
	}

	return &entries, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure *DBStructure) error {
	bytes, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, bytes, 0666)
	if err != nil {
		return err
	}
	return nil
}
