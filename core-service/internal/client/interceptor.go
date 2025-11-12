package client

import (
	"context"

	"github.com/axmdl1/MedicalDataExchange/core-service/internal/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthInterceptor - клиентский интерцептор для проброса Authorization токена в metadata
func AuthInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Извлекаем токен из контекста (который был добавлен middleware)
		if token, ok := middleware.GetAuthToken(ctx); ok && token != "" {
			// Добавляем токен в metadata для gRPC вызова
			md := metadata.Pairs("authorization", token)
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		// Вызываем gRPC метод с обновленным контекстом
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
