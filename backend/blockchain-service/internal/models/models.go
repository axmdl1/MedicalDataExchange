package models

import "time"

// AccessRequestInput для создания запроса на доступ
type AccessRequestInput struct {
	ID            string
	PatientID     string
	ClinicID      string
	MedicalDataID string
	RequestType   string // patient_view, clinic_transfer
	FromClinicID  string
	ToClinicID    string
}

// AccessRequest представляет запрос на доступ к медицинским данным
type AccessRequest struct {
	ID            string     `json:"id"`
	PatientID     string     `json:"patient_id"`
	ClinicID      string     `json:"clinic_id"`
	MedicalDataID string     `json:"medical_data_id"`
	RequestType   string     `json:"request_type"`
	Status        string     `json:"status"`
	RequestedAt   time.Time  `json:"requested_at"`
	ApprovedAt    *time.Time `json:"approved_at,omitempty"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	FromClinicID  string     `json:"from_clinic_id,omitempty"`
	ToClinicID    string     `json:"to_clinic_id,omitempty"`
}

// TemporaryAccess представляет временный доступ к данным
type TemporaryAccess struct {
	ID            string    `json:"id"`
	AccessToken   string    `json:"access_token"`
	PatientID     string    `json:"patient_id"`
	ClinicID      string    `json:"clinic_id"`
	MedicalDataID string    `json:"medical_data_id"`
	GrantedAt     time.Time `json:"granted_at"`
	ExpiresAt     time.Time `json:"expires_at"`
	IsRevoked     bool      `json:"is_revoked"`
}

// DataTransferInput для логирования передачи данных
type DataTransferInput struct {
	ID            string
	MedicalDataID string
	PatientID     string
	FromClinicID  string
	ToClinicID    string
	ApprovedBy    string
	DataHash      string
}

// DataTransferLog представляет запись о передаче данных между клиниками
type DataTransferLog struct {
	ID            string    `json:"id"`
	MedicalDataID string    `json:"medical_data_id"`
	PatientID     string    `json:"patient_id"`
	FromClinicID  string    `json:"from_clinic_id"`
	ToClinicID    string    `json:"to_clinic_id"`
	ApprovedBy    string    `json:"approved_by"`
	TransferredAt time.Time `json:"transferred_at"`
	DataHash      string    `json:"data_hash"`
}
