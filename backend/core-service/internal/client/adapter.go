package client

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	dtpb "github.com/axmdl1/MedicalDataExchange/core-service/pkg/gen/go/data-transfer"
)

// InMemoryAccessRequest — структура доступа, хранимая в адаптере для demo.
type InMemoryAccessRequest struct {
	RequestID string
	RecordID  string
	Patient   string
	Clinic    string
	Status    string // pending / approved / rejected / consumed / revoked
	Reason    string // optional (for reject)
	// you can add TokenHash, ExpiresAt etc. if you want to integrate with on-chain
}

// DataTransferAdapter — адаптер между handler'ами и реальным gRPC-клиентом.
// Он использует dtpb.* RPC методы для работы с медицинскими данными и
// локальный in-memory store для access requests (для демонстрации).
type DataTransferAdapter struct {
	Inner *DataTransferClient

	// accessRequests — in-memory store для access requests (demo)
	mu             sync.RWMutex
	accessRequests map[string]*InMemoryAccessRequest
}

// NewDataTransferAdapter creates an adapter wrapping the real DataTransferClient.
func NewDataTransferAdapter(inner *DataTransferClient) *DataTransferAdapter {
	return &DataTransferAdapter{
		Inner:          inner,
		accessRequests: make(map[string]*InMemoryAccessRequest),
	}
}

// -----------------------------
// Access Request methods (in-memory demo)
// -----------------------------

// CreateAccessRequest - сохраняет request в in-memory store.
// requestID — строка, recordID — строка (если у вас в dtpb Id — int64, то используйте число в recordID)
func (a *DataTransferAdapter) CreateAccessRequest(ctx context.Context, requestID, recordID, patient, clinic string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.accessRequests[requestID]; ok {
		return fmt.Errorf("request %s already exists", requestID)
	}
	a.accessRequests[requestID] = &InMemoryAccessRequest{
		RequestID: requestID,
		RecordID:  recordID,
		Patient:   patient,
		Clinic:    clinic,
		Status:    "pending",
	}
	return nil
}

// UpdateRequestStatus - обновляет статус существующего запроса
func (a *DataTransferAdapter) UpdateRequestStatus(ctx context.Context, requestID, status string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	ar, ok := a.accessRequests[requestID]
	if !ok {
		return fmt.Errorf("request %s not found", requestID)
	}
	ar.Status = status
	return nil
}

// GetAccessRequest - возвращает доступный запрос
func (a *DataTransferAdapter) GetAccessRequest(ctx context.Context, requestID string) (interface{}, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	ar, ok := a.accessRequests[requestID]
	if !ok {
		return nil, fmt.Errorf("request %s not found", requestID)
	}
	// можно возвращать конкретный тип, но handler ожидает interface{}.
	return ar, nil
}

// ListAccessRequests - возвращает все запросы (или примените фильтр по clinic)
func (a *DataTransferAdapter) ListAccessRequests(ctx context.Context) ([]interface{}, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	out := make([]interface{}, 0, len(a.accessRequests))
	for _, v := range a.accessRequests {
		out = append(out, v)
	}
	return out, nil
}

// RejectAccessRequest - отмечает запрос как rejected + optional reason
func (a *DataTransferAdapter) RejectAccessRequest(ctx context.Context, requestID, reason string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	ar, ok := a.accessRequests[requestID]
	if !ok {
		return fmt.Errorf("request %s not found", requestID)
	}
	ar.Status = "rejected"
	ar.Reason = reason
	return nil
}

// RevokeAccess - optional helper for revoking access (used from handlers)
func (a *DataTransferAdapter) RevokeAccess(ctx context.Context, requestID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	ar, ok := a.accessRequests[requestID]
	if !ok {
		return fmt.Errorf("request %s not found", requestID)
	}
	ar.Status = "revoked"
	return nil
}

// -----------------------------
// MedicalData / DataTransfer methods (actual gRPC calls to dtpb)
// -----------------------------

// GetMedicalData uses dtpb.GetMedicalData RPC. recordID is expected to be numeric string (id).
// If recordID is not a number we will try to return an error to make intent explicit.
func (a *DataTransferAdapter) GetMedicalData(ctx context.Context, recordID string) (interface{}, error) {
	if a.Inner == nil || a.Inner.Api == nil {
		return nil, errors.New("data-transfer client not configured")
	}

	// try parse numeric id
	id, err := strconv.ParseInt(recordID, 10, 64)
	if err != nil {
		// recordID not numeric — we could extend to search by other fields, but keep explicit for demo
		return nil, fmt.Errorf("GetMedicalData: recordID must be numeric id string, got %q", recordID)
	}

	resp, err := a.Inner.Api.GetMedicalData(ctx, &dtpb.GetMedicalDataRequest{Id: id})
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.MedicalData == nil {
		return nil, fmt.Errorf("medical data %d not found", id)
	}
	// return the protobuf message pointer (handler can convert to JSON)
	return resp.MedicalData, nil
}

// CreateMedicalData calls dtpb.CreateMedicalData (if such exists).
// In your proto we have CreateMedicalDataRequest/Response, so adapt accordingly.
func (a *DataTransferAdapter) CreateMedicalData(ctx context.Context, md *dtpb.MedicalData) (*dtpb.MedicalData, error) {
	if a.Inner == nil || a.Inner.Api == nil {
		return nil, errors.New("data-transfer client not configured")
	}
	req := &dtpb.CreateMedicalDataRequest{
		MedicalData: md,
	}
	resp, err := a.Inner.Api.CreateMedicalData(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("CreateMedicalData: empty response")
	}
	return resp.MedicalData, nil
}

// DeleteMedicalData calls dtpb.DeleteMedicalData.
func (a *DataTransferAdapter) DeleteMedicalData(ctx context.Context, id int64) (bool, error) {
	if a.Inner == nil || a.Inner.Api == nil {
		return false, errors.New("data-transfer client not configured")
	}
	resp, err := a.Inner.Api.DeleteMedicalData(ctx, &dtpb.DeleteMedicalDataRequest{Id: id})
	if err != nil {
		return false, err
	}
	if resp == nil {
		return false, errors.New("DeleteMedicalData: empty response")
	}
	return resp.Success, nil
}

// ListMedicalData wraps dtpb.ListMedicalData.
func (a *DataTransferAdapter) ListMedicalData(ctx context.Context, userId *int64, clinicId *int64) ([]*dtpb.MedicalData, error) {
	if a.Inner == nil || a.Inner.Api == nil {
		return nil, errors.New("data-transfer client not configured")
	}
	req := &dtpb.ListMedicalDataRequest{}
	if userId != nil {
		req.UserId = userId
	}
	if clinicId != nil {
		req.ClinicId = clinicId
	}
	resp, err := a.Inner.Api.ListMedicalData(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("ListMedicalData: empty response")
	}
	return resp.MedicalData, nil
}

// CreateDataTransfer wraps dtpb.CreateDataTransfer (to create transfer records between clinics)
func (a *DataTransferAdapter) CreateDataTransfer(ctx context.Context, medicalDataId, fromClinicId, toClinicId, userId int64) (*dtpb.DataTransfer, error) {
	if a.Inner == nil || a.Inner.Api == nil {
		return nil, errors.New("data-transfer client not configured")
	}
	req := &dtpb.CreateDataTransferRequest{
		MedicalDataId: medicalDataId,
		FromClinicId:  fromClinicId,
		ToClinicId:    toClinicId,
		UserId:        userId,
	}
	resp, err := a.Inner.Api.CreateDataTransfer(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("CreateDataTransfer: empty response")
	}
	return resp.Transfer, nil
}

// GetDataTransfer wraps dtpb.GetDataTransfer.
func (a *DataTransferAdapter) GetDataTransfer(ctx context.Context, id int64) (*dtpb.DataTransfer, error) {
	if a.Inner == nil || a.Inner.Api == nil {
		return nil, errors.New("data-transfer client not configured")
	}
	resp, err := a.Inner.Api.GetDataTransfer(ctx, &dtpb.GetDataTransferRequest{Id: id})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("GetDataTransfer: empty response")
	}
	return resp.Transfer, nil
}

// ListDataTransfers wraps dtpb.ListDataTransfers.
func (a *DataTransferAdapter) ListDataTransfers(ctx context.Context, clinicId *int64, userId *int64, status *string) ([]*dtpb.DataTransfer, error) {
	if a.Inner == nil || a.Inner.Api == nil {
		return nil, errors.New("data-transfer client not configured")
	}
	req := &dtpb.ListDataTransfersRequest{}
	if clinicId != nil {
		req.ClinicId = clinicId
	}
	if userId != nil {
		req.UserId = userId
	}
	if status != nil {
		req.Status = status
	}
	resp, err := a.Inner.Api.ListDataTransfers(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("ListDataTransfers: empty response")
	}
	return resp.Transfers, nil
}

// HandleDataTransferDecision wraps dtpb.HandleDataTransferDecision (confirm/reject transfer).
func (a *DataTransferAdapter) HandleDataTransferDecision(ctx context.Context, transferId, userId int64, confirm bool) (*dtpb.DataTransfer, error) {
	if a.Inner == nil || a.Inner.Api == nil {
		return nil, errors.New("data-transfer client not configured")
	}
	req := &dtpb.HandleDataTransferDecisionRequest{
		TransferId: transferId,
		UserId:     userId,
		Confirm:    confirm,
	}
	resp, err := a.Inner.Api.HandleDataTransferDecision(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("HandleDataTransferDecision: empty response")
	}
	return resp.Transfer, nil
}
