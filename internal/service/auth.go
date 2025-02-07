package service

import (
	"errors"

	"example/internal/models"
	"example/internal/repository"

	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService описывает бизнес-логику для регистрации и аутентификации.
type AuthService interface {
	Register(req models.RegisterRequest) error
	Login(req models.LoginRequest) (string, error)
}

type authService struct {
	userRepo  repository.UserRepository
	tokenAuth *jwtauth.JWTAuth
}

func NewAuthService(userRepo repository.UserRepository, tokenAuth *jwtauth.JWTAuth) AuthService {
	return &authService{
		userRepo:  userRepo,
		tokenAuth: tokenAuth,
	}
}

func (s *authService) Register(req models.RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.userRepo.CreateUser(req.Username, string(hashedPassword))
}

func (s *authService) Login(req models.LoginRequest) (string, error) {
	hashedPassword, err := s.userRepo.GetUser(req.Username)
	if err != nil {
		return "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return "", errors.New("invalid password")
	}

	_, tokenString, err := s.tokenAuth.Encode(map[string]interface{}{"user": req.Username})
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
