package service

import (
	"errors"
	"strings"
	"time"

	"github.com/nileshshrs/infinite-storage/model"
	"github.com/nileshshrs/infinite-storage/repository"
	"github.com/nileshshrs/infinite-storage/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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

// RegisterInput is the expected payload from the client
type RegisterInput struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	UserAgent string `json:"userAgent"` // optional client header
}

type LoginInput struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	Password        string `json:"password"`
	UserAgent       string `json:"userAgent"`
}

const refreshThreshold = 24 * time.Hour

// Register registers a new user, creates a session, and signs JWT tokens
func (s *AuthService) Register(input RegisterInput) (*model.User, string, string, error) {
	email := strings.ToLower(input.Email)
	username := strings.ToLower(input.Username)

	// 1. Check if email or username already exists
	if existing, _ := s.userRepo.FindByEmail(email); existing != nil {
		return nil, "", "", errors.New("email already exists")
	}
	if existing, _ := s.userRepo.FindByUsername(username); existing != nil {
		return nil, "", "", errors.New("username already exists")
	}

	// 2. Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", err
	}

	// 3. Create user
	user := &model.User{
		Email:     email,
		Username:  username,
		Password:  string(hashedPassword),
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert user and get MongoDB ObjectID
	insertResult, err := s.userRepo.Create(user)
	if err != nil {
		return nil, "", "", err
	}

	oid, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, "", "", errors.New("failed to convert inserted ID to ObjectID")
	}
	user.ID = oid

	// 4. Create session
	session := &model.Session{
		UserID:    user.ID,
		UserAgent: input.UserAgent,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, "", "", err
	}

	// 5. Sign JWT tokens using utils
	accessToken, refreshToken, err := utils.SignTokens(user.ID, session.ID)
	if err != nil {
		return nil, "", "", err
	}

	// 6. Remove password before returning
	user.Password = ""

	return user, accessToken, refreshToken, nil
}


func (s *AuthService) Login(input LoginInput) (*model.User, string, string, error) {
	identifier := strings.ToLower(input.UsernameOrEmail)

	var user *model.User
	var err error

	// Try email first
	user, err = s.userRepo.FindByEmail(identifier)
	if user == nil || err != nil {
		// Try username
		user, err = s.userRepo.FindByUsername(identifier)
		if user == nil || err != nil {
			return nil, "", "", errors.New("invalid email/username or password")
		}
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, "", "", errors.New("invalid email/username or password")
	}

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

	// Generate tokens
	accessToken, refreshToken, err := utils.SignTokens(user.ID, session.ID)
	if err != nil {
		return nil, "", "", err
	}

	user.Password = ""

	return user, accessToken, refreshToken, nil
}


func (s *AuthService) RefreshTokens(refreshToken string) (string, string, error) {
	// 1. Verify refresh token
	claims, err := utils.VerifyRefreshToken(refreshToken)
	if err != nil || claims == nil {
		return "", "", errors.New("invalid refresh token")
	}

	// 2. Extract session ID
	sessionIDHex, ok := claims["sessionID"].(string)
	if !ok {
		return "", "", errors.New("invalid refresh token payload")
	}

	sessionID, err := primitive.ObjectIDFromHex(sessionIDHex)
	if err != nil {
		return "", "", errors.New("invalid session ID in token")
	}

	// 3. Find session in DB
	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil || session == nil {
		return "", "", errors.New("session not found or expired")
	}

	now := time.Now()
	if session.ExpiresAt.Before(now) {
		return "", "", errors.New("session has expired")
	}

	// 4. Refresh session if it expires within next 24 hours
	if session.ExpiresAt.Sub(now) <= refreshThreshold {
		session.ExpiresAt = now.Add(30 * 24 * time.Hour) // extend 30 days
		if err := s.sessionRepo.Update(session.ID, map[string]interface{}{
			"expiresAt": session.ExpiresAt,
		}); err != nil {
			return "", "", errors.New("failed to update session")
		}
	}

	// 5. Always generate new tokens (both access and refresh)
	accessToken, refreshTokenNew, err := utils.SignTokens(session.UserID, session.ID)
	if err != nil {
		return "", "", errors.New("failed to sign tokens")
	}

	return accessToken, refreshTokenNew, nil
}
