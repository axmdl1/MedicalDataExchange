package handler

import (
	"net/http"
	"strconv"

	"github.com/axmdl1/MedicalDataExchange/core-service/internal/client"
	dtpb "github.com/axmdl1/MedicalDataExchange/core-service/pkg/gen/go/data-transfer"
	"github.com/axmdl1/MedicalDataExchange/core-service/pkg/json"
)

type DataTransferHandler struct {
	client *client.DataTransferClient
}

func NewDataTransferHandler(c *client.DataTransferClient) *DataTransferHandler {
	return &DataTransferHandler{client: c}
}

// POST /medical-data
func (h *DataTransferHandler) CreateMedicalData(w http.ResponseWriter, r *http.Request) {
	var req dtpb.CreateMedicalDataRequest
	if err := json.ParseProtoJSON(r.Body, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}
	resp, err := h.client.Api.CreateMedicalData(r.Context(), &req)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	json.WriteJSON(w, http.StatusCreated, resp)
}

// GET /medical-data/{id}
func (h *DataTransferHandler) GetMedicalData(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	resp, err := h.client.Api.GetMedicalData(r.Context(), &dtpb.GetMedicalDataRequest{Id: id})
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	json.WriteJSON(w, http.StatusOK, resp)
}

// GET /medical-data
func (h *DataTransferHandler) ListMedicalData(w http.ResponseWriter, r *http.Request) {
	userID := parseOptionalInt64(r.URL.Query().Get("user_id"))
	clinicID := parseOptionalInt64(r.URL.Query().Get("clinic_id"))

	resp, err := h.client.Api.ListMedicalData(r.Context(), &dtpb.ListMedicalDataRequest{
		UserId:   userID,   // передаем указатель, nil допустим
		ClinicId: clinicID, // передаем указатель, nil допустим
	})
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	json.WriteJSON(w, http.StatusOK, resp)
}

// DELETE /medical-data/{id}
func (h *DataTransferHandler) DeleteMedicalData(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	_, err := h.client.Api.DeleteMedicalData(r.Context(), &dtpb.DeleteMedicalDataRequest{Id: id})
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	json.WriteJSON(w, http.StatusOK, dtpb.DeleteMedicalDataResponse{Success: true})
}

// POST /transfers
func (h *DataTransferHandler) CreateDataTransfer(w http.ResponseWriter, r *http.Request) {
	var req dtpb.CreateDataTransferRequest
	if err := json.ParseProtoJSON(r.Body, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}
	resp, err := h.client.Api.CreateDataTransfer(r.Context(), &req)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	json.WriteJSON(w, http.StatusCreated, resp)
}

// GET /transfers/{id}
func (h *DataTransferHandler) GetDataTransfer(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	resp, err := h.client.Api.GetDataTransfer(r.Context(), &dtpb.GetDataTransferRequest{Id: id})
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	json.WriteJSON(w, http.StatusOK, resp)
}

// GET /transfers
func (h *DataTransferHandler) ListDataTransfers(w http.ResponseWriter, r *http.Request) {
	clinicID := parseOptionalInt64(r.URL.Query().Get("clinic_id"))
	userID := parseOptionalInt64(r.URL.Query().Get("user_id"))
	statusParam := r.URL.Query().Get("status")

	// Если строка пустая, передаем nil
	var status *string
	if statusParam != "" {
		status = &statusParam
	}

	resp, err := h.client.Api.ListDataTransfers(r.Context(), &dtpb.ListDataTransfersRequest{
		ClinicId: clinicID, // *int64
		UserId:   userID,   // *int64
		Status:   status,   // *string
	})
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	json.WriteJSON(w, http.StatusOK, resp)
}

// POST /transfers/decision
func (h *DataTransferHandler) HandleDecision(w http.ResponseWriter, r *http.Request) {
	var req dtpb.HandleDataTransferDecisionRequest
	if err := json.ParseProtoJSON(r.Body, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}
	resp, err := h.client.Api.HandleDataTransferDecision(r.Context(), &req)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	json.WriteJSON(w, http.StatusOK, resp)
}

// вспомогательная функция
func parseOptionalInt64(s string) *int64 {
	if s == "" {
		return nil
	}
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		return &v
	}
	return nil
}
