package grpc

import (
	"context"

	"github.com/medical-data-exchange/blockchain-service/internal/app"
	"github.com/medical-data-exchange/blockchain-service/internal/models"
	pb "github.com/axmdl1/MedicalDataExchange/core-service/pkg/gen/go/blockchain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BlockchainHandler struct {
	pb.UnimplementedBlockchainServiceServer
	service *app.BlockchainService
}

func NewBlockchainHandler(service *app.BlockchainService) *BlockchainHandler {
	return &BlockchainHandler{
		service: service,
	}
}

func (h *BlockchainHandler) CreateAccessRequest(ctx context.Context, req *pb.CreateAccessRequestRequest) (*pb.AccessRequestResponse, error) {
	input := &models.AccessRequestInput{
		ID:            req.Id,
		PatientID:     req.PatientId,
		ClinicID:      req.ClinicId,
		MedicalDataID: req.MedicalDataId,
		RequestType:   req.RequestType,
		FromClinicID:  req.FromClinicId,
		ToClinicID:    req.ToClinicId,
	}

	result, err := h.service.CreateAccessRequest(ctx, input)
	if err != nil {
		return nil, err
	}

	return toAccessRequestResponse(result), nil
}

func (h *BlockchainHandler) ApproveAccessRequest(ctx context.Context, req *pb.ApproveAccessRequestRequest) (*pb.AccessRequestResponse, error) {
	result, err := h.service.ApproveAccessRequest(ctx, req.RequestId, req.ApproverId)
	if err != nil {
		return nil, err
	}

	return toAccessRequestResponse(result), nil
}

func (h *BlockchainHandler) RejectAccessRequest(ctx context.Context, req *pb.RejectAccessRequestRequest) (*pb.AccessRequestResponse, error) {
	result, err := h.service.RejectAccessRequest(ctx, req.RequestId)
	if err != nil {
		return nil, err
	}

	return toAccessRequestResponse(result), nil
}

func (h *BlockchainHandler) GetAccessRequest(ctx context.Context, req *pb.GetAccessRequestRequest) (*pb.AccessRequestResponse, error) {
	result, err := h.service.GetAccessRequest(ctx, req.RequestId)
	if err != nil {
		return nil, err
	}

	return toAccessRequestResponse(result), nil
}

func (h *BlockchainHandler) ValidateAccess(ctx context.Context, req *pb.ValidateAccessRequest) (*pb.TemporaryAccessResponse, error) {
	result, err := h.service.ValidateAccess(ctx, req.AccessId)
	if err != nil {
		return nil, err
	}

	return &pb.TemporaryAccessResponse{
		Id:            result.ID,
		AccessToken:   result.AccessToken,
		PatientId:     result.PatientID,
		ClinicId:      result.ClinicID,
		MedicalDataId: result.MedicalDataID,
		GrantedAt:     timestamppb.New(result.GrantedAt),
		ExpiresAt:     timestamppb.New(result.ExpiresAt),
		IsRevoked:     result.IsRevoked,
	}, nil
}

func (h *BlockchainHandler) RevokeAccess(ctx context.Context, req *pb.RevokeAccessRequest) (*pb.RevokeAccessResponse, error) {
	err := h.service.RevokeAccess(ctx, req.AccessId)
	if err != nil {
		return &pb.RevokeAccessResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.RevokeAccessResponse{
		Success: true,
		Message: "Access revoked successfully",
	}, nil
}

func (h *BlockchainHandler) LogDataTransfer(ctx context.Context, req *pb.LogDataTransferRequest) (*pb.LogDataTransferResponse, error) {
	input := &models.DataTransferInput{
		ID:            req.Id,
		MedicalDataID: req.MedicalDataId,
		PatientID:     req.PatientId,
		FromClinicID:  req.FromClinicId,
		ToClinicID:    req.ToClinicId,
		ApprovedBy:    req.ApprovedBy,
		DataHash:      req.DataHash,
	}

	err := h.service.LogDataTransfer(ctx, input)
	if err != nil {
		return &pb.LogDataTransferResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.LogDataTransferResponse{
		Success: true,
		Message: "Transfer logged successfully",
	}, nil
}

func (h *BlockchainHandler) GetTransferLog(ctx context.Context, req *pb.GetTransferLogRequest) (*pb.DataTransferLogResponse, error) {
	result, err := h.service.GetTransferLog(ctx, req.TransferId)
	if err != nil {
		return nil, err
	}

	return &pb.DataTransferLogResponse{
		Id:            result.ID,
		MedicalDataId: result.MedicalDataID,
		PatientId:     result.PatientID,
		FromClinicId:  result.FromClinicID,
		ToClinicId:    result.ToClinicID,
		ApprovedBy:    result.ApprovedBy,
		TransferredAt: timestamppb.New(result.TransferredAt),
		DataHash:      result.DataHash,
	}, nil
}

func toAccessRequestResponse(req *models.AccessRequest) *pb.AccessRequestResponse {
	resp := &pb.AccessRequestResponse{
		Id:            req.ID,
		PatientId:     req.PatientID,
		ClinicId:      req.ClinicID,
		MedicalDataId: req.MedicalDataID,
		RequestType:   req.RequestType,
		Status:        req.Status,
		RequestedAt:   timestamppb.New(req.RequestedAt),
		FromClinicId:  req.FromClinicID,
		ToClinicId:    req.ToClinicID,
	}

	if req.ApprovedAt != nil {
		resp.ApprovedAt = timestamppb.New(*req.ApprovedAt)
	}

	if req.ExpiresAt != nil {
		resp.ExpiresAt = timestamppb.New(*req.ExpiresAt)
	}

	return resp
}
