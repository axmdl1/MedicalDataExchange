package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/medical-data-exchange/blockchain-service/internal/fabric"
	"github.com/medical-data-exchange/blockchain-service/internal/models"
)

type BlockchainService struct {
	fabricClient *fabric.FabricClient
}

func NewBlockchainService(fabricClient *fabric.FabricClient) *BlockchainService {
	return &BlockchainService{
		fabricClient: fabricClient,
	}
}

// CreateAccessRequest создает запрос на доступ к данным
func (s *BlockchainService) CreateAccessRequest(ctx context.Context, req *models.AccessRequestInput) (*models.AccessRequest, error) {
	contract := s.fabricClient.GetContract()

	_, err := contract.SubmitTransaction(
		"CreateAccessRequest",
		req.ID,
		req.PatientID,
		req.ClinicID,
		req.MedicalDataID,
		req.RequestType,
		req.FromClinicID,
		req.ToClinicID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create access request: %w", err)
	}

	// Получаем созданный запрос
	return s.GetAccessRequest(ctx, req.ID)
}

// ApproveAccessRequest одобряет запрос на доступ
func (s *BlockchainService) ApproveAccessRequest(ctx context.Context, requestID, approverID string) (*models.AccessRequest, error) {
	contract := s.fabricClient.GetContract()

	_, err := contract.SubmitTransaction(
		"ApproveAccessRequest",
		requestID,
		approverID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to approve access request: %w", err)
	}

	return s.GetAccessRequest(ctx, requestID)
}

// RejectAccessRequest отклоняет запрос на доступ
func (s *BlockchainService) RejectAccessRequest(ctx context.Context, requestID string) (*models.AccessRequest, error) {
	contract := s.fabricClient.GetContract()

	_, err := contract.SubmitTransaction(
		"RejectAccessRequest",
		requestID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to reject access request: %w", err)
	}

	return s.GetAccessRequest(ctx, requestID)
}

// GetAccessRequest получает запрос на доступ
func (s *BlockchainService) GetAccessRequest(ctx context.Context, requestID string) (*models.AccessRequest, error) {
	contract := s.fabricClient.GetContract()

	result, err := contract.EvaluateTransaction("GetAccessRequest", requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get access request: %w", err)
	}

	var request models.AccessRequest
	err = json.Unmarshal(result, &request)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal access request: %w", err)
	}

	return &request, nil
}

// ValidateAccess проверяет действительность токена доступа
func (s *BlockchainService) ValidateAccess(ctx context.Context, accessID string) (*models.TemporaryAccess, error) {
	contract := s.fabricClient.GetContract()

	result, err := contract.EvaluateTransaction("ValidateAccess", accessID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate access: %w", err)
	}

	var access models.TemporaryAccess
	err = json.Unmarshal(result, &access)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal temporary access: %w", err)
	}

	return &access, nil
}

// RevokeAccess отзывает доступ
func (s *BlockchainService) RevokeAccess(ctx context.Context, accessID string) error {
	contract := s.fabricClient.GetContract()

	_, err := contract.SubmitTransaction("RevokeAccess", accessID)
	if err != nil {
		return fmt.Errorf("failed to revoke access: %w", err)
	}

	return nil
}

// LogDataTransfer записывает передачу данных между клиниками в блокчейн
func (s *BlockchainService) LogDataTransfer(ctx context.Context, transfer *models.DataTransferInput) error {
	contract := s.fabricClient.GetContract()

	_, err := contract.SubmitTransaction(
		"LogDataTransfer",
		transfer.ID,
		transfer.MedicalDataID,
		transfer.PatientID,
		transfer.FromClinicID,
		transfer.ToClinicID,
		transfer.ApprovedBy,
		transfer.DataHash,
	)
	if err != nil {
		return fmt.Errorf("failed to log data transfer: %w", err)
	}

	return nil
}

// GetTransferLog получает лог передачи данных
func (s *BlockchainService) GetTransferLog(ctx context.Context, transferID string) (*models.DataTransferLog, error) {
	contract := s.fabricClient.GetContract()

	result, err := contract.EvaluateTransaction("GetTransferLog", transferID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transfer log: %w", err)
	}

	var log models.DataTransferLog
	err = json.Unmarshal(result, &log)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal transfer log: %w", err)
	}

	return &log, nil
}
