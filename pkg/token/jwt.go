package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Maker is an interface for managing tokens.
type Maker interface {
	// CreateToken creates a new token for a specific user ID.
	CreateToken(userID int64) (string, error)
	// VerifyToken checks if the token is valid or not.
	VerifyToken(tokenString string) (*Claims, error)
}

// Config holds configuration for JWT token generation.
type Config struct {
	SecretKey      string        `koanf:"secret_key"`
	ExpirationTime time.Duration `koanf:"expiration_time"`
}

// JWTMaker implements the Maker interface using JWT.
type JWTMaker struct {
	config Config
}

// NewJWTMaker creates a new JWTMaker from a config.
func NewJWTMaker(config Config) (Maker, error) {
	if len(config.SecretKey) < 32 {
		return nil, errors.New("invalid key size: must be at least 32 characters")
	}
	return &JWTMaker{config: config}, nil
}

// CreateToken creates a new JWT token using the duration from its config.
func (maker *JWTMaker) CreateToken(userID int64) (string, error) {
	claims := NewClaims(userID, maker.config.ExpirationTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(maker.config.SecretKey))
}

// VerifyToken verifies a JWT token.
func (maker *JWTMaker) VerifyToken(tokenString string) (*Claims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(maker.config.SecretKey), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		}
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
