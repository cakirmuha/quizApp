package database

import (
	"context"
	"quizApp/model"
	"sync"
)

type DB struct {
	ctx context.Context

	cache *dbCache
}

type dbCache struct {
	userCache map[string]model.User
	userMu    sync.RWMutex

	userScore   map[string]int
	userscoreMu sync.RWMutex
}

func New(ctx context.Context) (*DB, error) {
	c := dbCache{
		userCache: make(map[string]model.User),
		userScore: make(map[string]int),
	}

	return &DB{
		ctx:   ctx,
		cache: &c,
	}, nil
}
