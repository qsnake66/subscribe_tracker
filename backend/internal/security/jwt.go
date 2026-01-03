package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	Secret []byte
	TTL    time.Duration
}

func NewJWTManager(secret []byte, ttl time.Duration) JWTManager {
	return JWTManager{
		Secret: secret,
		TTL:    ttl,
	}
}

func (m JWTManager) Sign(userID, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(m.TTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.Secret)
}

func (m JWTManager) Parse(tokenValue string) (string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return m.Secret, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", errors.New("invalid token")
	}
	return userID, nil
}
