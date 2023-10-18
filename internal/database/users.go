package database

import (
	"fmt"

	"github.com/omn1vor/chirpy/internal/dto"
)

type User struct {
	Id      int    `json:"id"`
	Email   string `json:"email"`
	PwdHash string `json:"pwd_hash"`
}

func (db *DB) CreateUser(userRequest dto.UserRequest) (*dto.UserResponse, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	entries, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	maxId := len(entries.Users) + 1
	user := User{
		Id:      maxId,
		Email:   userRequest.Email,
		PwdHash: userRequest.Password,
	}
	entries.Users[maxId] = user

	db.writeDB(entries)

	return user.ToDto(), nil
}

func (db *DB) FindUserByEmail(email string) (*User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	entries, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	for _, user := range entries.Users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, nil
}

func (db *DB) UpdateUser(id int, userRequest dto.UserRequest) (*dto.UserResponse, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	entries, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	user, ok := entries.Users[id]
	if !ok {
		return nil, fmt.Errorf("User with ID %d not found", id)
	}

	user.Email = userRequest.Email
	user.PwdHash = userRequest.Password
	entries.Users[id] = user

	db.writeDB(entries)

	return user.ToDto(), nil
}

func (user *User) ToDto() *dto.UserResponse {
	return &dto.UserResponse{
		Id:    user.Id,
		Email: user.Email,
	}
}
