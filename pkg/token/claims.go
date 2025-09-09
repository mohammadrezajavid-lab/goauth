package token

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Claims contains the payload data of the token.
type Claims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"user_id"`
}

// NewClaims creates new JWT claims for a user.
func NewClaims(userID int64, duration time.Duration) *Claims {
	return &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	}
}
