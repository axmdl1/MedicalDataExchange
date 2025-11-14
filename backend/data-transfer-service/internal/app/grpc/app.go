package grpcapp

import (
	"fmt"
	"net"

	grpcapi "github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/grpc"
	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/service"
	"google.golang.org/grpc"
)

type App struct {
	srv  *grpc.Server
	port int
}

func New(port int, svc service.DataTransferService, jwtManager *service.JWTManager) *App {
	// Создаем gRPC сервер с интерцепторами для JWT и RBAC
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcapi.AuthInterceptor(jwtManager),
			grpcapi.RBACInterceptor(),
		),
	)
	grpcapi.RegisterServerAPI(s, svc)

	return &App{
		srv:  s,
		port: port,
	}
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return err
	}

	fmt.Printf("gRP server running on :%d\n", a.port)
	return a.srv.Serve(lis)
}

func (a *App) Stop() {
	fmt.Println("stopping gRPC server...")
	a.srv.GracefulStop()
}
