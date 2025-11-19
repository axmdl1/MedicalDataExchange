package client

import (
	"context"
	"errors"
	"strconv"
	"sync"

	dtpb "github.com/axmdl1/MedicalDataExchange/core-service/pkg/gen/go/data-transfer"
)

/*
   DataTransferAdapter — это мост между HTTP handlers и реальным gRPC клиентом.
   Именно он должен передаваться в handlers.
*/

type InMemoryAccessRequest struct {
	ID           int64
	PatientID    int64
	ClinicID     int64
	RecordID     int64
	Status       string
	BlockchainTx string
	RequestedAt  string
	ApprovedAt   string
	ExpiresAt    string
}

type DataTransferAdapter struct {
	api            dtpb.DataTransferServiceClient
	mu             sync.RWMutex
	accessRequests map[int64]*InMemoryAccessRequest
}

func (a *DataTransferAdapter) CreatePatientAccessRequest(ctx context.Context, req *dtpb.CreatePatientAccessRequestRequest) (*dtpb.CreatePatientAccessRequestResponse, error) {
	return a.api.CreatePatientAccessRequest(ctx, req)
}

func (a *DataTransferAdapter) ApprovePatientAccessRequest(ctx context.Context, req *dtpb.ApprovePatientAccessRequestRequest) (*dtpb.ApprovePatientAccessRequestResponse, error) {
	return a.api.ApprovePatientAccessRequest(ctx, req)
}

func (a *DataTransferAdapter) RejectPatientAccessRequest(ctx context.Context, req *dtpb.RejectPatientAccessRequestRequest) (*dtpb.RejectPatientAccessRequestResponse, error) {
	return a.api.RejectPatientAccessRequest(ctx, req)
}

func (a *DataTransferAdapter) GetPatientAccessRequest(ctx context.Context, req *dtpb.GetPatientAccessRequestRequest) (*dtpb.GetPatientAccessRequestResponse, error) {
	return a.api.GetPatientAccessRequest(ctx, req)
}

func (a *DataTransferAdapter) ListPatientAccessRequests(ctx context.Context, req *dtpb.ListPatientAccessRequestsRequest) (*dtpb.ListPatientAccessRequestsResponse, error) {
	return a.api.ListPatientAccessRequests(ctx, req)
}

func (a *DataTransferAdapter) GetTemporaryPatientData(ctx context.Context, req *dtpb.GetTemporaryPatientDataRequest) (*dtpb.GetTemporaryPatientDataResponse, error) {
	return a.api.GetTemporaryPatientData(ctx, req)
}

func (a *DataTransferAdapter) RevokePatientAccess(ctx context.Context, req *dtpb.RevokePatientAccessRequest) (*dtpb.RevokePatientAccessResponse, error) {
	return a.api.RevokePatientAccess(ctx, req)
}

func NewDataTransferAdapter(c *DataTransferClient) *DataTransferAdapter {
	return &DataTransferAdapter{
		api:            c.Api,
		accessRequests: make(map[int64]*InMemoryAccessRequest),
	}
}

// 1. CreateAccessRequest
func (a *DataTransferAdapter) CreateAccessRequest(
	ctx context.Context,
	patientID int64,
	clinicID int64,
	recordID int64,
) (int64, error) {

	req := &dtpb.CreatePatientAccessRequestRequest{
		PatientId:     patientID,
		ClinicId:      clinicID,
		MedicalDataId: recordID,
	}

	resp, err := a.api.CreatePatientAccessRequest(ctx, req)
	if err != nil {
		return 0, err
	}

	// сохранить копию в memory (не обязательно — только для удобства фронта)
	a.mu.Lock()
	defer a.mu.Unlock()

	a.accessRequests[resp.Request.Id] = &InMemoryAccessRequest{
		ID:        resp.Request.Id,
		PatientID: resp.Request.PatientId,
		ClinicID:  resp.Request.ClinicId,
		RecordID:  resp.Request.MedicalDataId,
		Status:    resp.Request.Status,
	}

	return resp.Request.Id, nil
}

//
// 2. UpdateRequestStatus
//
/*func (a *DataTransferAdapter) UpdateRequestStatus(ctx context.Context, requestID int64, status string) error {
	// gRPC вызов
	_, err := a.api.UpdatePatientAccessStatus(ctx, &dtpb.UpdatePatientAccessStatusRequest{
		RequestId: requestID,
		Status:    status,
	})
	if err != nil {
		return err
	}

	// обновить локально
	a.mu.Lock()
	defer a.mu.Unlock()
	if ar, ok := a.accessRequests[requestID]; ok {
		ar.Status = status
	}

	return nil
}*/

// 3. RejectAccessRequest
func (a *DataTransferAdapter) RejectAccessRequest(ctx context.Context, id int64) error {
	_, err := a.api.RejectPatientAccessRequest(ctx, &dtpb.RejectPatientAccessRequestRequest{
		RequestId: id,
	})
	if err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	if ar, ok := a.accessRequests[id]; ok {
		ar.Status = "rejected"
	}
	return nil
}

// 4. GetAccessRequest
func (a *DataTransferAdapter) GetAccessRequest(ctx context.Context, id int64) (interface{}, error) {
	resp, err := a.api.GetPatientAccessRequest(ctx, &dtpb.GetPatientAccessRequestRequest{
		RequestId: id,
	})
	if err != nil {
		return nil, err
	}
	return resp.Request, nil
}

// 5. ListAccessRequests
func (a *DataTransferAdapter) ListAccessRequests(ctx context.Context) ([]interface{}, error) {
	resp, err := a.api.ListPatientAccessRequests(ctx, &dtpb.ListPatientAccessRequestsRequest{})
	if err != nil {
		return nil, err
	}

	out := make([]interface{}, 0, len(resp.Requests))
	for _, r := range resp.Requests {
		out = append(out, r)
	}
	return out, nil
}

// 6. GetMedicalData
func (a *DataTransferAdapter) GetMedicalData(ctx context.Context, recordID string) (interface{}, error) {
	id, err := strconv.ParseInt(recordID, 10, 64)
	if err != nil {
		return nil, errors.New("recordID must be int64")
	}

	resp, err := a.api.GetMedicalData(ctx, &dtpb.GetMedicalDataRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return resp.MedicalData, nil
}
