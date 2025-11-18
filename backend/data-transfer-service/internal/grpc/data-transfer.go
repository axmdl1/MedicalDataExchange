package grpc

import (
	"context"
	"log"
	"time"

	dtpb "github.com/axmdl1/MedicalDataExchange/core-service/pkg/gen/go/data-transfer"
	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/model"
	"github.com/axmdl1/MedicalDataExchange/data-transfer-service/internal/service"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	dtpb.UnimplementedDataTransferServiceServer
	dataTransferService  service.DataTransferService
	patientAccessService service.PatientAccessService
}

func RegisterServerAPI(grpcServer *grpc.Server, dtService service.DataTransferService, patientAccessService service.PatientAccessService) {
	dtpb.RegisterDataTransferServiceServer(grpcServer, &serverAPI{
		dataTransferService:  dtService,
		patientAccessService: patientAccessService,
	})
}

// --- MedicalData ---

func (s *serverAPI) CreateMedicalData(ctx context.Context, req *dtpb.CreateMedicalDataRequest) (*dtpb.CreateMedicalDataResponse, error) {
	log := zerolog.Ctx(ctx).With().
		Str("component", "data_transfer_grpc").
		Str("method", "CreateMedicalData").
		Logger()

	if req.MedicalData == nil {
		return nil, status.Error(codes.InvalidArgument, "medical_data is required")
	}

	resp, err := s.dataTransferService.CreateMedicalData(ctx, req.MedicalData)
	if err != nil {
		log.Error().Err(err).Msg("service_failed")
		return nil, status.Errorf(codes.Internal, "failed to create medical data: %v", err)
	}

	return &dtpb.CreateMedicalDataResponse{MedicalData: resp}, nil
}

func (s *serverAPI) GetMedicalData(ctx context.Context, req *dtpb.GetMedicalDataRequest) (*dtpb.GetMedicalDataResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	resp, err := s.dataTransferService.GetMedicalData(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get medical data: %v", err)
	}

	return &dtpb.GetMedicalDataResponse{MedicalData: resp}, nil
}

func (s *serverAPI) ListMedicalData(ctx context.Context, req *dtpb.ListMedicalDataRequest) (*dtpb.ListMedicalDataResponse, error) {
	// Получаем user_id и role из контекста (из JWT токена)
	userID, ok := GetUserID(ctx)
	userRole, _ := GetUserRole(ctx)

	// Для пациента принудительно фильтруем по его ID
	var userIDFilter *int64
	if userRole == "patient" && ok {
		// Пациент может видеть только свои данные
		if req.UserId != nil && *req.UserId != userID {
			return nil, status.Error(codes.PermissionDenied, "you can only access your own medical data")
		}
		userIDFilter = &userID
		log.Printf("RBAC: Patient can only see their own medical data (user_id=%d)", userID)
	} else {
		// Для сотрудников и админов используем переданный фильтр
		userIDFilter = req.UserId
	}

	resp, err := s.dataTransferService.ListMedicalData(ctx, userIDFilter, req.ClinicId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list medical data: %v", err)
	}

	// Дополнительная фильтрация для пациента (на всякий случай)
	if userRole == "patient" && ok {
		filteredData := make([]*dtpb.MedicalData, 0)
		for _, data := range resp {
			if data.UserId == userID {
				filteredData = append(filteredData, data)
			}
		}
		resp = filteredData
	}

	return &dtpb.ListMedicalDataResponse{MedicalData: resp}, nil
}

func (s *serverAPI) DeleteMedicalData(ctx context.Context, req *dtpb.DeleteMedicalDataRequest) (*dtpb.DeleteMedicalDataResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.dataTransferService.DeleteMedicalData(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete medical data: %v", err)
	}

	return &dtpb.DeleteMedicalDataResponse{Success: true}, nil
}

// --- DataTransfer ---

func (s *serverAPI) CreateDataTransfer(ctx context.Context, req *dtpb.CreateDataTransferRequest) (*dtpb.CreateDataTransferResponse, error) {
	if req.MedicalDataId == 0 || req.FromClinicId == 0 || req.ToClinicId == 0 || req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "all fields are required")
	}

	transfer := &dtpb.DataTransfer{
		MedicalDataId: req.MedicalDataId,
		FromClinicId:  req.FromClinicId,
		ToClinicId:    req.ToClinicId,
		UserId:        req.UserId,
	}

	resp, err := s.dataTransferService.CreateDataTransfer(ctx, transfer)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create data transfer: %v", err)
	}

	return &dtpb.CreateDataTransferResponse{Transfer: resp}, nil
}

func (s *serverAPI) GetDataTransfer(ctx context.Context, req *dtpb.GetDataTransferRequest) (*dtpb.GetDataTransferResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	resp, err := s.dataTransferService.GetDataTransfer(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get data transfer: %v", err)
	}

	return &dtpb.GetDataTransferResponse{Transfer: resp}, nil
}

func (s *serverAPI) ListDataTransfers(ctx context.Context, req *dtpb.ListDataTransfersRequest) (*dtpb.ListDataTransfersResponse, error) {
	var statusFilter *string
	if req.Status != nil {
		statusFilter = req.Status
	}

	resp, err := s.dataTransferService.ListDataTransfers(ctx, req.ClinicId, req.UserId, statusFilter)
	if err != nil {
		return nil,
			status.Errorf(codes.Internal, "failed to list data transfers: %v", err)
	}

	return &dtpb.ListDataTransfersResponse{Transfers: resp}, nil
}

func (s *serverAPI) HandleDataTransferDecision(ctx context.Context, req *dtpb.HandleDataTransferDecisionRequest) (*dtpb.HandleDataTransferDecisionResponse, error) {
	if req.TransferId == 0 || req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "transfer_id and user_id are required")
	}

	resp, err := s.dataTransferService.HandleDataTransferDecision(ctx, req.TransferId, req.UserId, req.Confirm)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to handle data transfer decision: %v", err)
	}

	return &dtpb.HandleDataTransferDecisionResponse{Transfer: resp}, nil
}

// --- Clinic ---

func (s *serverAPI) CreateClinic(ctx context.Context, req *dtpb.CreateClinicRequest) (*dtpb.CreateClinicResponse, error) {
	log := zerolog.Ctx(ctx).With().
		Str("component", "data_transfer_grpc").
		Str("method", "CreateClinic").
		Logger()

	if req.Clinic == nil {
		return nil, status.Error(codes.InvalidArgument, "clinic is required")
	}

	resp, err := s.dataTransferService.CreateClinic(ctx, req.Clinic)
	if err != nil {
		log.Error().Err(err).Msg("service_failed")
		return nil, status.Errorf(codes.Internal, "failed to create clinic: %v", err)
	}

	return &dtpb.CreateClinicResponse{Clinic: resp}, nil
}

func (s *serverAPI) GetClinic(ctx context.Context, req *dtpb.GetClinicRequest) (*dtpb.GetClinicResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	resp, err := s.dataTransferService.GetClinic(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get clinic: %v", err)
	}

	return &dtpb.GetClinicResponse{Clinic: resp}, nil
}

func (s *serverAPI) UpdateClinic(ctx context.Context, req *dtpb.UpdateClinicRequest) (*dtpb.UpdateClinicResponse, error) {
	log := zerolog.Ctx(ctx).With().
		Str("component", "data_transfer_grpc").
		Str("method", "UpdateClinic").
		Logger()

	if req.Clinic == nil || req.Clinic.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "clinic with id is required")
	}

	resp, err := s.dataTransferService.UpdateClinic(ctx, req.Clinic)
	if err != nil {
		log.Error().Err(err).Msg("service_failed")
		return nil, status.Errorf(codes.Internal, "failed to update clinic: %v", err)
	}

	return &dtpb.UpdateClinicResponse{Clinic: resp}, nil
}

func (s *serverAPI) DeleteClinic(ctx context.Context, req *dtpb.DeleteClinicRequest) (*dtpb.DeleteClinicResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.dataTransferService.DeleteClinic(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete clinic: %v", err)
	}

	return &dtpb.DeleteClinicResponse{Success: true}, nil
}

func (s *serverAPI) ListClinics(ctx context.Context, req *dtpb.ListClinicsRequest) (*dtpb.ListClinicsResponse, error) {
	resp, err := s.dataTransferService.ListClinics(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list clinics: %v", err)
	}

	return &dtpb.ListClinicsResponse{Clinics: resp}, nil
}

// --- Patient Access ---

func (s *serverAPI) CreatePatientAccessRequest(ctx context.Context, req *dtpb.CreatePatientAccessRequestRequest) (*dtpb.CreatePatientAccessRequestResponse, error) {
	// Получаем user_id из контекста (из JWT токена)
	userID, ok := GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	// Для пациента принудительно используем его ID из токена
	userRole, _ := GetUserRole(ctx)
	if userRole == "patient" {
		if req.PatientId != 0 && req.PatientId != userID {
			return nil, status.Error(codes.PermissionDenied, "you can only create requests for yourself")
		}
		req.PatientId = userID
	}

	if req.PatientId == 0 || req.ClinicId == 0 || req.MedicalDataId == 0 {
		return nil, status.Error(codes.InvalidArgument, "patient_id, clinic_id and medical_data_id are required")
	}

	request, err := s.patientAccessService.CreateAccessRequest(ctx, req.PatientId, req.ClinicId, req.MedicalDataId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access request: %v", err)
	}

	return &dtpb.CreatePatientAccessRequestResponse{
		Request: convertPatientAccessRequestToProto(request),
	}, nil
}

func (s *serverAPI) ApprovePatientAccessRequest(ctx context.Context, req *dtpb.ApprovePatientAccessRequestRequest) (*dtpb.ApprovePatientAccessRequestResponse, error) {
	if req.RequestId == 0 || req.ApproverId == 0 {
		return nil, status.Error(codes.InvalidArgument, "request_id and approver_id are required")
	}

	tempData, err := s.patientAccessService.ApproveAccessRequest(ctx, req.RequestId, req.ApproverId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to approve access request: %v", err)
	}

	// Получаем обновленный запрос
	request, err := s.patientAccessService.GetAccessRequest(ctx, req.RequestId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get access request: %v", err)
	}

	return &dtpb.ApprovePatientAccessRequestResponse{
		Request:        convertPatientAccessRequestToProto(request),
		TemporaryData:  convertTemporaryPatientDataToProto(tempData),
	}, nil
}

func (s *serverAPI) RejectPatientAccessRequest(ctx context.Context, req *dtpb.RejectPatientAccessRequestRequest) (*dtpb.RejectPatientAccessRequestResponse, error) {
	if req.RequestId == 0 {
		return nil, status.Error(codes.InvalidArgument, "request_id is required")
	}

	if err := s.patientAccessService.RejectAccessRequest(ctx, req.RequestId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reject access request: %v", err)
	}

	request, err := s.patientAccessService.GetAccessRequest(ctx, req.RequestId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get access request: %v", err)
	}

	return &dtpb.RejectPatientAccessRequestResponse{
		Request: convertPatientAccessRequestToProto(request),
	}, nil
}

func (s *serverAPI) GetPatientAccessRequest(ctx context.Context, req *dtpb.GetPatientAccessRequestRequest) (*dtpb.GetPatientAccessRequestResponse, error) {
	if req.RequestId == 0 {
		return nil, status.Error(codes.InvalidArgument, "request_id is required")
	}

	request, err := s.patientAccessService.GetAccessRequest(ctx, req.RequestId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get access request: %v", err)
	}

	return &dtpb.GetPatientAccessRequestResponse{
		Request: convertPatientAccessRequestToProto(request),
	}, nil
}

func (s *serverAPI) ListPatientAccessRequests(ctx context.Context, req *dtpb.ListPatientAccessRequestsRequest) (*dtpb.ListPatientAccessRequestsResponse, error) {
	// Получаем user_id и role из контекста (из JWT токена)
	userID, ok := GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	userRole, _ := GetUserRole(ctx)

	var patientID, clinicID *int64
	var status *string

	// Для пациента принудительно фильтруем по его ID
	if userRole == "patient" {
		patientID = &userID
		log.Printf("RBAC: Patient can only see their own requests (patient_id=%d)", userID)
	} else {
		// Для сотрудников и админов используем переданные фильтры
		if req.PatientId != nil {
			patientID = req.PatientId
		}
	}

	if req.ClinicId != nil {
		clinicID = req.ClinicId
	}
	if req.Status != nil {
		status = req.Status
	}

	requests, err := s.patientAccessService.ListPatientAccessRequests(ctx, patientID, clinicID, status)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list access requests: %v", err)
	}

	// Дополнительная фильтрация для пациента (на всякий случай)
	if userRole == "patient" {
		filteredRequests := make([]*model.PatientAccessRequest, 0)
		for _, req := range requests {
			if req.PatientID == userID {
				filteredRequests = append(filteredRequests, req)
			}
		}
		requests = filteredRequests
	}

	protoRequests := make([]*dtpb.PatientAccessRequest, len(requests))
	for i, req := range requests {
		protoRequests[i] = convertPatientAccessRequestToProto(req)
	}

	return &dtpb.ListPatientAccessRequestsResponse{
		Requests: protoRequests,
	}, nil
}

func (s *serverAPI) GetTemporaryPatientData(ctx context.Context, req *dtpb.GetTemporaryPatientDataRequest) (*dtpb.GetTemporaryPatientDataResponse, error) {
	if req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access_token is required")
	}

	medData, err := s.patientAccessService.GetTemporaryData(ctx, req.AccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get temporary data: %v", err)
	}

	return &dtpb.GetTemporaryPatientDataResponse{
		MedicalData: convertMedicalDataToProto(medData),
	}, nil
}

func (s *serverAPI) RevokePatientAccess(ctx context.Context, req *dtpb.RevokePatientAccessRequest) (*dtpb.RevokePatientAccessResponse, error) {
	if req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access_token is required")
	}

	if err := s.patientAccessService.RevokeAccess(ctx, req.AccessToken); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to revoke access: %v", err)
	}

	return &dtpb.RevokePatientAccessResponse{Success: true}, nil
}

// Helper functions to convert models to proto

func convertPatientAccessRequestToProto(req *model.PatientAccessRequest) *dtpb.PatientAccessRequest {
	proto := &dtpb.PatientAccessRequest{
		Id:            req.ID,
		PatientId:     req.PatientID,
		ClinicId:      req.ClinicID,
		MedicalDataId: req.MedicalDataID,
		Status:        req.Status,
		BlockchainTxId: req.BlockchainTxID,
		RequestedAt:   req.RequestedAt.Format(time.RFC3339),
	}
	if req.ApprovedAt != nil {
		proto.ApprovedAt = req.ApprovedAt.Format(time.RFC3339)
	}
	if req.ExpiresAt != nil {
		proto.ExpiresAt = req.ExpiresAt.Format(time.RFC3339)
	}
	return proto
}

func convertTemporaryPatientDataToProto(temp *model.TemporaryPatientData) *dtpb.TemporaryPatientData {
	return &dtpb.TemporaryPatientData{
		Id:            temp.ID,
		AccessToken:   temp.AccessToken,
		PatientId:     temp.PatientID,
		ClinicId:      temp.ClinicID,
		MedicalDataId: temp.MedicalDataID,
		GrantedAt:     temp.GrantedAt.Format(time.RFC3339),
		ExpiresAt:     temp.ExpiresAt.Format(time.RFC3339),
		IsRevoked:     temp.IsRevoked,
	}
}

func convertMedicalDataToProto(med *model.MedicalData) *dtpb.MedicalData {
	return &dtpb.MedicalData{
		Id:          med.ID,
		UserId:      med.UserID,
		ClinicId:    med.ClinicID,
		Diagnosis:   med.Diagnosis,
		Complaint:   med.Complaint,
		Treatment:   med.Treatment,
		Medications: med.Medications,
		Allergies:   med.Allergies,
		DoctorNotes: med.DoctorNotes,
		LabResults:  med.LabResults,
		CreatedAt:   med.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   med.UpdatedAt.Format(time.RFC3339),
	}
}
