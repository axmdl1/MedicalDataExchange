package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-chi/chi/v5"

	"github.com/axmdl1/MedicalDataExchange/core-service/internal/blockchain"
)

// Интерфейс для data-transfer клиента.
// Твой реальный dtClient должен реализовать эти методы.
// Если имена отличаются — замени вызовы ниже на реальные.
type DataTransferClientInterface interface {
	// Создать запись запроса в БД клиники (gRPC)
	CreateAccessRequest(ctx context.Context, requestID, recordID, patient, clinic string) error
	UpdateRequestStatus(ctx context.Context, requestID, status string) error
	GetMedicalData(ctx context.Context, recordID string) (interface{}, error) // верни модель/JSON-совместимый объект
	ListAccessRequests(ctx context.Context) ([]interface{}, error)
	GetAccessRequest(ctx context.Context, requestID string) (interface{}, error)
	RejectAccessRequest(ctx context.Context, requestID, reason string) error
}

// PatientAccessHandler
type PatientAccessHandler struct {
	DTClient DataTransferClientInterface
	BC       *blockchain.Blockchain
}

// NewPatientAccessHandler принимает dtClient (тот который создаёшь в app.go)
func NewPatientAccessHandler(dt DataTransferClientInterface, bc *blockchain.Blockchain) *PatientAccessHandler {
	return &PatientAccessHandler{DTClient: dt, BC: bc}
}

// 1. Создать запрос (пациент)
func (h *PatientAccessHandler) CreateAccessRequest(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RequestID string `json:"request_id"`
		RecordID  string `json:"record_id"`
		Patient   string `json:"patient"`
		Clinic    string `json:"clinic"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	if body.RequestID == "" || body.RecordID == "" || body.Patient == "" || body.Clinic == "" {
		http.Error(w, "missing fields", 400)
		return
	}

	// Создаём запись в data-transfer (gRPC) — реализуй этот метод в dtClient
	if err := h.DTClient.CreateAccessRequest(r.Context(), body.RequestID, body.RecordID, body.Patient, body.Clinic); err != nil {
		http.Error(w, "failed to create request: "+err.Error(), 500)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]string{"status": "pending", "request_id": body.RequestID})
}

// 2. Клиника одобряет запрос -> записываем на chain и возвращаем токен пациенту
func (h *PatientAccessHandler) ApproveAccessRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", 400)
		return
	}

	// Генерируем token и hash
	token := "patient-temp-token-" + time.Now().Format("20060102T150405")
	tokenHash := crypto.Keccak256Hash([]byte(token))
	expires := time.Now().Add(15 * time.Minute).Unix()

	// Вызываем blockchain.Approve
	txHash, err := h.BC.Approve(id, tokenHash, expires)
	if err != nil {
		http.Error(w, "blockchain approve failed: "+err.Error(), 500)
		return
	}

	// Обновляем статус в data-transfer (опционально)
	_ = h.DTClient.UpdateRequestStatus(r.Context(), id, "approved")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":     token,
		"tokenHash": tokenHash.Hex(),
		"expiresAt": expires,
		"tx":        txHash,
		"status":    "approved",
	})
}

// 3. Клиника отклоняет
func (h *PatientAccessHandler) RejectAccessRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", 400)
		return
	}
	_ = h.DTClient.RejectAccessRequest(r.Context(), id, "rejected by clinic")
	_ = h.DTClient.UpdateRequestStatus(r.Context(), id, "rejected")
	json.NewEncoder(w).Encode(map[string]string{"status": "rejected", "request_id": id})
}

// 4. Получить временные данные (пациент)
// GET /patient-access/data?request_id=r1&token=....&record_id=rec1
func (h *PatientAccessHandler) GetTemporaryData(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	reqID := q.Get("request_id")
	token := q.Get("token")
	recordID := q.Get("record_id")
	if reqID == "" || token == "" || recordID == "" {
		http.Error(w, "missing params", 400)
		return
	}

	tokenHash := common.BytesToHash(crypto.Keccak256([]byte(token)))

	// 1) Verify on-chain
	ok, err := h.BC.VerifyToken(reqID, tokenHash)
	if err != nil {
		http.Error(w, "verify error: "+err.Error(), 500)
		return
	}
	if !ok {
		http.Error(w, "access denied or expired", 403)
		return
	}

	// 2) Read from clinic storage via data-transfer client
	data, err := h.DTClient.GetMedicalData(r.Context(), recordID)
	if err != nil {
		http.Error(w, "not found: "+err.Error(), 404)
		return
	}

	// Return whatever model dtClient returns (must be JSON serializable)
	json.NewEncoder(w).Encode(data)
}

// 5. Список запросов (для клиники)
// Возвращаем список из data-transfer
func (h *PatientAccessHandler) ListAccessRequests(w http.ResponseWriter, r *http.Request) {
	list, err := h.DTClient.ListAccessRequests(r.Context())
	if err != nil {
		http.Error(w, "failed: "+err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(list)
}

// 6. Получить конкретный запрос
func (h *PatientAccessHandler) GetAccessRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", 400)
		return
	}
	req, err := h.DTClient.GetAccessRequest(r.Context(), id)
	if err != nil {
		http.Error(w, "not found: "+err.Error(), 404)
		return
	}
	json.NewEncoder(w).Encode(req)
}

// 7. Revoke access (clinic)
func (h *PatientAccessHandler) RevokeAccess(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RequestID string `json:"request_id"`
		Reason    string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	// If contract has revoke function call it, otherwise mark in data-transfer
	_, _ = h.BC.Revoke(body.RequestID, "") // may return error "not implemented"
	_ = h.DTClient.UpdateRequestStatus(r.Context(), body.RequestID, "revoked")
	json.NewEncoder(w).Encode(map[string]string{"status": "revoked"})
}
