package client

import (
	userv1 "github.com/axmdl1/MedicalDataExchange/core-service/pkg/gen/go/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	API userv1.UserServiceClient
	// conn оставим, если захотите закрывать
	// conn *grpc.ClientConn
}

func NewUserClient(addr string) (*UserClient, error) {
	// Под ваш стиль: в data-transfer вы используете grpc.NewClient + insecure
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return &UserClient{
		API: userv1.NewUserServiceClient(conn),
		// conn: conn,
	}, nil
}

func (c *UserClient) Close() error {
	// если будете хранить conn, закрывайте его здесь
	return nil
}
