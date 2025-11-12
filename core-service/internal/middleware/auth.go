package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const AuthTokenKey contextKey = "auth_token"

// ExtractAuthToken извлекает Authorization header из HTTP запроса
// и сохраняет его в контексте для последующего использования в gRPC вызовах
func ExtractAuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем Authorization header
		authHeader := r.Header.Get("Authorization")

		// Если токен есть, добавляем его в контекст
		if authHeader != "" {
			ctx := context.WithValue(r.Context(), AuthTokenKey, authHeader)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// GetAuthToken извлекает токен из контекста
func GetAuthToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(AuthTokenKey).(string)
	return token, ok
}
