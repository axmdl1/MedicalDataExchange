package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(secret string, ttlMin int) *JWTManager {
	return &JWTManager{secret: []byte(secret), ttl: time.Duration(ttlMin) * time.Minute}
}

func (m *JWTManager) Sign(userID int64, role string) (string, time.Time, error) {
	exp := time.Now().Add(m.ttl)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  exp.Unix(),
		"iat":  time.Now().Unix(),
		"iss":  "user-service",
	})
	s, err := t.SignedString(m.secret)
	return s, exp, err
}

func (m *JWTManager) Verify(token string) (jwt.MapClaims, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("bad alg")
		}
		return m.secret, nil
	})
	if err != nil || !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}
