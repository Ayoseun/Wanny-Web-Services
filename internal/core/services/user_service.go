package services

import (
	"errors"
	"wanny-web-services/internal/core/domain"
	"wanny-web-services/internal/core/ports"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Username: username,
		Password: string(hashedPassword),
	}

	return s.repo.Create(user)
}

func (s *UserService) Authenticate(username, password string) (*domain.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}