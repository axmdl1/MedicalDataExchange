package model

import (
	"database/sql"
	"time"
)

// MedicalData представляет полноценные медицинские данные пациента
type MedicalData struct {
	ID       int64 `gorm:"column:id;primaryKey"`
	UserID   int64 `gorm:"column:user_id"`
	ClinicID int64 `gorm:"column:clinic_id"`

	// Visit information
	VisitDate       time.Time      `gorm:"column:visit_date"`
	VisitType       sql.NullString `gorm:"column:visit_type"`
	Department      sql.NullString `gorm:"column:department"`
	AttendingDoctor sql.NullString `gorm:"column:attending_doctor"`

	// Patient complaints and symptoms
	ChiefComplaint   sql.NullString `gorm:"column:chief_complaint"`
	Symptoms         sql.NullString `gorm:"column:symptoms"`
	SymptomDuration  sql.NullString `gorm:"column:symptom_duration"`
	PainLevel        sql.NullInt64  `gorm:"column:pain_level"`

	// Vital signs
	Temperature            sql.NullFloat64 `gorm:"column:temperature"`
	BloodPressureSystolic  sql.NullInt64   `gorm:"column:blood_pressure_systolic"`
	BloodPressureDiastolic sql.NullInt64   `gorm:"column:blood_pressure_diastolic"`
	HeartRate              sql.NullInt64   `gorm:"column:heart_rate"`
	RespiratoryRate        sql.NullInt64   `gorm:"column:respiratory_rate"`
	OxygenSaturation       sql.NullInt64   `gorm:"column:oxygen_saturation"`
	Weight                 sql.NullFloat64 `gorm:"column:weight"`
	Height                 sql.NullFloat64 `gorm:"column:height"`
	BMI                    sql.NullFloat64 `gorm:"column:bmi"`

	// Medical assessment
	Diagnosis            string         `gorm:"column:diagnosis"`
	DiagnosisCode        sql.NullString `gorm:"column:diagnosis_code"`
	SecondaryDiagnoses   sql.NullString `gorm:"column:secondary_diagnoses"`
	Severity             sql.NullString `gorm:"column:severity"`

	// Medical history
	MedicalHistory      sql.NullString `gorm:"column:medical_history"`
	SurgicalHistory     sql.NullString `gorm:"column:surgical_history"`
	FamilyHistory       sql.NullString `gorm:"column:family_history"`
	Allergies           sql.NullString `gorm:"column:allergies"`
	CurrentMedications  sql.NullString `gorm:"column:current_medications"`

	// Treatment and care plan
	TreatmentPlan         string         `gorm:"column:treatment_plan"`
	PrescribedMedications sql.NullString `gorm:"column:prescribed_medications"`
	ProceduresPerformed   sql.NullString `gorm:"column:procedures_performed"`
	LabTestsOrdered       sql.NullString `gorm:"column:lab_tests_ordered"`
	ImagingOrdered        sql.NullString `gorm:"column:imaging_ordered"`
	Referrals             sql.NullString `gorm:"column:referrals"`

	// Lab results
	LabResults        sql.NullString `gorm:"column:lab_results;type:jsonb"`
	LabResultsSummary sql.NullString `gorm:"column:lab_results_summary"`

	// Doctor's notes
	DoctorNotes            sql.NullString `gorm:"column:doctor_notes"`
	FollowUpInstructions   sql.NullString `gorm:"column:follow_up_instructions"`
	FollowUpDate           sql.NullTime   `gorm:"column:follow_up_date"`
	Restrictions           sql.NullString `gorm:"column:restrictions"`

	// Administrative
	InsuranceInfo sql.NullString  `gorm:"column:insurance_info"`
	BillingCode   sql.NullString  `gorm:"column:billing_code"`
	EstimatedCost sql.NullFloat64 `gorm:"column:estimated_cost"`

	// Metadata
	RecordStatus         sql.NullString `gorm:"column:record_status"`
	ConfidentialityLevel sql.NullString `gorm:"column:confidentiality_level"`
	DataHash             sql.NullString `gorm:"column:data_hash"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
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
