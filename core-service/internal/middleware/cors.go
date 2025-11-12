package middleware

import (
	"net/http"
)

// CORS middleware для разрешения запросов с фронтенда
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с любых источников (в продакшене укажите конкретный домен)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Разрешаем методы
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		// Разрешаем заголовки
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// Разрешаем credentials (куки, авторизация)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Время кеширования preflight запроса (в секундах)
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Обрабатываем preflight запрос
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
