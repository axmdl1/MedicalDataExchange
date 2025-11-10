package grpc

import (
	"context"

	dtpb "github.com/axmdl1/MedicalDataExchange/core-service/pkg/gen/go/data-transfer"
	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/service"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	dtpb.UnimplementedDataTransferServiceServer
	dataTransferService service.DataTransferService
}

func RegisterServerAPI(grpcServer *grpc.Server, dtService service.DataTransferService) {
	dtpb.RegisterDataTransferServiceServer(grpcServer, &serverAPI{
		dataTransferService: dtService,
	})
}

// --- MedicalData ---

func (s *serverAPI) CreateMedicalData(ctx context.Context, req *dtpb.CreateMedicalDataRequest) (*dtpb.CreateMedicalDataResponse, error) {
	log := zerolog.Ctx(ctx).With().
		Str("component", "data_transfer_grpc").
		Str("method", "CreateMedicalData").
		Logger()

	if req.MedicalData == nil {
		return nil, status.Error(codes.InvalidArgument, "medical_data is required")
	}

	resp, err := s.dataTransferService.CreateMedicalData(ctx, req.MedicalData)
	if err != nil {
		log.Error().Err(err).Msg("service_failed")
		return nil, status.Errorf(codes.Internal, "failed to create medical data: %v", err)
	}

	return &dtpb.CreateMedicalDataResponse{MedicalData: resp}, nil
}

func (s *serverAPI) GetMedicalData(ctx context.Context, req *dtpb.GetMedicalDataRequest) (*dtpb.GetMedicalDataResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	resp, err := s.dataTransferService.GetMedicalData(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get medical data: %v", err)
	}

	return &dtpb.GetMedicalDataResponse{MedicalData: resp}, nil
}

func (s *serverAPI) ListMedicalData(ctx context.Context, req *dtpb.ListMedicalDataRequest) (*dtpb.ListMedicalDataResponse, error) {
	resp, err := s.dataTransferService.ListMedicalData(ctx, req.UserId, req.ClinicId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list medical data: %v", err)
	}

	return &dtpb.ListMedicalDataResponse{MedicalData: resp}, nil
}

func (s *serverAPI) DeleteMedicalData(ctx context.Context, req *dtpb.DeleteMedicalDataRequest) (*dtpb.DeleteMedicalDataResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.dataTransferService.DeleteMedicalData(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete medical data: %v", err)
	}

	return &dtpb.DeleteMedicalDataResponse{Success: true}, nil
}

// --- DataTransfer ---

func (s *serverAPI) CreateDataTransfer(ctx context.Context, req *dtpb.CreateDataTransferRequest) (*dtpb.CreateDataTransferResponse, error) {
	if req.MedicalDataId == 0 || req.FromClinicId == 0 || req.ToClinicId == 0 || req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "all fields are required")
	}

	transfer := &dtpb.DataTransfer{
		MedicalDataId: req.MedicalDataId,
		FromClinicId:  req.FromClinicId,
		ToClinicId:    req.ToClinicId,
		UserId:        req.UserId,
	}

	resp, err := s.dataTransferService.CreateDataTransfer(ctx, transfer)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create data transfer: %v", err)
	}

	return &dtpb.CreateDataTransferResponse{Transfer: resp}, nil
}

func (s *serverAPI) GetDataTransfer(ctx context.Context, req *dtpb.GetDataTransferRequest) (*dtpb.GetDataTransferResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	resp, err := s.dataTransferService.GetDataTransfer(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get data transfer: %v", err)
	}

	return &dtpb.GetDataTransferResponse{Transfer: resp}, nil
}

func (s *serverAPI) ListDataTransfers(ctx context.Context, req *dtpb.ListDataTransfersRequest) (*dtpb.ListDataTransfersResponse, error) {
	var statusFilter *string
	if req.Status != nil {
		statusFilter = req.Status
	}

	resp, err := s.dataTransferService.ListDataTransfers(ctx, req.ClinicId, req.UserId, statusFilter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list data transfers: %v", err)
	}

	return &dtpb.ListDataTransfersResponse{Transfers: resp}, nil
}

func (s *serverAPI) HandleDataTransferDecision(ctx context.Context, req *dtpb.HandleDataTransferDecisionRequest) (*dtpb.HandleDataTransferDecisionResponse, error) {
	if req.TransferId == 0 || req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "transfer_id and user_id are required")
	}

	resp, err := s.dataTransferService.HandleDataTransferDecision(ctx, req.TransferId, req.UserId, req.Confirm)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to handle data transfer decision: %v", err)
	}

	return &dtpb.HandleDataTransferDecisionResponse{Transfer: resp}, nil
}
