package grpc

import (
	"context"
	"net"

	_ "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"

	userv1 "github.com/axmdl1/MedicalDataExchange/core-service/pkg/gen/go/user"
	"github.com/axmdl1/MedicalDataExchange/user-service/internal/models"
	"github.com/axmdl1/MedicalDataExchange/user-service/internal/service"
)

type Server struct {
	userv1.UnimplementedUserServiceServer
	svc  service.UserService
	jwtm *service.JWTManager
}

func NewServer(svc service.UserService, jwtm *service.JWTManager) *Server {
	return &Server{svc: svc, jwtm: jwtm}
}

func (s *Server) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	gs := grpc.NewServer(grpc.ChainUnaryInterceptor(authUnary(s.jwtm), rbacUnary()))
	userv1.RegisterUserServiceServer(gs, s)
	return gs.Serve(lis)
}

// ---- RPCs ----

func (s *Server) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	u := fromProto(req.GetUser())

	claims, hasAuth := claimsFromCtx(ctx)
	if !hasAuth {
		// Публичная регистрация: разрешаем только patient без clinic_id
		if u.Type != "patient" || u.ClinicID != nil {
			return nil, status.Error(codes.PermissionDenied, "only patient self-registration is allowed without auth")
		}
	} else {
		switch claims.Role {
		case "employee":
			// сотрудник может создавать только пациентов своей клиники
			if u.Type != "patient" {
				return nil, status.Error(codes.PermissionDenied, "employees can create only patients")
			}
			if claims.ClinicID == nil || u.ClinicID == nil || *claims.ClinicID != *u.ClinicID {
				return nil, status.Error(codes.PermissionDenied, "clinic mismatch")
			}
		case "admin":
			// без ограничений
		case "patient":
			// хоть RBAC разрешил, но логически — пациент не создает других
			return nil, status.Error(codes.PermissionDenied, "patients cannot create users")
		}
	}

	created, err := s.svc.Create(ctx, u)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	out := toProto(created)
	out.Password = ""
	return &userv1.CreateUserResponse{User: out}, nil
}

func (s *Server) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := s.svc.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if claims, ok := claimsFromCtx(ctx); ok {
		switch claims.Role {
		case "patient":
			if claims.UserID != user.ID {
				return nil, status.Error(codes.PermissionDenied, "patients can view only themselves")
			}
		case "employee":
			// сотрудник может смотреть себя всегда
			if claims.UserID != user.ID {
				// чужого — только в своей клинике
				if claims.ClinicID == nil || user.ClinicID == nil || *claims.ClinicID != *user.ClinicID {
					return nil, status.Error(codes.PermissionDenied, "user not in your clinic")
				}
			}
		}
	}
	out := toProto(user)
	out.Password = ""
	return &userv1.GetUserResponse{User: out}, nil
}

func (s *Server) ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	var clinicPtr *int64
	if req.ClinicId != nil {
		v := req.ClinicId.Value
		clinicPtr = &v
	}

	if claims, ok := claimsFromCtx(ctx); ok && claims.Role == "employee" {
		// для сотрудника игнорируем чужие клиники и принудительно подставляем его клинику
		if claims.ClinicID == nil {
			return nil, status.Error(codes.PermissionDenied, "employee has no clinic assigned")
		}
		clinicPtr = claims.ClinicID
	}

	users, err := s.svc.List(ctx, clinicPtr, req.GetType())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := make([]*userv1.User, 0, len(users))
	for i := range users {
		u := toProto(&users[i])
		u.Password = ""
		resp = append(resp, u)
	}
	return &userv1.ListUsersResponse{Users: resp}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	in := fromProto(req.GetUser())

	// прочитаем текущего, чтобы проверить правила
	current, err := s.svc.Get(ctx, in.ID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	if claims, ok := claimsFromCtx(ctx); ok {
		switch claims.Role {
		case "patient":
			if claims.UserID != in.ID {
				return nil, status.Error(codes.PermissionDenied, "patients can update only themselves")
			}
			// пациент не может менять свой type/clinic_id
			in.Type = current.Type
			in.ClinicID = current.ClinicID
		case "employee":
			// сотрудник может обновлять только пациентов своей клиники
			if current.Type != "patient" {
				return nil, status.Error(codes.PermissionDenied, "employees can update only patients")
			}
			if claims.ClinicID == nil || current.ClinicID == nil || *claims.ClinicID != *current.ClinicID {
				return nil, status.Error(codes.PermissionDenied, "user not in your clinic")
			}
			// не даем сотруднику менять тип пациента и clinic_id на чужую
			in.Type = "patient"
			in.ClinicID = current.ClinicID
		case "admin":
			// ок
		}
	}

	updated, err := s.svc.Update(ctx, in)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	out := toProto(updated)
	out.Password = ""
	return &userv1.UpdateUserResponse{User: out}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	target, err := s.svc.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if claims, ok := claimsFromCtx(ctx); ok {
		switch claims.Role {
		case "patient":
			return nil, status.Error(codes.PermissionDenied, "patients cannot delete users")
		case "employee":
			// сотрудник может удалять только пациентов своей клиники
			if target.Type != "patient" {
				return nil, status.Error(codes.PermissionDenied, "employees can delete only patients")
			}
			if claims.ClinicID == nil || target.ClinicID == nil || *claims.ClinicID != *target.ClinicID {
				return nil, status.Error(codes.PermissionDenied, "user not in your clinic")
			}
		case "admin":
			// ок
		}
	}
	if err := s.svc.Delete(ctx, req.GetId()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &userv1.DeleteUserResponse{Success: true}, nil
}

func (s *Server) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	u, token, _, err := s.svc.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid email or password")
	}
	out := toProto(u)
	out.Password = ""
	return &userv1.LoginResponse{
		AccessToken: token,
		User:        out,
	}, nil
}

// ---- mapping ----

func fromProto(u *userv1.User) *models.User {
	if u == nil {
		return &models.User{}
	}
	var clinicPtr *int64
	if u.ClinicId != nil {
		v := u.ClinicId.Value
		clinicPtr = &v
	}
	return &models.User{
		ID:          u.Id,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Type:        u.Type,
		Password:    u.Password, // захешируется в сервисе
		ClinicID:    clinicPtr,
	}
}

func toProto(u *models.User) *userv1.User {
	if u == nil {
		return &userv1.User{}
	}
	var cl *wrapperspb.Int64Value
	if u.ClinicID != nil {
		cl = wrapperspb.Int64(*u.ClinicID)
	}
	return &userv1.User{
		Id:          u.ID,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Type:        u.Type,
		ClinicId:    cl,
	}
}
