package handler

import (
	"net/http"
	"strconv"

	"github.com/axmdl1/MedicalDataExchange/core-service/internal/client"
	dtpb "github.com/axmdl1/MedicalDataExchange/core-service/pkg/gen/go/data-transfer"
	"github.com/axmdl1/MedicalDataExchange/core-service/pkg/json"
)

type ClinicHandler struct {
	client *client.DataTransferClient
}

func NewClinicHandler(c *client.DataTransferClient) *ClinicHandler {
	return &ClinicHandler{client: c}
}

// POST /clinics
func (h *ClinicHandler) CreateClinic(w http.ResponseWriter, r *http.Request) {
	var req dtpb.CreateClinicRequest
	if err := json.ParseProtoJSON(r.Body, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}
	resp, err := h.client.Api.CreateClinic(r.Context(), &req)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	json.WriteJSON(w, http.StatusCreated, resp)
}

// GET /clinics/{id}
func (h *ClinicHandler) GetClinic(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	resp, err := h.client.Api.GetClinic(r.Context(), &dtpb.GetClinicRequest{Id: id})
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	json.WriteJSON(w, http.StatusOK, resp)
}

// GET /clinics
func (h *ClinicHandler) ListClinics(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.Api.ListClinics(r.Context(), &dtpb.ListClinicsRequest{})
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	json.WriteJSON(w, http.StatusOK, resp)
}

// PATCH /clinics/{id}
func (h *ClinicHandler) UpdateClinic(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	var req dtpb.UpdateClinicRequest
	if err := json.ParseProtoJSON(r.Body, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// Set the ID from path
	if req.Clinic == nil {
		req.Clinic = &dtpb.Clinic{}
	}
	req.Clinic.Id = id

	resp, err := h.client.Api.UpdateClinic(r.Context(), &req)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	json.WriteJSON(w, http.StatusOK, resp)
}

// DELETE /clinics/{id}
func (h *ClinicHandler) DeleteClinic(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	_, err := h.client.Api.DeleteClinic(r.Context(), &dtpb.DeleteClinicRequest{Id: id})
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	json.WriteJSON(w, http.StatusOK, dtpb.DeleteClinicResponse{Success: true})
}
