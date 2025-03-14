package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"geo-service/internal/entities"
	"github.com/lib/pq"
)

type UserRepository interface {
	Register(user *entities.User) error
	FindByUsername(username string) (*entities.User, error)
}

type DBUserRepository struct {
	db *sql.DB
}

func NewDBUserRepository(db *sql.DB) *DBUserRepository {
	return &DBUserRepository{db: db}
}

func (r *DBUserRepository) Register(user *entities.User) error {
	query := `
        INSERT INTO users (username, password)
        VALUES ($1, $2)
    `
	_, err := r.db.ExecContext(context.Background(), query, user.Username, user.Password)
	if err != nil {
		if isDuplicateError(err) {
			return fmt.Errorf("user already exists")
		}
		return fmt.Errorf("failed to register user: %w", err)
	}
	return nil
}

func (r *DBUserRepository) FindByUsername(username string) (*entities.User, error) {
	query := `
        SELECT id, username, password
        FROM users
        WHERE username = $1
    `
	var user entities.User
	err := r.db.QueryRowContext(context.Background(), query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

func isDuplicateError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return true
	}
	return false
}
