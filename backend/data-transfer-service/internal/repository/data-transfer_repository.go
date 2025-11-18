package repository

import (
	"context"
	"time"

	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/model"
	"gorm.io/gorm"
)

type DataTransferRepository interface {
	// MedicalData CRUD
	CreateMedicalData(ctx context.Context, data *model.MedicalData) error
	GetMedicalData(ctx context.Context, id int64) (*model.MedicalData, error)
	ListMedicalData(ctx context.Context, userID, clinicID *int64) ([]model.MedicalData, error)
	DeleteMedicalData(ctx context.Context, id int64) error

	// DataTransfer CRUD
	CreateDataTransfer(ctx context.Context, transfer *model.DataTransfer) error
	GetDataTransfer(ctx context.Context, id int64) (*model.DataTransfer, error)
	ListDataTransfers(ctx context.Context, clinicID, userID *int64, status *string) ([]model.DataTransfer, error)
	UpdateDataTransfer(ctx context.Context, transfer *model.DataTransfer) error

	// Clinic CRUD
	CreateClinic(ctx context.Context, clinic *model.Clinic) error
	GetClinic(ctx context.Context, id int64) (*model.Clinic, error)
	UpdateClinic(ctx context.Context, clinic *model.Clinic) error
	DeleteClinic(ctx context.Context, id int64) error
	ListClinics(ctx context.Context) ([]model.Clinic, error)

	// PatientAccessRequest CRUD
	CreatePatientAccessRequest(ctx context.Context, request *model.PatientAccessRequest) error
	GetPatientAccessRequest(ctx context.Context, id int64) (*model.PatientAccessRequest, error)
	UpdatePatientAccessRequest(ctx context.Context, request *model.PatientAccessRequest) error
	ListPatientAccessRequests(ctx context.Context, patientID, clinicID *int64, status *string) ([]*model.PatientAccessRequest, error)
	UpdateExpiredAccessRequests(ctx context.Context) (int, error)

	// TemporaryPatientData CRUD
	CreateTemporaryPatientData(ctx context.Context, data *model.TemporaryPatientData) error
	GetTemporaryDataByToken(ctx context.Context, accessToken string) (*model.TemporaryPatientData, error)
	UpdateTemporaryPatientData(ctx context.Context, data *model.TemporaryPatientData) error
	DeleteTemporaryPatientData(ctx context.Context, id int64) error
	DeleteExpiredTemporaryData(ctx context.Context) (int, error)
}

type dataTransferRepo struct {
	db *gorm.DB
}

func NewDataTransferRepository(db *gorm.DB) DataTransferRepository {
	return &dataTransferRepo{db: db}
}

// --- MedicalData ---
func (r *dataTransferRepo) CreateMedicalData(ctx context.Context, data *model.MedicalData) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *dataTransferRepo) GetMedicalData(ctx context.Context, id int64) (*model.MedicalData, error) {
	var data model.MedicalData
	if err := r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *dataTransferRepo) ListMedicalData(ctx context.Context, userID, clinicID *int64) ([]model.MedicalData, error) {
	var data []model.MedicalData
	query := r.db.WithContext(ctx).Model(&model.MedicalData{})
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if clinicID != nil {
		query = query.Where("clinic_id = ?", *clinicID)
	}
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *dataTransferRepo) DeleteMedicalData(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.MedicalData{}, id).Error
}

// --- DataTransfer ---
func (r *dataTransferRepo) CreateDataTransfer(ctx context.Context, transfer *model.DataTransfer) error {
	transfer.Status = "pending"
	transfer.CreatedAt = time.Now()
	transfer.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(transfer).Error
}

func (r *dataTransferRepo) GetDataTransfer(ctx context.Context, id int64) (*model.DataTransfer, error) {
	var transfer model.DataTransfer
	if err := r.db.WithContext(ctx).First(&transfer, id).Error; err != nil {
		return nil, err
	}
	return &transfer, nil
}

func (r *dataTransferRepo) ListDataTransfers(ctx context.Context, clinicID, userID *int64, status *string) ([]model.DataTransfer, error) {
	var transfers []model.DataTransfer
	query := r.db.WithContext(ctx).Model(&model.DataTransfer{})
	if clinicID != nil {
		query = query.Where("to_clinic_id = ? OR from_clinic_id = ?", *clinicID, *clinicID)
	}
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Find(&transfers).Error; err != nil {
		return nil, err
	}
	return transfers, nil
}

func (r *dataTransferRepo) UpdateDataTransfer(ctx context.Context, transfer *model.DataTransfer) error {
	transfer.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(transfer).Error
}

// --- Clinic ---
func (r *dataTransferRepo) CreateClinic(ctx context.Context, clinic *model.Clinic) error {
	return r.db.WithContext(ctx).Create(clinic).Error
}

func (r *dataTransferRepo) GetClinic(ctx context.Context, id int64) (*model.Clinic, error) {
	var clinic model.Clinic
	if err := r.db.WithContext(ctx).First(&clinic, id).Error; err != nil {
		return nil, err
	}
	return &clinic, nil
}

func (r *dataTransferRepo) UpdateClinic(ctx context.Context, clinic *model.Clinic) error {
	return r.db.WithContext(ctx).Save(clinic).Error
}

func (r *dataTransferRepo) DeleteClinic(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Clinic{}, id).Error
}

func (r *dataTransferRepo) ListClinics(ctx context.Context) ([]model.Clinic, error) {
	var clinics []model.Clinic
	if err := r.db.WithContext(ctx).Find(&clinics).Error; err != nil {
		return nil, err
	}
	return clinics, nil
}

// --- PatientAccessRequest ---
func (r *dataTransferRepo) CreatePatientAccessRequest(ctx context.Context, request *model.PatientAccessRequest) error {
	return r.db.WithContext(ctx).Create(request).Error
}

func (r *dataTransferRepo) GetPatientAccessRequest(ctx context.Context, id int64) (*model.PatientAccessRequest, error) {
	var request model.PatientAccessRequest
	if err := r.db.WithContext(ctx).First(&request, id).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *dataTransferRepo) UpdatePatientAccessRequest(ctx context.Context, request *model.PatientAccessRequest) error {
	return r.db.WithContext(ctx).Save(request).Error
}

func (r *dataTransferRepo) ListPatientAccessRequests(ctx context.Context, patientID, clinicID *int64, status *string) ([]*model.PatientAccessRequest, error) {
	var requests []*model.PatientAccessRequest
	query := r.db.WithContext(ctx).Model(&model.PatientAccessRequest{})

	if patientID != nil {
		query = query.Where("patient_id = ?", *patientID)
	}
	if clinicID != nil {
		query = query.Where("clinic_id = ?", *clinicID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Order("requested_at DESC").Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *dataTransferRepo) UpdateExpiredAccessRequests(ctx context.Context) (int, error) {
	result := r.db.WithContext(ctx).
		Model(&model.PatientAccessRequest{}).
		Where("expires_at < ? AND status = ?", time.Now(), "approved").
		Update("status", "expired")

	return int(result.RowsAffected), result.Error
}

// --- TemporaryPatientData ---
func (r *dataTransferRepo) CreateTemporaryPatientData(ctx context.Context, data *model.TemporaryPatientData) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *dataTransferRepo) GetTemporaryDataByToken(ctx context.Context, accessToken string) (*model.TemporaryPatientData, error) {
	var data model.TemporaryPatientData
	if err := r.db.WithContext(ctx).
		Where("access_token = ?", accessToken).
		First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *dataTransferRepo) UpdateTemporaryPatientData(ctx context.Context, data *model.TemporaryPatientData) error {
	return r.db.WithContext(ctx).Save(data).Error
}

func (r *dataTransferRepo) DeleteTemporaryPatientData(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.TemporaryPatientData{}, id).Error
}

func (r *dataTransferRepo) DeleteExpiredTemporaryData(ctx context.Context) (int, error) {
	result := r.db.WithContext(ctx).
		Where("expires_at < ? OR is_revoked = ?", time.Now(), true).
		Delete(&model.TemporaryPatientData{})

	return int(result.RowsAffected), result.Error
}
