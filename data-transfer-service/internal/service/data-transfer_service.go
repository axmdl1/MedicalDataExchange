package service

import (
	"context"
	"fmt"
	"time"

	dtpb "github.com/axmdl1/MedicalDataExchange/core-service/pkg/gen/go/data-transfer"
	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/model"
	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/repository"
	"github.com/jinzhu/copier"
	"github.com/rs/zerolog"
)

type DataTransferService interface {
	CreateMedicalData(ctx context.Context, data *dtpb.MedicalData) (*dtpb.MedicalData, error)
	GetMedicalData(ctx context.Context, id int64) (*dtpb.MedicalData, error)
	ListMedicalData(ctx context.Context, userID, clinicID *int64) ([]*dtpb.MedicalData, error)
	DeleteMedicalData(ctx context.Context, id int64) error

	CreateDataTransfer(ctx context.Context, transfer *dtpb.DataTransfer) (*dtpb.DataTransfer, error)
	GetDataTransfer(ctx context.Context, id int64) (*dtpb.DataTransfer, error)
	ListDataTransfers(ctx context.Context, clinicID, userID *int64, status *string) ([]*dtpb.DataTransfer, error)
	HandleDataTransferDecision(ctx context.Context, transferID, userID int64, confirm bool) (*dtpb.DataTransfer, error)

	CreateClinic(ctx context.Context, clinic *dtpb.Clinic) (*dtpb.Clinic, error)
	GetClinic(ctx context.Context, id int64) (*dtpb.Clinic, error)
	UpdateClinic(ctx context.Context, clinic *dtpb.Clinic) (*dtpb.Clinic, error)
	DeleteClinic(ctx context.Context, id int64) error
	ListClinics(ctx context.Context) ([]*dtpb.Clinic, error)
}

type dataTransferService struct {
	repo repository.DataTransferRepository
}

func NewDataTransferService(repo repository.DataTransferRepository) DataTransferService {
	return &dataTransferService{
		repo: repo,
	}
}

// --- MedicalData CRUD ---

func (s *dataTransferService) CreateMedicalData(ctx context.Context, data *dtpb.MedicalData) (*dtpb.MedicalData, error) {
	log := zerolog.Ctx(ctx).With().Str("component", "data_transfer_service").Str("method", "CreateMedicalData").Logger()

	dbModel := &model.MedicalData{}
	if err := copier.Copy(dbModel, data); err != nil {
		log.Error().Err(err).Msg("copy_failed_proto_to_model")
		return nil, err
	}

	if data.CreatedAt != "" {
		if t, err := time.Parse(time.RFC3339, data.CreatedAt); err == nil {
			dbModel.CreatedAt = t
		}
	}
	if data.UpdatedAt != "" {
		if t, err := time.Parse(time.RFC3339, data.UpdatedAt); err == nil {
			dbModel.UpdatedAt = t
		}
	}

	if err := s.repo.CreateMedicalData(ctx, dbModel); err != nil {
		log.Error().Err(err).Msg("repo_create_failed")
		return nil, err
	}

	var resp dtpb.MedicalData
	if err := copier.Copy(&resp, dbModel); err != nil {
		log.Error().Err(err).Msg("copy_failed_model_to_proto")
		return nil, err
	}

	resp.CreatedAt = dbModel.CreatedAt.Format(time.RFC3339)
	resp.UpdatedAt = dbModel.UpdatedAt.Format(time.RFC3339)
	return &resp, nil
}

func (s *dataTransferService) GetMedicalData(ctx context.Context, id int64) (*dtpb.MedicalData, error) {
	dbModel, err := s.repo.GetMedicalData(ctx, id)
	if err != nil {
		return nil, err
	}

	var resp dtpb.MedicalData
	if err := copier.Copy(&resp, dbModel); err != nil {
		return nil, err
	}

	resp.CreatedAt = dbModel.CreatedAt.Format(time.RFC3339)
	resp.UpdatedAt = dbModel.UpdatedAt.Format(time.RFC3339)
	return &resp, nil
}

func (s *dataTransferService) ListMedicalData(ctx context.Context, userID, clinicID *int64) ([]*dtpb.MedicalData, error) {
	dbModels, err := s.repo.ListMedicalData(ctx, userID, clinicID)
	if err != nil {
		return nil, err
	}

	var result []*dtpb.MedicalData
	for i := range dbModels {
		md := &dtpb.MedicalData{}
		if err := copier.Copy(md, &dbModels[i]); err != nil {
			return nil, err
		}
		md.CreatedAt = dbModels[i].CreatedAt.Format(time.RFC3339)
		md.UpdatedAt = dbModels[i].UpdatedAt.Format(time.RFC3339)
		result = append(result, md)
	}
	return result, nil
}

func (s *dataTransferService) DeleteMedicalData(ctx context.Context, id int64) error {
	return s.repo.DeleteMedicalData(ctx, id)
}

// --- DataTransfer CRUD ---

func (s *dataTransferService) CreateDataTransfer(ctx context.Context, transfer *dtpb.DataTransfer) (*dtpb.DataTransfer, error) {
	dbModel := &model.DataTransfer{}
	if err := copier.Copy(dbModel, transfer); err != nil {
		return nil, err
	}

	dbModel.Status = "pending"
	dbModel.CreatedAt = time.Now()
	dbModel.UpdatedAt = time.Now()

	if err := s.repo.CreateDataTransfer(ctx, dbModel); err != nil {
		return nil, err
	}

	var resp dtpb.DataTransfer
	copier.Copy(&resp, dbModel)
	resp.CreatedAt = dbModel.CreatedAt.Format(time.RFC3339)
	resp.UpdatedAt = dbModel.UpdatedAt.Format(time.RFC3339)
	return &resp, nil
}

func (s *dataTransferService) GetDataTransfer(ctx context.Context, id int64) (*dtpb.DataTransfer, error) {
	dbModel, err := s.repo.GetDataTransfer(ctx, id)
	if err != nil {
		return nil, err
	}

	var resp dtpb.DataTransfer
	copier.Copy(&resp, dbModel)
	resp.CreatedAt = dbModel.CreatedAt.Format(time.RFC3339)
	resp.UpdatedAt = dbModel.UpdatedAt.Format(time.RFC3339)
	return &resp, nil
}

func (s *dataTransferService) ListDataTransfers(ctx context.Context, clinicID, userID *int64, status *string) ([]*dtpb.DataTransfer, error) {
	dbModels, err := s.repo.ListDataTransfers(ctx, clinicID, userID, status)
	if err != nil {
		return nil, err
	}

	var result []*dtpb.DataTransfer
	for i := range dbModels {
		dt := &dtpb.DataTransfer{}
		copier.Copy(dt, &dbModels[i])
		dt.CreatedAt = dbModels[i].CreatedAt.Format(time.RFC3339)
		dt.UpdatedAt = dbModels[i].UpdatedAt.Format(time.RFC3339)
		result = append(result, dt)
	}

	return result, nil
}

func (s *dataTransferService) HandleDataTransferDecision(ctx context.Context, transferID, userID int64, confirm bool) (*dtpb.DataTransfer, error) {
	log := zerolog.Ctx(ctx).With().Str("component", "data_transfer_service").Str("method", "HandleDataTransferDecision").Logger()

	transfer, err := s.repo.GetDataTransfer(ctx, transferID)
	if err != nil {
		return nil, err
	}

	if transfer.UserID != userID {
		return nil, fmt.Errorf("user is not authorized to decide on this transfer")
	}

	if confirm {
		transfer.Status = "confirmed"

		// Получаем оригинальные медицинские данные
		originalData, err := s.repo.GetMedicalData(ctx, transfer.MedicalDataID)
		if err != nil {
			return nil, fmt.Errorf("failed to get medical data: %w", err)
		}

		// Создаём новый объект без ID, чтобы база сгенерировала уникальный ключ
		newMedicalData := &model.MedicalData{
			UserID:      originalData.UserID,
			ClinicID:    transfer.ToClinicID,
			Diagnosis:   originalData.Diagnosis,
			Complaint:   originalData.Complaint,
			Treatment:   originalData.Treatment,
			Medications: originalData.Medications,
			Allergies:   originalData.Allergies,
			DoctorNotes: originalData.DoctorNotes,
			LabResults:  originalData.LabResults,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := s.repo.CreateMedicalData(ctx, newMedicalData); err != nil {
			return nil, fmt.Errorf("failed to transfer medical data: %w", err)
		}

	} else {
		transfer.Status = "rejected"
	}

	transfer.UpdatedAt = time.Now()
	if err := s.repo.UpdateDataTransfer(ctx, transfer); err != nil {
		return nil, err
	}

	var resp dtpb.DataTransfer
	copier.Copy(&resp, transfer)
	resp.CreatedAt = transfer.CreatedAt.Format(time.RFC3339)
	resp.UpdatedAt = transfer.UpdatedAt.Format(time.RFC3339)

	log.Info().
		Int64("transfer_id", transferID).
		Str("status", transfer.Status).
		Msg("data_transfer_decision_handled")

	return &resp, nil
}

// --- Clinic CRUD ---

func (s *dataTransferService) CreateClinic(ctx context.Context, clinic *dtpb.Clinic) (*dtpb.Clinic, error) {
	log := zerolog.Ctx(ctx).With().Str("component", "data_transfer_service").Str("method", "CreateClinic").Logger()

	dbModel := &model.Clinic{}
	if err := copier.Copy(dbModel, clinic); err != nil {
		log.Error().Err(err).Msg("copy_failed_proto_to_model")
		return nil, err
	}

	// Ensure ID is 0 for new records (will be auto-generated by database)
	dbModel.ID = 0

	if err := s.repo.CreateClinic(ctx, dbModel); err != nil {
		log.Error().Err(err).Msg("repo_create_failed")
		return nil, err
	}

	var resp dtpb.Clinic
	if err := copier.Copy(&resp, dbModel); err != nil {
		log.Error().Err(err).Msg("copy_failed_model_to_proto")
		return nil, err
	}

	resp.CreatedAt = dbModel.CreatedAt.Format(time.RFC3339)
	resp.UpdatedAt = dbModel.UpdatedAt.Format(time.RFC3339)
	return &resp, nil
}

func (s *dataTransferService) GetClinic(ctx context.Context, id int64) (*dtpb.Clinic, error) {
	dbModel, err := s.repo.GetClinic(ctx, id)
	if err != nil {
		return nil, err
	}

	var resp dtpb.Clinic
	if err := copier.Copy(&resp, dbModel); err != nil {
		return nil, err
	}

	resp.CreatedAt = dbModel.CreatedAt.Format(time.RFC3339)
	resp.UpdatedAt = dbModel.UpdatedAt.Format(time.RFC3339)
	return &resp, nil
}

func (s *dataTransferService) UpdateClinic(ctx context.Context, clinic *dtpb.Clinic) (*dtpb.Clinic, error) {
	log := zerolog.Ctx(ctx).With().Str("component", "data_transfer_service").Str("method", "UpdateClinic").Logger()

	dbModel := &model.Clinic{}
	if err := copier.Copy(dbModel, clinic); err != nil {
		log.Error().Err(err).Msg("copy_failed_proto_to_model")
		return nil, err
	}

	if err := s.repo.UpdateClinic(ctx, dbModel); err != nil {
		log.Error().Err(err).Msg("repo_update_failed")
		return nil, err
	}

	var resp dtpb.Clinic
	if err := copier.Copy(&resp, dbModel); err != nil {
		log.Error().Err(err).Msg("copy_failed_model_to_proto")
		return nil, err
	}

	resp.CreatedAt = dbModel.CreatedAt.Format(time.RFC3339)
	resp.UpdatedAt = dbModel.UpdatedAt.Format(time.RFC3339)
	return &resp, nil
}

func (s *dataTransferService) DeleteClinic(ctx context.Context, id int64) error {
	return s.repo.DeleteClinic(ctx, id)
}

func (s *dataTransferService) ListClinics(ctx context.Context) ([]*dtpb.Clinic, error) {
	dbModels, err := s.repo.ListClinics(ctx)
	if err != nil {
		return nil, err
	}

	var result []*dtpb.Clinic
	for i := range dbModels {
		c := &dtpb.Clinic{}
		if err := copier.Copy(c, &dbModels[i]); err != nil {
			return nil, err
		}
		c.CreatedAt = dbModels[i].CreatedAt.Format(time.RFC3339)
		c.UpdatedAt = dbModels[i].UpdatedAt.Format(time.RFC3339)
		result = append(result, c)
	}
	return result, nil
}
