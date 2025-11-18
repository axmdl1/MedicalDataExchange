package grpc

import (
	"context"
	"log"
	"reflect"
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

// AuthInterceptor
func AuthInterceptor(jwtManager *service.JWTManager) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}

		authHeader := values[0]
		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization format")
		}
		token := strings.TrimPrefix(authHeader, prefix)

		claims, err := jwtManager.ParseClaims(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token: "+err.Error())
		}

		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

		log.Printf("AUTH: user_id=%d, role=%s", claims.UserID, claims.Role)
		return handler(ctx, req)
	}
}

// RBACInterceptor
func RBACInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		role, ok := ctx.Value(UserRoleKey).(string)
		if !ok {
			log.Println("RBAC: missing role in context")
			return nil, status.Error(codes.PermissionDenied, "missing role in context")
		}

		userID, _ := ctx.Value(UserIDKey).(int64)
		method := info.FullMethod

		log.Printf("RBAC: method=%s, role=%s, user_id=%d", method, role, userID)

		// Создание записей
		if strings.Contains(method, "CreateMedicalData") {
			if role == "patient" {
				log.Println("RBAC: patient tried to create medical data → denied")
				return nil, status.Error(codes.PermissionDenied, "patients cannot create medical records")
			}
			log.Println("RBAC: employee/admin → create allowed")
			return handler(ctx, req)
		}

		// employee / admin
		if role == "employee" || role == "admin" {
			log.Println("RBAC: employee/admin → full access")
			return handler(ctx, req)
		}

		// patient
		if role == "patient" {
			log.Println("RBAC: patient access check")

			// Проверяем Patient Access методы
			if strings.Contains(method, "PatientAccess") {
				// Для Patient Access проверяем patient_id
				hasPatientID := hasPatientIDInRequest(req)
				log.Printf("RBAC: hasPatientIDInRequest = %v", hasPatientID)

				if hasPatientID {
					matches := patientIDMatches(req, userID)
					log.Printf("RBAC: patientIDMatches = %v (expected: %d)", matches, userID)
					if !matches {
						log.Println("RBAC: patient_id mismatch → denied")
						return nil, status.Error(codes.PermissionDenied, "you can only access your own requests")
					}
				} else {
					// Если patient_id не указан, принудительно устанавливаем из токена
					log.Printf("RBAC: no patient_id in request, will be set from token (user_id=%d)", userID)
				}
			} else {
				// Для других методов проверяем user_id
				hasID := hasUserIDInRequest(req)
				log.Printf("RBAC: hasUserIDInRequest = %v", hasID)

				if hasID {
					matches := userIDMatches(req, userID)
					log.Printf("RBAC: userIDMatches = %v (expected: %d)", matches, userID)
					if !matches {
						log.Println("RBAC: user_id mismatch → denied")
						return nil, status.Error(codes.PermissionDenied, "you can only access your own data")
					}
				} else {
					log.Println("RBAC: no user_id in request → allowing (service will auto-fill)")
				}
			}

			log.Println("RBAC: patient access granted")
			return handler(ctx, req)
		}

		log.Println("RBAC: unknown role → denied")
		return nil, status.Error(codes.PermissionDenied, "insufficient permissions")
	}
}

// hasUserIDInRequest — проверяет, заполнено ли поле user_id
func hasUserIDInRequest(req interface{}) bool {
	v := reflect.ValueOf(req)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		log.Println("RBAC: req is not struct")
		return false
	}

	log.Printf("RBAC: checking request type: %s", v.Type().Name())

	// Проверяем UserId (protobuf генерирует как UserId)
	if f := v.FieldByName("UserId"); f.IsValid() {
		// Используем Interface() + приведение
		if val, ok := f.Interface().(int64); ok && val != 0 {
			log.Printf("RBAC: UserId = %d (via Interface)", val)
			return true
		}
		log.Printf("RBAC: UserId is zero or not int64")
	}

	// Проверяем в Filter
	if filter := v.FieldByName("Filter"); filter.IsValid() && !filter.IsNil() {
		log.Println("RBAC: found Filter")
		if filter.Kind() == reflect.Ptr {
			filter = filter.Elem()
		}
		if f := filter.FieldByName("UserId"); f.IsValid() {
			if val, ok := f.Interface().(int64); ok && val != 0 {
				log.Printf("RBAC: Filter.UserId = %d", val)
				return true
			}
		}
	}

	log.Println("RBAC: no user_id found")
	return false
}

// userIDMatches — сравнивает user_id с expected
func userIDMatches(req interface{}, expected int64) bool {
	v := reflect.ValueOf(req)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Проверяем UserId
	if f := v.FieldByName("UserId"); f.IsValid() {
		if val, ok := f.Interface().(int64); ok {
			log.Printf("RBAC: UserId = %d, expected = %d", val, expected)
			return val == expected
		}
	}

	// Проверяем в Filter
	if filter := v.FieldByName("Filter"); filter.IsValid() && !filter.IsNil() {
		if filter.Kind() == reflect.Ptr {
			filter = filter.Elem()
		}
		if f := filter.FieldByName("UserId"); f.IsValid() {
			if val, ok := f.Interface().(int64); ok {
				log.Printf("RBAC: Filter.UserId = %d, expected = %d", val, expected)
				return val == expected
			}
		}
	}

	log.Println("RBAC: no user_id field found for comparison")
	return false
}

// hasPatientIDInRequest — проверяет, заполнено ли поле patient_id
func hasPatientIDInRequest(req interface{}) bool {
	v := reflect.ValueOf(req)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false
	}

	// Проверяем PatientId (protobuf генерирует как PatientId)
	if f := v.FieldByName("PatientId"); f.IsValid() {
		if val, ok := f.Interface().(*int64); ok && val != nil && *val != 0 {
			log.Printf("RBAC: PatientId = %d", *val)
			return true
		}
		// Также проверяем как int64 (не указатель)
		if val, ok := f.Interface().(int64); ok && val != 0 {
			log.Printf("RBAC: PatientId = %d (non-pointer)", val)
			return true
		}
	}

	return false
}

// patientIDMatches — сравнивает patient_id с expected
func patientIDMatches(req interface{}, expected int64) bool {
	v := reflect.ValueOf(req)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Проверяем PatientId
	if f := v.FieldByName("PatientId"); f.IsValid() {
		// Проверяем как указатель
		if val, ok := f.Interface().(*int64); ok && val != nil {
			log.Printf("RBAC: PatientId = %d, expected = %d", *val, expected)
			return *val == expected
		}
		// Проверяем как int64
		if val, ok := f.Interface().(int64); ok {
			log.Printf("RBAC: PatientId = %d, expected = %d", val, expected)
			return val == expected
		}
	}

	return false
}

// GetUserID / GetUserRole
func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

func GetUserRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(UserRoleKey).(string)
	return role, ok
}
