package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//go:generate mockgen -source=service.go -package=auth -destination=mock_service.go
type Service interface {
	GenerateToken(userID uint) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type TokenConfig struct {
	SecretKey string
	Issuer    string
}

type serviceImpl struct {
	tokenConfig TokenConfig
}

func NewService(tokenConfig TokenConfig) Service {
	return &serviceImpl{
		tokenConfig: tokenConfig,
	}
}

func (s *serviceImpl) GenerateToken(userID uint) (string, error) {
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

func (s *serviceImpl) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// return the secret signing key
		return []byte(s.tokenConfig.SecretKey), nil
	})
}
