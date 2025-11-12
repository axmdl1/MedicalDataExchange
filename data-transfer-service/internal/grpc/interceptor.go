package grpc

import (
	"context"
	"strings"

	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

// AuthInterceptor проверяет JWT токен и добавляет claims в context
func AuthInterceptor(jwtManager *service.JWTManager) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Извлекаем metadata из контекста
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		// Извлекаем Authorization header
		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}

		// Извлекаем токен из "Bearer <token>"
		authHeader := values[0]
		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization format")
		}
		token := strings.TrimPrefix(authHeader, prefix)

		// Проверяем токен и парсим claims
		claims, err := jwtManager.ParseClaims(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token: "+err.Error())
		}

		// Добавляем claims в контекст
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

		// Продолжаем обработку запроса
		return handler(ctx, req)
	}
}

// RBACInterceptor проверяет права доступа на основе роли
func RBACInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Извлекаем роль из контекста
		role, ok := ctx.Value(UserRoleKey).(string)
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "missing role in context")
		}

		// Проверяем права доступа
		// Только employee и admin могут работать с медицинскими данными
		if role != "employee" && role != "admin" {
			return nil, status.Error(codes.PermissionDenied, "insufficient permissions")
		}

		// Продолжаем обработку запроса
		return handler(ctx, req)
	}
}

// GetUserID извлекает user_id из контекста
func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

// GetUserRole извлекает роль из контекста
func GetUserRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(UserRoleKey).(string)
	return role, ok
}
