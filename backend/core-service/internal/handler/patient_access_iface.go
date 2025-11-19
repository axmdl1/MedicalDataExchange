package handler

import (
	"context"

	dtpb "github.com/axmdl1/MedicalDataExchange/core-service/pkg/gen/go/data-transfer"
)

type PatientAccessDT interface {
	CreatePatientAccessRequest(ctx context.Context, req *dtpb.CreatePatientAccessRequestRequest) (*dtpb.CreatePatientAccessRequestResponse, error)
	ApprovePatientAccessRequest(ctx context.Context, req *dtpb.ApprovePatientAccessRequestRequest) (*dtpb.ApprovePatientAccessRequestResponse, error)
	RejectPatientAccessRequest(ctx context.Context, req *dtpb.RejectPatientAccessRequestRequest) (*dtpb.RejectPatientAccessRequestResponse, error)
	GetPatientAccessRequest(ctx context.Context, req *dtpb.GetPatientAccessRequestRequest) (*dtpb.GetPatientAccessRequestResponse, error)
	ListPatientAccessRequests(ctx context.Context, req *dtpb.ListPatientAccessRequestsRequest) (*dtpb.ListPatientAccessRequestsResponse, error)
	GetTemporaryPatientData(ctx context.Context, req *dtpb.GetTemporaryPatientDataRequest) (*dtpb.GetTemporaryPatientDataResponse, error)
	RevokePatientAccess(ctx context.Context, req *dtpb.RevokePatientAccessRequest) (*dtpb.RevokePatientAccessResponse, error)
}
