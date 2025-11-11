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

func (m *JWTManager) Sign(userID int64, role string, clinicID *int64, email string) (string, time.Time, error) {
	expirationTime := time.Now().Add(m.ttl)

	claims := jwt.MapClaims{
		"sub":   userID,                // subject — ID пользователя
		"role":  role,                  // patient | employee | admin
		"email": email,                 // для быстрой идентификации
		"exp":   expirationTime.Unix(), // время истечения токена
		"iat":   time.Now().Unix(),     // время выдачи
		"iss":   "user-service",        // кто выдал токен
	}

	if clinicID != nil {
		claims["clinic_id"] = *clinicID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
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
