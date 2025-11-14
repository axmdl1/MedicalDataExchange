package client

import (
	dtpb "github.com/axmdl1/MedicalDataExchange/core-service/pkg/gen/go/data-transfer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DataTransferClient struct {
	Api dtpb.DataTransferServiceClient
}

func NewDataTransferClient(addr string) (*DataTransferClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(AuthInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	return &DataTransferClient{
		Api: dtpb.NewDataTransferServiceClient(conn),
	}, nil
}

func (c *DataTransferClient) Close() error {
	// В реальном проекте — закрывать conn
	return nil
}
