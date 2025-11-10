package model

import "time"

// MedicalData представляет мед. данные пациента
type MedicalData struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	UserID      int64     `gorm:"column:user_id"`
	ClinicID    int64     `gorm:"column:clinic_id"`
	Diagnosis   string    `gorm:"column:diagnosis"`
	Complaint   string    `gorm:"column:complaint"`
	Treatment   string    `gorm:"column:treatment"`
	Medications string    `gorm:"column:medications"`
	Allergies   string    `gorm:"column:allergies"`
	DoctorNotes string    `gorm:"column:doctor_notes"`
	LabResults  string    `gorm:"column:lab_results"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (MedicalData) TableName() string {
	return "medical_data"
}

// DataTransfer — запрос на передачу мед. данных между клиниками
type DataTransfer struct {
	ID            int64     `gorm:"column:id;primaryKey"`
	MedicalDataID int64     `gorm:"column:medical_data_id"`
	FromClinicID  int64     `gorm:"column:from_clinic_id"`
	ToClinicID    int64     `gorm:"column:to_clinic_id"`
	UserID        int64     `gorm:"column:user_id"`
	Status        string    `gorm:"column:status"` // pending, confirmed, rejected
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (DataTransfer) TableName() string {
	return "data_transfers"
}
