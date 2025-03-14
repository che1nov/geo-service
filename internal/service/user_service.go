package service

import (
	"errors"
	"geo-service/internal/entities"
	"geo-service/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(username, password string) error {
	user := &entities.User{
		Username: username,
		Password: password,
	}
	if err := user.HashPassword(); err != nil {
		return err
	}
	return s.repo.Register(user)
}

func (s *UserService) Login(username, password string) (*entities.User, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}
	if !user.CheckPassword(password) {
		return nil, errors.New("invalid username or password")
	}
	return user, nil
}
