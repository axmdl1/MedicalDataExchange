package repository

import (
	"errors"

	"gorm.io/gorm"
)

type MedicalData struct {
	ID        string `gorm:"column:id"`
	UserID    int64
	ClinicID  int64
	Diagnosis string
	Complaint string
	Treatment string
}

type MedicalRepository struct {
	DB *gorm.DB
}

func NewMedicalRepository(db *gorm.DB) MedicalRepository {
	return MedicalRepository{DB: db}
}

func (r MedicalRepository) GetMedicalData(recordID string) (*MedicalData, error) {
	var data MedicalData
	err := r.DB.Where("id = ?", recordID).First(&data).Error
	if err != nil {
		return nil, errors.New("record not found")
	}
	return &data, nil
}
