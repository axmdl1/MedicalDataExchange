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
