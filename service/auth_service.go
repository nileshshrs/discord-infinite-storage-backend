package service

import (
	"errors"
	"strings"
	"time"

	"github.com/nileshshrs/infinite-storage/model"
	"github.com/nileshshrs/infinite-storage/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userRepo    *repository.UserRepo
	sessionRepo *repository.SessionRepo
}

func NewAuthService(userRepo *repository.UserRepo, sessionRepo *repository.SessionRepo) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

type RegisterInput struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	UserAgent string `json:"userAgent"`
}

// Secret keys (store in env in production)
var (
	accessSecret  = []byte("access-secret-key")
	refreshSecret = []byte("refresh-secret-key")
)

func SignTokens(userID, sessionID primitive.ObjectID) (string, string, error) {
	accessClaims := jwt.MapClaims{
		"userID":    userID.Hex(),
		"sessionID": sessionID.Hex(),
		"exp":       time.Now().Add(15 * time.Minute).Unix(),
	}
	refreshClaims := jwt.MapClaims{
		"sessionID": sessionID.Hex(),
		"exp":       time.Now().Add(30 * 24 * time.Hour).Unix(),
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(accessSecret)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(refreshSecret)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *AuthService) Register(input RegisterInput) (*model.User, string, string, error) {
	email := strings.ToLower(input.Email)
	username := strings.ToLower(input.Username)

	// Check if email or username exists
	if existing, _ := s.userRepo.FindByEmail(email); existing != nil {
		return nil, "", "", errors.New("email already exists")
	}
	if existing, _ := s.userRepo.FindByUsername(username); existing != nil {
		return nil, "", "", errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", err
	}

	// Create user
	user := &model.User{
		Email:     email,
		Username:  username,
		Password:  string(hashedPassword),
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert user and get ID
	res, err := s.userRepo.Create(user)
	if err != nil {
		return nil, "", "", err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, "", "", errors.New("failed to convert inserted ID to ObjectID")
	}
	user.ID = oid

	// Create session
	session := &model.Session{
		UserID:    user.ID,
		UserAgent: input.UserAgent,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, "", "", err
	}

	// Sign tokens
	accessToken, refreshToken, err := SignTokens(user.ID, session.ID)
	if err != nil {
		return nil, "", "", err
	}

	user.Password = "" // hide password

	return user, accessToken, refreshToken, nil
}
