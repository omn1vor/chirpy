package database

import (
	"sort"

	"github.com/omn1vor/chirpy/internal/dto"
)

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(chirpDto dto.ChirpDto) (*Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	entries, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	maxId := len(entries.Chirps) + 1
	chirp := Chirp{
		Id:       maxId,
		Body:     chirpDto.Body,
		AuthorId: chirpDto.AuthorId,
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

func (db *DB) GetChirp(id int) (*Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	entries, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirp, ok := entries.Chirps[id]
	if !ok {
		return nil, nil
	}
	return &chirp, nil
}
