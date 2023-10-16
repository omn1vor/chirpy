package database

import (
	"encoding/json"
	"os"
	"sort"
	"sync"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
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
	if err := db.ensureDB(); err != nil {
		return nil, err
	}

	return &db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (*Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	entries, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	maxId := len(entries.Chirps) + 1
	chirp := Chirp{
		Id:   maxId,
		Body: body,
	}
	entries.Chirps[maxId] = chirp

	db.writeDB(entries)
	return &chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	entries, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := make([]Chirp, 0, len(entries.Chirps))
	for _, v := range entries.Chirps {
		chirps = append(chirps, v)
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})
	return chirps, nil
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
			Chirps: map[int]Chirp{},
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
