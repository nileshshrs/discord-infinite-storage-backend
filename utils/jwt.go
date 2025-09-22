package utils

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/golang-jwt/jwt/v5"
)

var (
	defaultAudience = []string{"user"}

	AccessTokenOptions = TokenOptions{
		Secret:    []byte("access-secret-key"), // replace with env var
		ExpiresIn: 15 * time.Minute,
		Audience:  defaultAudience,
	}

	RefreshTokenOptions = TokenOptions{
		Secret:    []byte("refresh-secret-key"), // replace with env var
		ExpiresIn: 15 * 24 * time.Hour,
		Audience:  defaultAudience,
	}
)

type TokenOptions struct {
	Secret    []byte
	ExpiresIn time.Duration
	Audience  []string
}

// SignAccessToken signs an access token
func SignAccessToken(userID, sessionID primitive.ObjectID) (string, error) {
	claims := jwt.MapClaims{
		"userID":    userID.Hex(),
		"sessionID": sessionID.Hex(),
		"aud":       AccessTokenOptions.Audience,
		"exp":       time.Now().Add(AccessTokenOptions.ExpiresIn).Unix(),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(AccessTokenOptions.Secret)
}

// SignRefreshToken signs a refresh token
func SignRefreshToken(sessionID primitive.ObjectID) (string, error) {
	claims := jwt.MapClaims{
		"sessionID": sessionID.Hex(),
		"aud":       RefreshTokenOptions.Audience,
		"exp":       time.Now().Add(RefreshTokenOptions.ExpiresIn).Unix(),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(RefreshTokenOptions.Secret)
}

// SignTokens signs both access and refresh tokens
func SignTokens(userID, sessionID primitive.ObjectID) (string, string, error) {
	accessToken, err := SignAccessToken(userID, sessionID)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := SignRefreshToken(sessionID)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// VerifyAccessToken verifies an access token
func VerifyAccessToken(tokenStr string) (map[string]interface{}, error) {
	return verifyToken(tokenStr, AccessTokenOptions)
}

// VerifyRefreshToken verifies a refresh token
func VerifyRefreshToken(tokenStr string) (map[string]interface{}, error) {
	return verifyToken(tokenStr, RefreshTokenOptions)
}

func verifyToken(tokenStr string, opts TokenOptions) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return opts.Secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
