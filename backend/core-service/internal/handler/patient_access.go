package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/axmdl1/MedicalDataExchange/core-service/internal/blockchain"
	"github.com/axmdl1/MedicalDataExchange/core-service/internal/middleware"
	dtpb "github.com/axmdl1/MedicalDataExchange/core-service/pkg/gen/go/data-transfer"
	"github.com/go-chi/chi/v5"
)

// -------------------------
// Logging setup (file)
// -------------------------
var (
	paLog     *log.Logger
	paLogOnce sync.Once
	logFile   *os.File
)

func initPatientAccessLogger() {
	// ensure only once
	paLogOnce.Do(func() {
		// put log file next to binary (working dir)
		const fname = "patient_access.log"
		// open (create) file
		f, err := os.OpenFile(filepath.Clean(fname), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			// fallback to stderr logger if file cannot be opened
			log.Printf("patient_access: cannot open log file %s: %v", fname, err)
			paLog = log.New(os.Stderr, "patient_access ", log.LstdFlags|log.Lmicroseconds|log.LUTC)
			return
		}
		logFile = f
		paLog = log.New(f, "patient_access ", log.LstdFlags|log.Lmicroseconds|log.LUTC)
		paLog.Printf("=== patient_access logger started: %s ===", time.Now().UTC().Format(time.RFC3339))
	})
}

// call on shutdown if you want to close file
func closePatientAccessLogger() {
	if logFile != nil {
		paLog.Printf("=== patient_access logger stopping: %s ===", time.Now().UTC().Format(time.RFC3339))
		_ = logFile.Close()
	}
}

// -------------------------
// Helpers
// -------------------------

// writeJSON writes a JSON response and logs it (info).
func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	_ = enc.Encode(payload)
}

// writeJSONError logs error with stack and writes standard JSON error response.
func writeJSONError(w http.ResponseWriter, status int, userMsg string, internalErr error, ctxInfo string) {
	// ensure logger is initialized
	initPatientAccessLogger()

	// build log message with stack
	var stack string
	if internalErr != nil {
		stack = string(debug.Stack())
		paLog.Printf("ERROR: status=%d msg=%s internal=%v ctx=%s\nSTACK:\n%s", status, userMsg, internalErr, ctxInfo, stack)
	} else {
		paLog.Printf("ERROR: status=%d msg=%s ctx=%s", status, userMsg, ctxInfo)
	}

	// send safe message to client
	writeJSON(w, status, map[string]string{"error": userMsg})
}

// logInfo - informational log
func logInfo(format string, args ...interface{}) {
	initPatientAccessLogger()
	paLog.Printf("INFO: "+format, args...)
}

// -------------------------
// Handler
// -------------------------

// PatientAccessHandler — полностью новая версия
type PatientAccessHandler struct {
	DT PatientAccessDT
	BC *blockchain.Blockchain
}

// NewPatientAccessHandler создаёт handler (логгер инициализируется лениво)
func NewPatientAccessHandler(dt PatientAccessDT, bc *blockchain.Blockchain) *PatientAccessHandler {
	// initialize logger now
	initPatientAccessLogger()
	return &PatientAccessHandler{
		DT: dt,
		BC: bc,
	}
}

// ------------------------------------
// 1. Create request (patient)
// POST /patient-access/request
// ------------------------------------
func (h *PatientAccessHandler) CreateAccessRequest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		logInfo("CreateAccessRequest finished in %s from %s", time.Since(start), r.RemoteAddr)
	}()

	var body struct {
		PatientID     int64 `json:"patient_id"`
		ClinicID      int64 `json:"clinic_id"`
		MedicalDataID int64 `json:"medical_data_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON", err, "CreateAccessRequest decode")
		return
	}

	if body.PatientID == 0 || body.ClinicID == 0 || body.MedicalDataID == 0 {
		writeJSONError(w, http.StatusBadRequest, "missing fields", nil, fmt.Sprintf("payload=%+v", body))
		return
	}

	req := &dtpb.CreatePatientAccessRequestRequest{
		PatientId:     body.PatientID,
		ClinicId:      body.ClinicID,
		MedicalDataId: body.MedicalDataID,
	}

	logInfo("CreateAccessRequest: patient=%d clinic=%d record=%d", body.PatientID, body.ClinicID, body.MedicalDataID)
	resp, err := h.DT.CreatePatientAccessRequest(r.Context(), req)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to create access request", err, fmt.Sprintf("patient=%d clinic=%d record=%d", body.PatientID, body.ClinicID, body.MedicalDataID))
		return
	}

	logInfo("CreateAccessRequest: created request id=%d", resp.Request.GetId())
	writeJSON(w, http.StatusCreated, resp.Request)
}

// ------------------------------------
// 2. Approve request (clinic)
// POST /patient-access/approve/{id}
// ------------------------------------
func (h *PatientAccessHandler) ApproveAccessRequest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		logInfo("ApproveAccessRequest finished in %s from %s", time.Since(start), r.RemoteAddr)
	}()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, 400, "invalid id", err, idStr)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		writeJSONError(w, 401, "unauthorized", nil, "missing user in context")
		return
	}

	req := &dtpb.ApprovePatientAccessRequestRequest{
		RequestId:  id,
		ApproverId: user.Id,
	}

	resp, err := h.DT.ApprovePatientAccessRequest(r.Context(), req)
	if err != nil {
		writeJSONError(w, 500, "failed to approve", err, fmt.Sprintf("id=%d", id))
		return
	}

	// ---------- FIX: правильный JSON ----------
	out := map[string]interface{}{
		"access_token":    resp.TemporaryData.GetAccessToken(),
		"expires_at":      resp.TemporaryData.GetExpiresAt(),
		"request_id":      resp.Request.GetId(),
		"patient_id":      resp.Request.GetPatientId(),
		"clinic_id":       resp.Request.GetClinicId(),
		"medical_data_id": resp.Request.GetMedicalDataId(),
		"status":          resp.Request.GetStatus(),
	}

	if resp.Request.GetBlockchainTxId() != "" {
		out["blockchain_tx_id"] = resp.Request.GetBlockchainTxId()
	}

	writeJSON(w, http.StatusOK, out)
}

// ------------------------------------
// 3. Reject request
// ------------------------------------
func (h *PatientAccessHandler) RejectAccessRequest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		logInfo("RejectAccessRequest finished in %s from %s", time.Since(start), r.RemoteAddr)
	}()

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "missing id", nil, "RejectAccessRequest url param")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id", err, idStr)
		return
	}

	logInfo("RejectAccessRequest: request=%d", id)

	req := &dtpb.RejectPatientAccessRequestRequest{
		RequestId: id,
	}

	resp, err := h.DT.RejectPatientAccessRequest(r.Context(), req)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to reject request", err, fmt.Sprintf("request=%d", id))
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// ------------------------------------
// 4. Temporary data access
// GET /patient-access/data?token=xxx
// ------------------------------------
func (h *PatientAccessHandler) GetTemporaryData(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		logInfo("GetTemporaryData finished in %s from %s", time.Since(start), r.RemoteAddr)
	}()

	token := r.URL.Query().Get("token")
	if token == "" {
		writeJSONError(w, http.StatusBadRequest, "missing token", nil, "GetTemporaryData query")
		return
	}

	logInfo("GetTemporaryData: token(len)=%d", len(token))

	req := &dtpb.GetTemporaryPatientDataRequest{
		AccessToken: token,
	}

	resp, err := h.DT.GetTemporaryPatientData(r.Context(), req)
	if err != nil {
		writeJSONError(w, http.StatusForbidden, "access denied or expired", err, fmt.Sprintf("token len=%d", len(token)))
		return
	}

	writeJSON(w, http.StatusOK, resp.MedicalData)
}

// ------------------------------------
// 5. List access requests
// ------------------------------------
func (h *PatientAccessHandler) ListAccessRequests(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		logInfo("ListAccessRequests finished in %s from %s", time.Since(start), r.RemoteAddr)
	}()

	var patientID *int64
	var clinicID *int64
	var status *string

	if v := r.URL.Query().Get("patient_id"); v != "" {
		x, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			patientID = &x
		} else {
			writeJSONError(w, http.StatusBadRequest, "invalid patient_id", err, v)
			return
		}
	}
	if v := r.URL.Query().Get("clinic_id"); v != "" {
		x, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			clinicID = &x
		} else {
			writeJSONError(w, http.StatusBadRequest, "invalid clinic_id", err, v)
			return
		}
	}
	if v := r.URL.Query().Get("status"); v != "" {
		status = &v
	}

	req := &dtpb.ListPatientAccessRequestsRequest{
		PatientId: patientID,
		ClinicId:  clinicID,
		Status:    status,
	}

	resp, err := h.DT.ListPatientAccessRequests(r.Context(), req)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to list requests", err, fmt.Sprintf("q=%v", r.URL.RawQuery))
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// ------------------------------------
// 6. Revoke
// ------------------------------------
func (h *PatientAccessHandler) RevokeAccess(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		logInfo("RevokeAccess finished in %s from %s", time.Since(start), r.RemoteAddr)
	}()

	var body struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON", err, "RevokeAccess decode")
		return
	}

	if body.AccessToken == "" {
		writeJSONError(w, http.StatusBadRequest, "missing access_token", nil, "RevokeAccess body empty")
		return
	}

	req := &dtpb.RevokePatientAccessRequest{
		AccessToken: body.AccessToken,
	}

	resp, err := h.DT.RevokePatientAccess(r.Context(), req)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to revoke access", err, fmt.Sprintf("token len=%d", len(body.AccessToken)))
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
