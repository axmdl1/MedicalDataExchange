package grpc

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/axmdl1/MedicalDataExchange/user-service/internal/service"
)

func authUnary(jwtm *service.JWTManager) grpc.UnaryServerInterceptor {
	public := map[string]bool{
		"/user.v1.UserService/Login":      true,
		"/user.v1.UserService/CreateUser": true,
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if public[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}
		auths := md.Get("authorization")
		if len(auths) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization")
		}
		token := strings.TrimSpace(strings.TrimPrefix(auths[0], "Bearer"))
		claims, err := jwtm.ParseClaims(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}
		return handler(withClaims(ctx, claims), req)
	}
}

type policy struct {
	AllowPatients  bool
	AllowEmployees bool
	AllowAdmins    bool
	// Публично без токена
	Public bool
}

var methodPolicies = map[string]policy{
	"/user.v1.UserService/Login":      {Public: true},
	"/user.v1.UserService/CreateUser": {Public: true, AllowEmployees: true, AllowAdmins: true},
	"/user.v1.UserService/GetUser":    {AllowPatients: true, AllowEmployees: true, AllowAdmins: true},
	"/user.v1.UserService/UpdateUser": {AllowPatients: true, AllowEmployees: true, AllowAdmins: true},
	"/user.v1.UserService/DeleteUser": {AllowEmployees: true, AllowAdmins: true},
	"/user.v1.UserService/ListUsers":  {AllowEmployees: true, AllowAdmins: true},
}

func rbacUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		pol, ok := methodPolicies[info.FullMethod]
		if !ok {
			// по умолчанию закрываем
			return nil, status.Error(codes.PermissionDenied, "no policy")
		}
		if pol.Public {
			// CreateUser публичен только для регистрации пациента — это проверим в handler.
			return handler(ctx, req)
		}
		claims, ok := claimsFromCtx(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "no auth")
		}
		switch claims.Role {
		case "admin":
			if !pol.AllowAdmins {
				return nil, status.Error(codes.PermissionDenied, "forbidden")
			}
		case "employee":
			if !pol.AllowEmployees {
				return nil, status.Error(codes.PermissionDenied, "forbidden")
			}
		case "patient":
			if !pol.AllowPatients {
				return nil, status.Error(codes.PermissionDenied, "forbidden")
			}
		default:
			return nil, status.Error(codes.PermissionDenied, "unknown role")
		}
		return handler(ctx, req)
	}
}
