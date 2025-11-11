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
	svc service.UserService
}

func NewServer(svc service.UserService) *Server { return &Server{svc: svc} }

func (s *Server) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	gs := grpc.NewServer()
	userv1.RegisterUserServiceServer(gs, s)
	return gs.Serve(lis)
}

// ---- RPCs ----

func (s *Server) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	u := fromProto(req.User)
	created, err := s.svc.Create(ctx, u)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	out := toProto(created)
	out.Password = ""
	return &userv1.CreateUserResponse{User: out}, nil
}

func (s *Server) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	u, err := s.svc.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	out := toProto(u)
	out.Password = ""
	return &userv1.GetUserResponse{User: out}, nil
}

func (s *Server) ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	var clinicPtr *int64
	if req.ClinicId != nil {
		v := req.ClinicId.Value
		clinicPtr = &v
	}
	users, err := s.svc.List(ctx, clinicPtr, req.Type)
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
	u := fromProto(req.User)
	upd, err := s.svc.Update(ctx, u)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	out := toProto(upd)
	out.Password = ""
	return &userv1.UpdateUserResponse{User: out}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	if err := s.svc.Delete(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &userv1.DeleteUserResponse{Success: true}, nil
}

func (s *Server) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	u, token, _, err := s.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid email or password")
	}
	out := toProto(u)
	out.Password = ""
	return &userv1.LoginResponse{AccessToken: token, User: out}, nil
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
