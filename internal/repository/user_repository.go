package repository

import (
	"errors"
	"sync"
)

// UserRepository определяет интерфейс работы с пользователями.
type UserRepository interface {
	CreateUser(username, hashedPassword string) error
	GetUser(username string) (string, error)
}

type inMemoryUserRepository struct {
	users map[string]string
	mu    sync.Mutex
}

func NewUserRepository() UserRepository {
	return &inMemoryUserRepository{
		users: make(map[string]string),
	}
}

func (r *inMemoryUserRepository) CreateUser(username, hashedPassword string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[username]; exists {
		return errors.New("user already exists")
	}
	r.users[username] = hashedPassword
	return nil
}

func (r *inMemoryUserRepository) GetUser(username string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if password, exists := r.users[username]; exists {
		return password, nil
	}
	return "", errors.New("user not found")
}
