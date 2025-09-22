package service

import (
	"errors"
	"strings"
	"fmt"

	"github.com/nileshshrs/infinite-storage/model"
	"github.com/nileshshrs/infinite-storage/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UserRepo
}

func NewAuthService(repo *repository.UserRepo) *AuthService {
	return &AuthService{repo: repo}
}

type RegisterInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *AuthService) Register(input RegisterInput) (*model.User, error) {
	fmt.Println(input)
	email := strings.ToLower(input.Email)
	username := strings.ToLower(input.Username)

	if existing, _ := s.repo.FindByEmail(email); existing != nil {
		return nil, errors.New("email already exists")
	}

	if existing, _ := s.repo.FindByUsername(username); existing != nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
		Role:     "user",
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}
