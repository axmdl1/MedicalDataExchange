package service

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type JWTManager struct {
	secret []byte
}

type Claims struct {
	UserID   int64
	Role     string // "patient"|"employee"|"admin"
	ClinicID *int64 // nil для patient/admin, обязателен для employee
	Email    string
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{secret: []byte(secret)}
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

func (m *JWTManager) ParseClaims(tokenStr string) (*Claims, error) {
	claimsMap, err := m.Verify(tokenStr)
	if err != nil {
		return nil, err
	}

	var c Claims
	// sub может быть float64 или string в зависимости от encoder — аккуратно парсим
	switch v := claimsMap["sub"].(type) {
	case float64:
		c.UserID = int64(v)
	case int64:
		c.UserID = v
	case string: // если решите хранить строкой — распарсить
		// ...
	}

	if r, ok := claimsMap["role"].(string); ok {
		c.Role = r
	}

	if cid, ok := claimsMap["clinic_id"].(float64); ok {
		t := int64(cid)
		c.ClinicID = &t
	}
	if em, ok := claimsMap["email"].(string); ok {
		c.Email = em
	}

	return &c, nil
}
