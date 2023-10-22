package database

import (
	"fmt"
	"sort"
	"strings"

	"github.com/omn1vor/chirpy/internal/dto"
	"github.com/omn1vor/chirpy/internal/errs"
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
func (db *DB) GetChirps(authorId int, sorting string) ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	entries, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := make([]Chirp, 0)
	for _, v := range entries.Chirps {
		if authorId != -1 && v.AuthorId != authorId {
			continue
		}
		chirps = append(chirps, v)
	}
	sort.Slice(chirps, func(i, j int) bool {
		if strings.ToLower(sorting) == "desc" {
			return chirps[i].Id > chirps[j].Id
		} else {
			return chirps[i].Id < chirps[j].Id
		}
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
		return nil, &errs.ErrNotFound{
			Msg: fmt.Sprintf("Chirp with ID %d not found", id),
		}
	}
	return &chirp, nil
}

func (db *DB) DeleteChirp(id int) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	entries, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(entries.Chirps, id)

	return nil
}
