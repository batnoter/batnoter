package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Service represents an auth service.
// It provides methods to auth token generation and validation etc.
//go:generate mockgen -source=service.go -package=auth -destination=mock_service.go
type Service interface {
	GenerateToken(userID uint) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

// TokenConfig fields will be used to generate & parse the JWT token.
type TokenConfig struct {
	SecretKey string
	Issuer    string
}

type service struct {
	tokenConfig TokenConfig
}

// NewService creates and returns a new auth service.
func NewService(tokenConfig TokenConfig) Service {
	return &service{
		tokenConfig: tokenConfig,
	}
}

// GenerateToken creates and returns a jwt string token for a user id.
// It returns jwt token string along with any error occurred while creating the token.
func (s *service) GenerateToken(userID uint) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   strconv.FormatUint(uint64(userID), 10),
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		Issuer:    s.tokenConfig.Issuer,
		IssuedAt:  time.Now().Unix(),
	}

	// create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// generate signed token using the secret signing key
	t, err := token.SignedString([]byte(s.tokenConfig.SecretKey))
	if err != nil {
		return "", err
	}
	return t, nil
}

// ValidateToken checks the validity of given jwt token string.
// It returns a jwt token along with any error occurred while validating the token.
func (s *service) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// return the secret signing key
		return []byte(s.tokenConfig.SecretKey), nil
	})
}
