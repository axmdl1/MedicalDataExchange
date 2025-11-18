package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	dtpb "github.com/axmdl1/MedicalDataExchange/core-service/pkg/gen/go/data-transfer"
	"github.com/axmdl1/MedicalDataExchange/core-service/internal/client"
	jsonpkg "github.com/axmdl1/MedicalDataExchange/core-service/pkg/json"
)

type PatientAccessHandler struct {
	dtClient *client.DataTransferClient
}

func NewPatientAccessHandler(dtClient *client.DataTransferClient) *PatientAccessHandler {
	return &PatientAccessHandler{
		dtClient: dtClient,
	}
}

// Структуры для запросов/ответов
type CreateAccessRequestInput struct {
	PatientID     int64 `json:"patient_id"`
	ClinicID      int64 `json:"clinic_id"`
	MedicalDataID int64 `json:"medical_data_id"`
}

type CreateAccessRequestResponse struct {
	RequestID int64  `json:"request_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

type ApproveAccessRequestInput struct {
	ApproverID int64 `json:"approver_id"`
}

type ApproveAccessRequestResponse struct {
	RequestID   int64  `json:"request_id"`
	AccessToken string `json:"access_token"`
	ExpiresAt   string `json:"expires_at"`
	Status      string `json:"status"`
	Message     string `json:"message"`
}

type AccessRequestInfo struct {
	ID            int64  `json:"id"`
	PatientID     int64  `json:"patient_id"`
	ClinicID      int64  `json:"clinic_id"`
	MedicalDataID int64  `json:"medical_data_id"`
	Status        string `json:"status"`
	RequestedAt   string `json:"requested_at"`
	ApprovedAt    string `json:"approved_at,omitempty"`
	ExpiresAt     string `json:"expires_at,omitempty"`
}

type MedicalDataResponse struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	ClinicID    int64  `json:"clinic_id"`
	Diagnosis   string `json:"diagnosis"`
	Complaint   string `json:"complaint"`
	Treatment   string `json:"treatment"`
	Medications string `json:"medications"`
	Allergies   string `json:"allergies"`
	DoctorNotes string `json:"doctor_notes"`
	LabResults  string `json:"lab_results"`
	CreatedAt   string `json:"created_at"`
}

// POST /api/patient-access/request - Создать запрос на доступ
func (h *PatientAccessHandler) CreateAccessRequest(w http.ResponseWriter, r *http.Request) {
	var input CreateAccessRequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonpkg.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Вызов к data-transfer-service
	// patient_id будет автоматически установлен из JWT токена в data-transfer-service
	// для пациента, так что здесь просто передаем то, что пришло
	ctx := r.Context()
	req := &dtpb.CreatePatientAccessRequestRequest{
		PatientId:     input.PatientID,
		ClinicId:      input.ClinicID,
		MedicalDataId: input.MedicalDataID,
	}

	resp, err := h.dtClient.Api.CreatePatientAccessRequest(ctx, req)
	if err != nil {
		jsonpkg.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := CreateAccessRequestResponse{
		RequestID: resp.Request.Id,
		Status:    resp.Request.Status,
		Message:   "Запрос на доступ создан. Ожидайте одобрения клиникой.",
	}

	jsonpkg.WriteJSON(w, http.StatusCreated, response)
}

// POST /api/patient-access/approve/{id} - Одобрить запрос
func (h *PatientAccessHandler) ApproveAccessRequest(w http.ResponseWriter, r *http.Request) {
	requestID, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	var input ApproveAccessRequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonpkg.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Вызов к data-transfer-service для одобрения
	ctx := r.Context()
	req := &dtpb.ApprovePatientAccessRequestRequest{
		RequestId:  requestID,
		ApproverId: input.ApproverID,
	}

	resp, err := h.dtClient.Api.ApprovePatientAccessRequest(ctx, req)
	if err != nil {
		jsonpkg.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := ApproveAccessRequestResponse{
		RequestID:   resp.Request.Id,
		AccessToken: resp.TemporaryData.AccessToken,
		ExpiresAt:   resp.TemporaryData.ExpiresAt,
		Status:      resp.Request.Status,
		Message:     "Доступ одобрен. Токен действителен 15 минут.",
	}

	jsonpkg.WriteJSON(w, http.StatusOK, response)
}

// GET /api/patient-access/data?token=xxx - Получить данные по токену
func (h *PatientAccessHandler) GetTemporaryData(w http.ResponseWriter, r *http.Request) {
	accessToken := r.URL.Query().Get("token")
	if accessToken == "" {
		jsonpkg.WriteJSON(w, http.StatusBadRequest,
			map[string]string{"error": "token parameter is required"})
		return
	}

	// Вызов к data-transfer-service для получения временных данных
	ctx := r.Context()
	req := &dtpb.GetTemporaryPatientDataRequest{
		AccessToken: accessToken,
	}

	resp, err := h.dtClient.Api.GetTemporaryPatientData(ctx, req)
	if err != nil {
		jsonpkg.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	medData := resp.MedicalData
	response := MedicalDataResponse{
		ID:          medData.Id,
		UserID:      medData.UserId,
		ClinicID:    medData.ClinicId,
		Diagnosis:   medData.Diagnosis,
		Complaint:   medData.Complaint,
		Treatment:   medData.Treatment,
		Medications: medData.Medications,
		Allergies:   medData.Allergies,
		DoctorNotes: medData.DoctorNotes,
		LabResults:  medData.LabResults,
		CreatedAt:   medData.CreatedAt,
	}

	jsonpkg.WriteJSON(w, http.StatusOK, response)
}

// GET /api/patient-access/requests - Список запросов
func (h *PatientAccessHandler) ListAccessRequests(w http.ResponseWriter, r *http.Request) {
	patientIDStr := r.URL.Query().Get("patient_id")
	clinicIDStr := r.URL.Query().Get("clinic_id")
	status := r.URL.Query().Get("status")

	ctx := r.Context()
	req := &dtpb.ListPatientAccessRequestsRequest{}

	// patient_id будет автоматически установлен из JWT токена в data-transfer-service
	// для пациента, но мы все равно передаем его для совместимости
	if patientIDStr != "" {
		if patientID, err := strconv.ParseInt(patientIDStr, 10, 64); err == nil {
			req.PatientId = &patientID
		}
	}
	if clinicIDStr != "" {
		if clinicID, err := strconv.ParseInt(clinicIDStr, 10, 64); err == nil {
			req.ClinicId = &clinicID
		}
	}
	if status != "" {
		req.Status = &status
	}

	resp, err := h.dtClient.Api.ListPatientAccessRequests(ctx, req)
	if err != nil {
		jsonpkg.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	requests := make([]AccessRequestInfo, len(resp.Requests))
	for i, req := range resp.Requests {
		requests[i] = AccessRequestInfo{
			ID:            req.Id,
			PatientID:     req.PatientId,
			ClinicID:      req.ClinicId,
			MedicalDataID: req.MedicalDataId,
			Status:        req.Status,
			RequestedAt:   req.RequestedAt,
			ApprovedAt:    req.ApprovedAt,
			ExpiresAt:     req.ExpiresAt,
		}
	}

	jsonpkg.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"requests": requests,
		"total":    len(requests),
	})
}

// POST /api/patient-access/revoke - Отозвать доступ
func (h *PatientAccessHandler) RevokeAccess(w http.ResponseWriter, r *http.Request) {
	var input struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonpkg.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Вызов к data-transfer-service для отзыва доступа
	ctx := r.Context()
	req := &dtpb.RevokePatientAccessRequest{
		AccessToken: input.AccessToken,
	}

	_, err := h.dtClient.Api.RevokePatientAccess(ctx, req)
	if err != nil {
		jsonpkg.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Доступ успешно отозван",
	})
}

// GET /api/patient-access/request/{id} - Получить конкретный запрос
func (h *PatientAccessHandler) GetAccessRequest(w http.ResponseWriter, r *http.Request) {
	requestID, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	// Вызов к data-transfer-service
	ctx := r.Context()
	req := &dtpb.GetPatientAccessRequestRequest{
		RequestId: requestID,
	}

	resp, err := h.dtClient.Api.GetPatientAccessRequest(ctx, req)
	if err != nil {
		jsonpkg.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	reqProto := resp.Request
	request := AccessRequestInfo{
		ID:            reqProto.Id,
		PatientID:     reqProto.PatientId,
		ClinicID:      reqProto.ClinicId,
		MedicalDataID: reqProto.MedicalDataId,
		Status:        reqProto.Status,
		RequestedAt:   reqProto.RequestedAt,
		ApprovedAt:    reqProto.ApprovedAt,
		ExpiresAt:     reqProto.ExpiresAt,
	}

	jsonpkg.WriteJSON(w, http.StatusOK, request)
}

// POST /api/patient-access/reject/{id} - Отклонить запрос
func (h *PatientAccessHandler) RejectAccessRequest(w http.ResponseWriter, r *http.Request) {
	requestID, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	// Вызов к data-transfer-service
	ctx := r.Context()
	req := &dtpb.RejectPatientAccessRequestRequest{
		RequestId: requestID,
	}

	resp, err := h.dtClient.Api.RejectPatientAccessRequest(ctx, req)
	if err != nil {
		jsonpkg.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	reqProto := resp.Request
	request := AccessRequestInfo{
		ID:            reqProto.Id,
		PatientID:     reqProto.PatientId,
		ClinicID:      reqProto.ClinicId,
		MedicalDataID: reqProto.MedicalDataId,
		Status:        reqProto.Status,
		RequestedAt:   reqProto.RequestedAt,
		ApprovedAt:    reqProto.ApprovedAt,
		ExpiresAt:     reqProto.ExpiresAt,
	}

	jsonpkg.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Запрос отклонен",
		"request": request,
	})
}
