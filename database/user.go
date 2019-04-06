package database

import (
	"fmt"
	"quizApp/model"
)

func (db *DB) AddUser(user model.User) error {
	db.cache.userMu.RLock()
	_, ok := db.cache.userCache[user.Username]
	db.cache.userMu.RUnlock()
	if ok {
		return fmt.Errorf("User already exists")
	}

	db.cache.userMu.Lock()
	db.cache.userCache[user.Username] = user
	db.cache.userMu.Unlock()

	return nil
}

func (db *DB) GetUser(username string) bool {
	db.cache.userMu.RLock()
	_, ok := db.cache.userCache[username]
	db.cache.userMu.RUnlock()
	return ok
}
