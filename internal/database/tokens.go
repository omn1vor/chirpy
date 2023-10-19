package database

import "time"

func (db *DB) TokenIsNotRevoked(token string) (bool, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	entries, err := db.loadDB()
	if err != nil {
		return false, err
	}

	_, ok := entries.RevokedTokens[token]
	return !ok, nil
}

func (db *DB) RevokeToken(token string) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	entries, err := db.loadDB()
	if err != nil {
		return err
	}

	entries.RevokedTokens[token] = time.Now()

	if err = db.writeDB(entries); err != nil {
		return err
	}

	return nil
}
