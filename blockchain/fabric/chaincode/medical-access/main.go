package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// MedicalAccessContract управляет доступом к медицинским данным
type MedicalAccessContract struct {
	contractapi.Contract
}

// AccessRequest представляет запрос на доступ к медицинским данным
type AccessRequest struct {
	ID           string    `json:"id"`
	PatientID    string    `json:"patient_id"`
	ClinicID     string    `json:"clinic_id"`
	MedicalDataID string   `json:"medical_data_id"`
	RequestType  string    `json:"request_type"` // patient_view, clinic_transfer
	Status       string    `json:"status"`       // pending, approved, rejected, expired
	RequestedAt  time.Time `json:"requested_at"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	FromClinicID string    `json:"from_clinic_id,omitempty"`
	ToClinicID   string    `json:"to_clinic_id,omitempty"`
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

// DataTransferLog представляет запись о передаче данных между клиниками
type DataTransferLog struct {
	ID            string    `json:"id"`
	MedicalDataID string    `json:"medical_data_id"`
	PatientID     string    `json:"patient_id"`
	FromClinicID  string    `json:"from_clinic_id"`
	ToClinicID    string    `json:"to_clinic_id"`
	ApprovedBy    string    `json:"approved_by"` // patient ID
	TransferredAt time.Time `json:"transferred_at"`
	DataHash      string    `json:"data_hash"` // Hash данных для аудита
}

// InitLedger инициализирует леджер
func (c *MedicalAccessContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

// CreateAccessRequest создает запрос на доступ к медицинским данным
func (c *MedicalAccessContract) CreateAccessRequest(
	ctx contractapi.TransactionContextInterface,
	id string,
	patientID string,
	clinicID string,
	medicalDataID string,
	requestType string,
	fromClinicID string,
	toClinicID string,
) error {
	request := AccessRequest{
		ID:           id,
		PatientID:    patientID,
		ClinicID:     clinicID,
		MedicalDataID: medicalDataID,
		RequestType:  requestType,
		Status:       "pending",
		RequestedAt:  time.Now(),
		FromClinicID: fromClinicID,
		ToClinicID:   toClinicID,
	}

	requestJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(fmt.Sprintf("REQUEST_%s", id), requestJSON)
}

// ApproveAccessRequest одобряет запрос на доступ (вызывается клиникой)
func (c *MedicalAccessContract) ApproveAccessRequest(
	ctx contractapi.TransactionContextInterface,
	requestID string,
	approverID string,
) error {
	requestJSON, err := ctx.GetStub().GetState(fmt.Sprintf("REQUEST_%s", requestID))
	if err != nil {
		return fmt.Errorf("failed to read request: %v", err)
	}
	if requestJSON == nil {
		return fmt.Errorf("request %s does not exist", requestID)
	}

	var request AccessRequest
	err = json.Unmarshal(requestJSON, &request)
	if err != nil {
		return err
	}

	if request.Status != "pending" {
		return fmt.Errorf("request is not in pending status")
	}

	now := time.Now()
	expiresAt := now.Add(15 * time.Minute) // Доступ на 15 минут

	request.Status = "approved"
	request.ApprovedAt = &now
	request.ExpiresAt = &expiresAt

	updatedRequestJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(fmt.Sprintf("REQUEST_%s", requestID), updatedRequestJSON)
	if err != nil {
		return err
	}

	// Создаем временный доступ
	if request.RequestType == "patient_view" {
		accessToken := fmt.Sprintf("TOKEN_%s_%d", requestID, time.Now().Unix())

		temporaryAccess := TemporaryAccess{
			ID:            fmt.Sprintf("ACCESS_%s", requestID),
			AccessToken:   accessToken,
			PatientID:     request.PatientID,
			ClinicID:      request.ClinicID,
			MedicalDataID: request.MedicalDataID,
			GrantedAt:     now,
			ExpiresAt:     expiresAt,
			IsRevoked:     false,
		}

		accessJSON, err := json.Marshal(temporaryAccess)
		if err != nil {
			return err
		}

		return ctx.GetStub().PutState(fmt.Sprintf("ACCESS_%s", requestID), accessJSON)
	}

	return nil
}

// RejectAccessRequest отклоняет запрос на доступ
func (c *MedicalAccessContract) RejectAccessRequest(
	ctx contractapi.TransactionContextInterface,
	requestID string,
) error {
	requestJSON, err := ctx.GetStub().GetState(fmt.Sprintf("REQUEST_%s", requestID))
	if err != nil {
		return fmt.Errorf("failed to read request: %v", err)
	}
	if requestJSON == nil {
		return fmt.Errorf("request %s does not exist", requestID)
	}

	var request AccessRequest
	err = json.Unmarshal(requestJSON, &request)
	if err != nil {
		return err
	}

	if request.Status != "pending" {
		return fmt.Errorf("request is not in pending status")
	}

	request.Status = "rejected"

	updatedRequestJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(fmt.Sprintf("REQUEST_%s", requestID), updatedRequestJSON)
}

// ValidateAccess проверяет действительность токена доступа
func (c *MedicalAccessContract) ValidateAccess(
	ctx contractapi.TransactionContextInterface,
	accessID string,
) (*TemporaryAccess, error) {
	accessJSON, err := ctx.GetStub().GetState(fmt.Sprintf("ACCESS_%s", accessID))
	if err != nil {
		return nil, fmt.Errorf("failed to read access: %v", err)
	}
	if accessJSON == nil {
		return nil, fmt.Errorf("access %s does not exist", accessID)
	}

	var access TemporaryAccess
	err = json.Unmarshal(accessJSON, &access)
	if err != nil {
		return nil, err
	}

	// Проверяем не истек ли доступ
	if time.Now().After(access.ExpiresAt) {
		access.IsRevoked = true
		// Обновляем статус в леджере
		updatedJSON, _ := json.Marshal(access)
		ctx.GetStub().PutState(fmt.Sprintf("ACCESS_%s", accessID), updatedJSON)
		return nil, fmt.Errorf("access token has expired")
	}

	if access.IsRevoked {
		return nil, fmt.Errorf("access token has been revoked")
	}

	return &access, nil
}

// RevokeAccess отзывает доступ до истечения времени
func (c *MedicalAccessContract) RevokeAccess(
	ctx contractapi.TransactionContextInterface,
	accessID string,
) error {
	accessJSON, err := ctx.GetStub().GetState(fmt.Sprintf("ACCESS_%s", accessID))
	if err != nil {
		return fmt.Errorf("failed to read access: %v", err)
	}
	if accessJSON == nil {
		return fmt.Errorf("access %s does not exist", accessID)
	}

	var access TemporaryAccess
	err = json.Unmarshal(accessJSON, &access)
	if err != nil {
		return err
	}

	access.IsRevoked = true

	updatedJSON, err := json.Marshal(access)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(fmt.Sprintf("ACCESS_%s", accessID), updatedJSON)
}

// LogDataTransfer записывает передачу данных между клиниками
func (c *MedicalAccessContract) LogDataTransfer(
	ctx contractapi.TransactionContextInterface,
	id string,
	medicalDataID string,
	patientID string,
	fromClinicID string,
	toClinicID string,
	approvedBy string,
	dataHash string,
) error {
	log := DataTransferLog{
		ID:            id,
		MedicalDataID: medicalDataID,
		PatientID:     patientID,
		FromClinicID:  fromClinicID,
		ToClinicID:    toClinicID,
		ApprovedBy:    approvedBy,
		TransferredAt: time.Now(),
		DataHash:      dataHash,
	}

	logJSON, err := json.Marshal(log)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(fmt.Sprintf("TRANSFER_%s", id), logJSON)
}

// GetAccessRequest получает запрос на доступ
func (c *MedicalAccessContract) GetAccessRequest(
	ctx contractapi.TransactionContextInterface,
	requestID string,
) (*AccessRequest, error) {
	requestJSON, err := ctx.GetStub().GetState(fmt.Sprintf("REQUEST_%s", requestID))
	if err != nil {
		return nil, fmt.Errorf("failed to read request: %v", err)
	}
	if requestJSON == nil {
		return nil, fmt.Errorf("request %s does not exist", requestID)
	}

	var request AccessRequest
	err = json.Unmarshal(requestJSON, &request)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

// GetTransferLog получает лог передачи данных
func (c *MedicalAccessContract) GetTransferLog(
	ctx contractapi.TransactionContextInterface,
	transferID string,
) (*DataTransferLog, error) {
	logJSON, err := ctx.GetStub().GetState(fmt.Sprintf("TRANSFER_%s", transferID))
	if err != nil {
		return nil, fmt.Errorf("failed to read transfer log: %v", err)
	}
	if logJSON == nil {
		return nil, fmt.Errorf("transfer log %s does not exist", transferID)
	}

	var log DataTransferLog
	err = json.Unmarshal(logJSON, &log)
	if err != nil {
		return nil, err
	}

	return &log, nil
}

// GetPatientAccessHistory получает историю доступа пациента к данным
func (c *MedicalAccessContract) GetPatientAccessHistory(
	ctx contractapi.TransactionContextInterface,
	patientID string,
) ([]*AccessRequest, error) {
	// В реальной реализации использовать CouchDB rich queries
	// Для упрощения возвращаем пустой массив
	return []*AccessRequest{}, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&MedicalAccessContract{})
	if err != nil {
		fmt.Printf("Error creating medical access chaincode: %v\n", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting medical access chaincode: %v\n", err)
	}
}