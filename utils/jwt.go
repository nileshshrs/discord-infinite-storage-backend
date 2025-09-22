package utils

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/golang-jwt/jwt/v5"
)

// TokenOptions defines signing options
type TokenOptions struct {
	Secret    []byte
	ExpiresIn time.Duration
	Audience  []string
}

var (
	accessTokenOptions = TokenOptions{
		Secret:    []byte("access-secret-key"),
		ExpiresIn: 15 * time.Minute,
		Audience:  []string{"user"},
	}

	refreshTokenOptions = TokenOptions{
		Secret:    []byte("refresh-secret-key"),
		ExpiresIn: 30 * 24 * time.Hour,
		Audience:  []string{"user"},
	}
)

// SignTokens generates access and refresh tokens
func SignTokens(userID, sessionID primitive.ObjectID) (string, string, error) {
	accessClaims := jwt.MapClaims{
		"userID":    userID.Hex(),
		"sessionID": sessionID.Hex(),
		"aud":       accessTokenOptions.Audience,
		"exp":       time.Now().Add(accessTokenOptions.ExpiresIn).Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"sessionID": sessionID.Hex(),
		"aud":       refreshTokenOptions.Audience,
		"exp":       time.Now().Add(refreshTokenOptions.ExpiresIn).Unix(),
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(accessTokenOptions.Secret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(refreshTokenOptions.Secret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// VerifyAccessToken verifies access JWT
func VerifyAccessToken(tokenStr string) (jwt.MapClaims, error) {
	return verifyToken(tokenStr, accessTokenOptions)
}

// VerifyRefreshToken verifies refresh JWT
func VerifyRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	return verifyToken(tokenStr, refreshTokenOptions)
}

func verifyToken(tokenStr string, opts TokenOptions) (jwt.MapClaims, error) {
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
