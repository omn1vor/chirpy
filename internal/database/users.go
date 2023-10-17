package database

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) CreateUser(email string) (*User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	entries, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	maxId := len(entries.Users) + 1
	user := User{
		Id:    maxId,
		Email: email,
	}
	entries.Users[maxId] = user

	db.writeDB(entries)
	return &user, nil
}
