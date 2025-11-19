package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("supersecret") // <-- должен совпадать с user-service config.secret!

func ExtractAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization")
		if auth == "" {
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			next.ServeHTTP(w, r)
			return
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			next.ServeHTTP(w, r)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		var u UserInfo

		if v, ok := claims["sub"].(float64); ok {
			u.Id = int64(v)
		}

		if v, ok := claims["email"].(string); ok {
			u.Email = v
		}

		if v, ok := claims["role"].(string); ok {
			u.Role = v
		}

		if v, ok := claims["clinic_id"].(float64); ok {
			id := int64(v)
			u.ClinicID = &id
		}

		ctx := SetUser(r.Context(), u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
