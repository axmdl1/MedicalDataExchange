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
	return "data_transfer"
}

// Clinic — клиника
type Clinic struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Address   string    `gorm:"column:address"`
	Phone     string    `gorm:"column:phone"`
	Email     string    `gorm:"column:email"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Clinic) TableName() string {
	return "clinics"
}

// PatientAccessRequest представляет запрос пациента на просмотр своих данных
type PatientAccessRequest struct {
	ID              int64     `gorm:"column:id;primaryKey"`
	PatientID       int64     `gorm:"column:patient_id"`
	ClinicID        int64     `gorm:"column:clinic_id"`
	MedicalDataID   int64     `gorm:"column:medical_data_id"`
	Status          string    `gorm:"column:status"` // pending, approved, rejected, expired
	BlockchainTxID  string    `gorm:"column:blockchain_tx_id"`
	RequestedAt     time.Time `gorm:"column:requested_at;autoCreateTime"`
	ApprovedAt      *time.Time `gorm:"column:approved_at"`
	ExpiresAt       *time.Time `gorm:"column:expires_at"`
}

func (PatientAccessRequest) TableName() string {
	return "patient_access_requests"
}

// TemporaryPatientData представляет временную копию мед. данных для пациента
type TemporaryPatientData struct {
	ID            int64     `gorm:"column:id;primaryKey"`
	AccessToken   string    `gorm:"column:access_token;uniqueIndex"`
	PatientID     int64     `gorm:"column:patient_id"`
	ClinicID      int64     `gorm:"column:clinic_id"`
	MedicalDataID int64     `gorm:"column:medical_data_id"`
	// Копия данных (зашифрованных)
	EncryptedData string    `gorm:"column:encrypted_data;type:text"`
	GrantedAt     time.Time `gorm:"column:granted_at;autoCreateTime"`
	ExpiresAt     time.Time `gorm:"column:expires_at"`
	IsRevoked     bool      `gorm:"column:is_revoked;default:false"`
}

func (TemporaryPatientData) TableName() string {
	return "temporary_patient_data"
}
