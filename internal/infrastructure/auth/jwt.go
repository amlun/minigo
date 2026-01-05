package auth

import (
	"time"

	"minigo/internal/infrastructure/config"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int64  `json:"userId"`
	UserRole string `json:"userRole"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int64, userRole string, ttl time.Duration) (string, error) {
	secret := config.GetJWTSecret()
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	claims := Claims{
		UserID:   userID,
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenStr string) (*Claims, error) {
	secret := config.GetJWTSecret()
	var claims Claims
	_, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return &claims, nil
}
